package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/firefly-software-mt/standard-template/internal/database"
	"github.com/firefly-software-mt/standard-template/internal/mail"
	"github.com/firefly-software-mt/standard-template/internal/view"
)

// tripLabels maps form values to human-readable labels for the email.
var tripLabels = map[string]string{
	"full_day_single": "Full Day — Single Boat ($675)",
	"half_day_single": "Half Day — Single Boat ($575)",
	"early_season":    "Early Season Full Day ($500)",
	"multiple_boats":  "Full Day — Multiple Boats ($675/boat)",
	"heroes":          "Heroes Rate — Full Day ($500)",
}

var experienceLabels = map[string]string{
	"never":       "Never fished — first time",
	"some":        "Some experience",
	"comfortable": "Comfortable — fishes regularly",
	"advanced":    "Advanced angler",
}

// Contact handles GET /contact and renders the booking form.
func Contact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := make(map[string]string)

		// Pre-select trip type from query param (e.g. /contact?rate=heroes)
		if rate := r.URL.Query().Get("rate"); rate == "heroes" {
			values["trip_type"] = "heroes"
		}

		if err := view.ContactPage(nil, values, false).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// ContactSubmit handles POST /contact, validates input, saves to database, and sends booking email.
func ContactSubmit(mailer *mail.Client, turnstileSecret string, db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		values := map[string]string{
			"trip_type":      strings.TrimSpace(r.FormValue("trip_type")),
			"preferred_date": strings.TrimSpace(r.FormValue("preferred_date")),
			"alternate_date": strings.TrimSpace(r.FormValue("alternate_date")),
			"angler_count":   strings.TrimSpace(r.FormValue("angler_count")),
			"youth_count":    strings.TrimSpace(r.FormValue("youth_count")),
			"experience":     strings.TrimSpace(r.FormValue("experience")),
			"referred_by":    strings.TrimSpace(r.FormValue("referred_by")),
			"client_notes":   strings.TrimSpace(r.FormValue("client_notes")),
			"client_name":    strings.TrimSpace(r.FormValue("client_name")),
			"client_email":   strings.TrimSpace(r.FormValue("client_email")),
			"client_phone":   strings.TrimSpace(r.FormValue("client_phone")),
		}

		errors := validateBooking(values)

		if len(errors) > 0 {
			if err := view.BookingForm(errors, values, false).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		// Verify Turnstile token
		if turnstileSecret != "" {
			token := r.FormValue("cf-turnstile-response")
			if !verifyTurnstile(turnstileSecret, token, r.RemoteAddr) {
				errors = map[string]string{"form": "Verification failed. Please try again."}
				if err := view.BookingForm(errors, values, false).Render(r.Context(), w); err != nil {
					slog.Error("render error", "err", err)
				}
				return
			}
		}

		// Save booking to database
		booking := database.BookingRequest{
			IPAddress:     r.RemoteAddr,
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
			slog.Error("database insert error", "err", err)
			errors = map[string]string{"form": "Failed to save your request. Please try again or call us directly."}
			if err := view.BookingForm(errors, values, false).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		if err := db.InsertEvent(bookingID, "submitted", "Booking request submitted via website"); err != nil {
			slog.Error("event insert error", "err", err)
		}

		// Send notification email
		if mailer != nil {
			msg := mail.Message{
				Name:    values["client_name"],
				Email:   values["client_email"],
				Subject: fmt.Sprintf("New booking request — %s · %s · %s", tripLabel(values["trip_type"]), values["preferred_date"], values["client_name"]),
				Body:    formatBookingEmail(values),
			}
			if err := mailer.Send(msg); err != nil {
				slog.Error("postmark send error", "err", err)
				if err := db.InsertEvent(bookingID, "email_failed", err.Error()); err != nil {
					slog.Error("event insert error", "err", err)
				}
			} else {
				if err := db.SetEmailedAt(bookingID, time.Now()); err != nil {
					slog.Error("set emailed_at error", "err", err)
				}
				if err := db.InsertEvent(bookingID, "email_sent", "Notification email sent"); err != nil {
					slog.Error("event insert error", "err", err)
				}
			}
		}

		if err := view.BookingForm(nil, nil, true).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

func validateBooking(v map[string]string) map[string]string {
	errors := make(map[string]string)

	if v["trip_type"] == "" {
		errors["trip_type"] = "Please select a trip type."
	} else if _, ok := tripLabels[v["trip_type"]]; !ok {
		errors["trip_type"] = "Please select a valid trip type."
	}

	if v["preferred_date"] == "" {
		errors["preferred_date"] = "Preferred date is required."
	} else if d, err := time.Parse("2006-01-02", v["preferred_date"]); err == nil {
		if d.Before(time.Now().Truncate(24 * time.Hour)) {
			errors["preferred_date"] = "Please select a future date."
		}
	}

	if v["angler_count"] == "" {
		errors["angler_count"] = "Number of anglers is required."
	}

	if v["experience"] == "" {
		errors["experience"] = "Please select your experience level."
	}

	if v["client_name"] == "" {
		errors["client_name"] = "Name is required."
	}

	if v["client_email"] == "" {
		errors["client_email"] = "Email is required."
	} else if !strings.Contains(v["client_email"], "@") {
		errors["client_email"] = "Enter a valid email address."
	}

	if v["client_phone"] == "" {
		errors["client_phone"] = "Phone number is required."
	}

	return errors
}

func tripLabel(val string) string {
	if label, ok := tripLabels[val]; ok {
		return label
	}
	return val
}

func formatBookingEmail(v map[string]string) string {
	var b strings.Builder

	b.WriteString("BOOKING REQUEST\n")
	b.WriteString("===============\n\n")

	b.WriteString("TRIP\n")
	b.WriteString(fmt.Sprintf("  Type:           %s\n", tripLabel(v["trip_type"])))
	b.WriteString(fmt.Sprintf("  Preferred date: %s\n", v["preferred_date"]))
	if v["alternate_date"] != "" {
		b.WriteString(fmt.Sprintf("  Alternate date: %s\n", v["alternate_date"]))
	}
	b.WriteString(fmt.Sprintf("  Anglers:        %s\n", v["angler_count"]))
	if v["youth_count"] != "" && v["youth_count"] != "0" {
		b.WriteString(fmt.Sprintf("  Youth (under 16): %s\n", v["youth_count"]))
	}
	if v["trip_type"] == "heroes" {
		b.WriteString("  ** Heroes rate requested **\n")
	}

	b.WriteString("\nPARTY\n")
	if label, ok := experienceLabels[v["experience"]]; ok {
		b.WriteString(fmt.Sprintf("  Experience:     %s\n", label))
	}
	if v["referred_by"] != "" {
		b.WriteString(fmt.Sprintf("  Referred by:    %s\n", v["referred_by"]))
	}
	if v["client_notes"] != "" {
		b.WriteString(fmt.Sprintf("\n  Notes:\n  %s\n", v["client_notes"]))
	}

	b.WriteString("\nCONTACT\n")
	b.WriteString(fmt.Sprintf("  Name:  %s\n", v["client_name"]))
	b.WriteString(fmt.Sprintf("  Email: %s\n", v["client_email"]))
	b.WriteString(fmt.Sprintf("  Phone: %s\n", v["client_phone"]))

	return b.String()
}

// verifyTurnstile checks a Turnstile token against the Cloudflare API.
func verifyTurnstile(secret, token, remoteIP string) bool {
	if token == "" {
		slog.Warn("turnstile token is empty — widget may not have loaded or HTMX did not include it")
		return false
	}

	// Note: remoteip is omitted — behind a reverse proxy r.RemoteAddr is the
	// proxy's IP, which can cause Cloudflare to reject the token.
	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{
		"secret":   {secret},
		"response": {token},
	})
	if err != nil {
		slog.Error("turnstile verify request failed", "err", err)
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Success    bool     `json:"success"`
		ErrorCodes []string `json:"error-codes"`
		Hostname   string   `json:"hostname"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("turnstile verify decode failed", "err", err)
		return false
	}

	if !result.Success {
		slog.Warn("turnstile verification failed",
			"error_codes", result.ErrorCodes,
			"hostname", result.Hostname,
			"remote_ip", remoteIP,
		)
	}
	return result.Success
}
