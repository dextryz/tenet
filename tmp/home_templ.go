// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.590
package tmp

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Index() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><link rel=\"stylesheet\" href=\"https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css\"><link href=\"https://fonts.googleapis.com/css2?family=Fira+Code&amp;display=swap\" rel=\"stylesheet\"><link rel=\"stylesheet\" href=\"/static/style.css\" type=\"text/css\"><script src=\"https://unpkg.com/htmx.org@1.9.2\"></script></head><body hx-boost=\"true\"><main><form class=\"search-container\" hx-get=\"/list\" hx-target=\"#cards\" hx-indicator=\"#spinner\" hx-swap=\"innerHTML\"><input class=\"search-bar\" name=\"search\" type=\"search\" placeholder=\"Paste naddr or npub and Enter\"></form><div id=\"cards\"><div id=\"spinner\" class=\"loading-container\"><img class=\"htmx-indicator\" src=\"/img/bars.svg\"><p>Requesting Notes...</p></div></div><footer><p>Made with <i class=\"fas fa-heart\"></i> by <a href=\"https://github.com/dextryz\">dextryz</a></p></footer></main></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
