package handler

import (
	"log/slog"
	"net/http"

	"github.com/firefly-software-mt/standard-template/internal/view"
)

// Policies handles GET /policies and renders the policies & FAQ page.
func Policies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.PoliciesPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
