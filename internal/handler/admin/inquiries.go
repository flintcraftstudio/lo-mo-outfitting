package admin

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/firefly-software-mt/standard-template/internal/database"
	adminview "github.com/firefly-software-mt/standard-template/internal/view/admin"
)

// Inquiries renders the inquiry board page.
func Inquiries(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("status")

		counts, err := db.CountByStatus()
		if err != nil {
			slog.Error("count by status", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		bookings, err := db.ListOpenBookings(filter)
		if err != nil {
			slog.Error("list bookings", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if err := adminview.InquiriesPage(counts, bookings, filter).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// InquiryDetail renders the detail pane for a single inquiry (htmx partial).
func InquiryDetail(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		booking, err := db.GetBookingDetail(id)
		if err != nil {
			slog.Error("get booking detail", "err", err, "id", id)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		events, err := db.ListBookingEvents(id)
		if err != nil {
			slog.Error("list booking events", "err", err, "id", id)
		}

		guides, err := db.ListActiveGuides()
		if err != nil {
			slog.Error("list guides", "err", err)
		}

		// Get assigned guide name if set
		var guideName string
		if booking.GuideID != nil {
			if guide, err := db.GetGuide(*booking.GuideID); err == nil {
				guideName = guide.Name
			}
		}

		if err := adminview.InquiryDetailPane(booking, events, guides, guideName).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// StatusUpdate handles POST /admin/inquiries/{id}/status.
func StatusUpdate(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		newStatus := r.FormValue("status")
		if err := db.UpdateStatus(id, newStatus); err != nil {
			slog.Error("update status", "err", err, "id", id)
			http.Error(w, "failed to update status", http.StatusInternalServerError)
			return
		}

		renderDetailWithOOB(w, r, db, id)
	}
}

// GuideAssign handles POST /admin/inquiries/{id}/guide.
func GuideAssign(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		guideID, err := strconv.ParseInt(r.FormValue("guide_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid guide_id", http.StatusBadRequest)
			return
		}

		if err := db.AssignGuide(id, guideID); err != nil {
			slog.Error("assign guide", "err", err, "id", id)
			http.Error(w, "failed to assign guide", http.StatusInternalServerError)
			return
		}

		renderDetailWithOOB(w, r, db, id)
	}
}

// NoteAdd handles POST /admin/inquiries/{id}/note.
func NoteAdd(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		note := r.FormValue("note")
		if err := db.AddNote(id, note); err != nil {
			slog.Error("add note", "err", err, "id", id)
			http.Error(w, "failed to add note", http.StatusInternalServerError)
			return
		}

		renderDetailWithOOB(w, r, db, id)
	}
}

// PaymentMethod handles POST /admin/inquiries/{id}/payment.
func PaymentMethod(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		method := r.FormValue("payment_method")
		if err := db.SetPaymentMethod(id, method); err != nil {
			slog.Error("set payment method", "err", err, "id", id)
			http.Error(w, "failed to set payment method", http.StatusInternalServerError)
			return
		}

		renderDetailWithOOB(w, r, db, id)
	}
}

// renderDetailWithOOB re-renders the detail pane and includes OOB-swapped status pills.
func renderDetailWithOOB(w http.ResponseWriter, r *http.Request, db *database.DB, id int64) {
	booking, err := db.GetBookingDetail(id)
	if err != nil {
		slog.Error("get booking detail", "err", err, "id", id)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	events, err := db.ListBookingEvents(id)
	if err != nil {
		slog.Error("list booking events", "err", err)
	}

	guides, err := db.ListActiveGuides()
	if err != nil {
		slog.Error("list guides", "err", err)
	}

	counts, err := db.CountByStatus()
	if err != nil {
		slog.Error("count by status", "err", err)
	}

	var guideName string
	if booking.GuideID != nil {
		if guide, err := db.GetGuide(*booking.GuideID); err == nil {
			guideName = guide.Name
		}
	}

	// Render detail pane + OOB status pills
	if err := adminview.InquiryDetailPane(booking, events, guides, guideName).Render(r.Context(), w); err != nil {
		slog.Error("render error", "err", err)
		return
	}
	if err := adminview.StatusPillsOOB(counts).Render(r.Context(), w); err != nil {
		slog.Error("render oob error", "err", err)
	}
}
