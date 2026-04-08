package admin

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/firefly-software-mt/standard-template/internal/database"
	adminview "github.com/firefly-software-mt/standard-template/internal/view/admin"
)

const perPage = 25

// AllBookings renders the searchable, paginated all-bookings table.
func AllBookings(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		bookings, total, err := db.SearchBookings(query, page, perPage)
		if err != nil {
			slog.Error("search bookings", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		totalPages := (total + perPage - 1) / perPage
		if totalPages < 1 {
			totalPages = 1
		}

		// If htmx request, render just the table body + pagination
		if r.Header.Get("HX-Request") == "true" {
			if err := adminview.AllBookingsTableBody(bookings, query, page, totalPages).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		if err := adminview.AllBookingsPage(bookings, query, page, totalPages).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
