package handler

import (
	"log/slog"
	"net/http"

	"github.com/firefly-software-mt/standard-template/internal/view"
)

// Guides handles GET /guides and renders the guides page.
func Guides() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.GuidesPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
