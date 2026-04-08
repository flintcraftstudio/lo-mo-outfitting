package handler

import (
	"log/slog"
	"net/http"

	"github.com/firefly-software-mt/standard-template/internal/view"
)

// About handles GET /about and renders the about page.
func About() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.AboutPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
