package admin

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/firefly-software-mt/standard-template/internal/database"
	adminview "github.com/firefly-software-mt/standard-template/internal/view/admin"
)

var manualTripTypes = map[string]bool{
	"full_day_single": true,
	"half_day_single": true,
	"early_season":    true,
	"multiple_boats":  true,
	"heroes":          true,
}

var manualSourceLabels = map[string]string{
	"phone":     "Phone",
	"in_person": "In person",
	"email":     "Email",
	"referral":  "Referral",
	"other":     "Other",
}

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

// NewBookingForm renders the manual booking entry form.
func NewBookingForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := map[string]string{"source": "phone"}
		if err := adminview.NewBookingPage(values, nil).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// NewBookingSubmit validates and inserts a manually entered booking.
func NewBookingSubmit(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		values := map[string]string{
			"source":         strings.TrimSpace(r.FormValue("source")),
			"client_name":    strings.TrimSpace(r.FormValue("client_name")),
			"client_email":   strings.TrimSpace(r.FormValue("client_email")),
			"client_phone":   strings.TrimSpace(r.FormValue("client_phone")),
			"trip_type":      strings.TrimSpace(r.FormValue("trip_type")),
			"preferred_date": strings.TrimSpace(r.FormValue("preferred_date")),
			"alternate_date": strings.TrimSpace(r.FormValue("alternate_date")),
			"angler_count":   strings.TrimSpace(r.FormValue("angler_count")),
			"youth_count":    strings.TrimSpace(r.FormValue("youth_count")),
			"experience":     strings.TrimSpace(r.FormValue("experience")),
			"referred_by":    strings.TrimSpace(r.FormValue("referred_by")),
			"client_notes":   strings.TrimSpace(r.FormValue("client_notes")),
		}

		errors := validateManualBooking(values)
		if len(errors) > 0 {
			if err := adminview.NewBookingPage(values, errors).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		booking := database.BookingRequest{
			IPAddress:     "manual",
			TripType:      values["trip_type"],
			PreferredDate: values["preferred_date"],
			AlternateDate: values["alternate_date"],
			AnglerCount:   values["angler_count"],
			YouthCount:    values["youth_count"],
			Heroes:        values["trip_type"] == "heroes",
			Experience:    values["experience"],
			ClientNotes:   values["client_notes"],
			ReferredBy:    values["referred_by"],
			ClientName:    values["client_name"],
			ClientEmail:   values["client_email"],
			ClientPhone:   values["client_phone"],
		}

		bookingID, err := db.InsertBooking(&booking)
		if err != nil {
			slog.Error("manual booking insert", "err", err)
			errors = map[string]string{"form": "Failed to save booking. Please try again."}
			if err := adminview.NewBookingPage(values, errors).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		sourceLabel := manualSourceLabels[values["source"]]
		if sourceLabel == "" {
			sourceLabel = values["source"]
		}
		if err := db.InsertEvent(bookingID, "manual_created", fmt.Sprintf("Booking manually added (source: %s)", sourceLabel)); err != nil {
			slog.Error("event insert error", "err", err)
		}
		if err := db.UpdateStatus(bookingID, "contacted"); err != nil {
			slog.Error("manual booking status update", "err", err)
		}

		http.Redirect(w, r, fmt.Sprintf("/admin/inquiries?status=contacted#booking-%d", bookingID), http.StatusSeeOther)
	}
}

func validateManualBooking(v map[string]string) map[string]string {
	errors := make(map[string]string)

	if v["client_name"] == "" {
		errors["client_name"] = "Name is required."
	}
	if v["trip_type"] == "" {
		errors["trip_type"] = "Trip type is required."
	} else if !manualTripTypes[v["trip_type"]] {
		errors["trip_type"] = "Select a valid trip type."
	}
	if v["preferred_date"] == "" {
		errors["preferred_date"] = "Preferred date is required."
	} else if _, err := time.Parse("2006-01-02", v["preferred_date"]); err != nil {
		errors["preferred_date"] = "Use YYYY-MM-DD."
	}
	if v["angler_count"] == "" {
		errors["angler_count"] = "Angler count is required."
	} else if n, err := strconv.Atoi(v["angler_count"]); err != nil || n < 1 {
		errors["angler_count"] = "Must be a positive number."
	}
	if v["client_email"] != "" && !strings.Contains(v["client_email"], "@") {
		errors["client_email"] = "Enter a valid email address."
	}

	return errors
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
			writeToast(w, "save-error", "Note save failed — try again.")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		writeToast(w, "save-success", "Note saved")
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
			writeToast(w, "save-error", "Payment update failed — try again.")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		writeToast(w, "save-success", "Payment updated")
		renderDetailWithOOB(w, r, db, id)
	}
}

// writeToast sets an HX-Trigger header so the client-side toast fires after the swap.
// event is the custom event name (e.g. "save-success"); message is shown in the toast.
// Must be called before any response body is written.
func writeToast(w http.ResponseWriter, event, message string) {
	payload, err := json.Marshal(map[string]map[string]string{
		event: {"message": message},
	})
	if err != nil {
		return
	}
	w.Header().Set("HX-Trigger", string(payload))
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
