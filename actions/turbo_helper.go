package actions

import (
	"io"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo/render"
)

func isTurboFrame(request *http.Request, frameId string) bool {
	for _, acceptValue := range request.Header["Turbo-Frame"] {
		if strings.Contains(acceptValue, frameId) {
			return true
		}
	}
	return false
}

func createTurboWriter(template, action, target string) render.RendererFunc {
	return func(w io.Writer, d render.Data) error {
		d["action"] = action
		d["target"] = target
		r.HTML(template, "turbo/turbo_stream.plush.html").Render(w, d)
		return nil
	}
}

func createTurboPlain(template string) render.RendererFunc {
	return func(w io.Writer, d render.Data) error {
		r.HTML(template, "turbo/plain.plush.html").Render(w, d)
		return nil
	}
}
