package handler

import (
	"log/slog"
	"net/http"

	"github.com/firefly-software-mt/standard-template/internal/view"
)

// Privacy handles GET /privacy and renders the privacy policy page.
func Privacy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.PrivacyPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
