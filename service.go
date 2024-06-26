package notezero

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/dextryz/notezero/badger"
	"github.com/fiatjaf/eventstore"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

type eventService struct {
	db     eventstore.Store
	cache  *badger.Cache
	relays []string
}

func NewEventService(db eventstore.Store, cache *badger.Cache, relays []string) eventService {
	return eventService{
		db:     db,
		cache:  cache,
		relays: relays,
	}
}

// 1. Check if the event is in the cache
// 2. If not, request event from the set of relays
func (s eventService) RequestEvent(ctx context.Context, code string) (*nostr.Event, error) {

	// Wrap the cache db to be used with a relay interface
	wdb := eventstore.RelayWrapper{Store: s.db}

	// Create a nostr filter from the NIP-19 code
	prefix, data, err := nip19.Decode(code)
	if err != nil {
		return nil, err
	}

	var filter nostr.Filter

	switch v := data.(type) {
	case nostr.EntityPointer:
		filter.Authors = []string{v.PublicKey}
		filter.Tags = nostr.TagMap{
			"d": []string{v.Identifier},
		}
		if v.Kind != 0 {
			filter.Kinds = append(filter.Kinds, v.Kind)
		}
	case string:
		if prefix == "npub" {
			filter.Authors = []string{v}
			filter.Kinds = []int{0}
		}
	default:
		return nil, fmt.Errorf("code type not supported: %s", code)
	}

	// Try to fetch in our internal eventstore (cache) first
	events, err := wdb.QuerySync(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(events) != 0 {
		return events[0], nil
	}

	// No events found in cache, request relays and publish to cache
	events = s.queryRelays(ctx, filter)
	for _, e := range events {
		err := wdb.Publish(ctx, *e)
		if err != nil {
			return nil, err
		}
	}

	return events[0], nil
}

func (s eventService) AuthorArticles(ctx context.Context, npub string) ([]*nostr.Event, error) {

	_, pk, err := nip19.Decode(npub)
	if err != nil {
		return nil, err
	}

	filter := nostr.Filter{
		Kinds:   []int{nostr.KindArticle},
		Authors: []string{pk.(string)},
		Limit:   500,
	}

	// fetch from local store if available
	wdb := eventstore.RelayWrapper{Store: s.db}

	// Try to fetch in our internal eventstore (cache) first
	events, err := wdb.QuerySync(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(events) != 0 {
		return events, nil
	}

	// No events found in cache, request relays and publish to cache
	events = s.queryRelays(ctx, filter)
	for _, e := range events {
		err := wdb.Publish(ctx, *e)
		if err != nil {
			return nil, err
		}
	}

	// sort before returning
	slices.SortFunc(events, func(a, b *nostr.Event) int { return int(b.CreatedAt - a.CreatedAt) })

	return events, nil
}

func (s eventService) ArticleHighlights(ctx context.Context, kind int, pubkey, identifier string) ([]*nostr.Event, error) {

	wdb := eventstore.RelayWrapper{Store: s.db}

	pool := nostr.NewSimplePool(context.Background())

	// 2. Article is cached, so pull highlights

	tag := fmt.Sprintf("%d:%s:%s", kind, pubkey, identifier)

	filter := nostr.Filter{
		Kinds: []int{9802},
		Tags: nostr.TagMap{
			"a": []string{tag},
		},
		Limit: 500,
	}

	var lastNotes []*nostr.Event

	// fetch from external relays asynchronously
	external := make(chan []*nostr.Event)
	go func() {
		notes := make([]*nostr.Event, 0, filter.Limit)
		defer func() {
			external <- notes
		}()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		ch := pool.SubManyEose(ctx, s.relays, nostr.Filters{filter})
		for {
			select {
			case ie, more := <-ch:
				if !more {
					return
				}
				notes = append(notes, ie.Event)
				s.db.SaveEvent(ctx, ie.Event)
				s.cache.Set(identifier, []byte{})
			case <-ctx.Done():
				return
			}
		}
	}()

	// fetch from local store if available
	if _, found := s.cache.Get(identifier); found {
		lastNotes, _ = wdb.QuerySync(ctx, filter)
	} else {
		// if we didn't get enough notes (or if we didn't even query the local store), wait for the external relays
		lastNotes = <-external
		s.cache.Set(identifier, []byte{})

		// 		tags := nostr.Tags{
		// 			{"a", tag},
		// 		}
		//
		// 		e := nostr.Event{
		// 			Kind:      9802,
		// 			PubKey:    pubkey,
		// 			Content:   "",
		// 			CreatedAt: nostr.Now(),
		// 			Tags:      tags,
		// 		}
		//
		// 		// USe the server secret key, makes it easy to filer using pubkey
		// 		sk := os.Getenv("NOSTR_SK")
		// 		_ = e.Sign(sk)
		//
		// 		s.db.SaveEvent(ctx, &e)
		//
		// 		// Add a dummy
		// 		lastNotes = append(lastNotes, &e)
	}

	return lastNotes, nil
}

func (s *eventService) queryRelays(ctx context.Context, filter nostr.Filter) (ev []*nostr.Event) {

	var m sync.Map
	var wg sync.WaitGroup
	for _, url := range s.relays {

		wg.Add(1)
		go func(wg *sync.WaitGroup, url string) {
			defer wg.Done()

			r, err := nostr.RelayConnect(ctx, url)
			if err != nil {
				panic(err)
			}

			events, err := r.QuerySync(ctx, filter)
			if err != nil {
				// TODO
				return
			}

			for _, e := range events {
				m.Store(e.ID, e)
			}

		}(&wg, url)
	}
	wg.Wait()

	m.Range(func(_, v any) bool {
		ev = append(ev, v.(*nostr.Event))
		return true
	})

	return ev
}
