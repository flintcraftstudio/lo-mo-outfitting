package handler

import (
	"log/slog"
	"net/http"

	"github.com/firefly-software-mt/standard-template/internal/view"
)

// Store handles GET /store and renders the store stub page.
func Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.StorePage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
