package handler

import (
	"log/slog"
	"net/http"

	"github.com/firefly-software-mt/standard-template/internal/view"
)

// Reviews handles GET /reviews and renders the reviews page.
func Reviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.ReviewsPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
