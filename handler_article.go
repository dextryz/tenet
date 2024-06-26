package notezero

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

// 1. Highlights are encoded into data.Notes
// 2. Process the highlights into the data.Content string
func (s *Handler) ContentHandler(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("naddr")
	fmt.Println("Contetn handler")

	data, err := s.requestData(r.Context(), code, true)
	if err != nil {
		s.log.Error("failed to get events", slog.Any("error", err))
		http.Error(w, "failed to get counts", http.StatusInternalServerError)
		return
	}

	s.log.Info("rendering content view", "author", data.Npub, "highlightCount", len(data.Notes))

	var component templ.Component

	switch data.TemplateId {
	case Article:
		component = ContentTemplate(ArticleParams{
			Event:   data.Event,
			Content: template.HTML(data.Content), // data.Content is converted from Md to Html in data service.
		})
	default:
		s.log.Error("unable to render template", "templateId", data.TemplateId)
		http.Error(w, "tried to render an unsupported template", 500)
		return
	}

	err = component.Render(r.Context(), w)
	if err != nil {
		s.log.Error("error rendering tmpl", "error", err.Error())
	}
}

// Poplated the data.Notes field with a list of requested notes based on the search field.
func (s *Handler) ArticleHandler(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("naddr")
	npub := r.PathValue("npub")

	s.log.Info("handler for article", "naddr", code, "npub", npub)

	data, err := s.requestData(r.Context(), code, false)
	if err != nil {
		s.log.Error("failed to get events", slog.Any("error", err))
		http.Error(w, "failed to get counts", http.StatusInternalServerError)
		return
	}

	s.log.Info("rendering article view", "author", data.Npub)

	var component templ.Component

	switch data.TemplateId {
	case Article:
		component = ArticleTemplate(ArticleParams{
			Event:   data.Event,
			Content: template.HTML(data.Content), // data.Content is converted from Md to Html in data service.
		})
	default:
		s.log.Error("unable to render template", "templateId", data.TemplateId)
		http.Error(w, "tried to render an unsupported template", 500)
		return
	}

	err = component.Render(r.Context(), w)
	if err != nil {
		s.log.Error("error rendering tmpl", "error", err.Error())
	}
}
