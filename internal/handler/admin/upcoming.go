package admin

import (
	"log/slog"
	"net/http"

	"github.com/firefly-software-mt/standard-template/internal/database"
	adminview "github.com/firefly-software-mt/standard-template/internal/view/admin"
)

const upcomingWindowDays = 60

// Upcoming renders the upcoming confirmed trips view.
func Upcoming(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		summary, err := db.GetUpcomingSummary(upcomingWindowDays)
		if err != nil {
			slog.Error("upcoming summary", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		bookings, err := db.ListUpcomingConfirmed(upcomingWindowDays)
		if err != nil {
			slog.Error("list upcoming", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if err := adminview.UpcomingPage(summary, bookings).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
