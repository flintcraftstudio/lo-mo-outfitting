package database

import (
	"time"
)

// BookingRequest represents a row in the booking_requests table.
type BookingRequest struct {
	ID              int64
	CreatedAt       time.Time
	EmailedAt       *time.Time
	IPAddress       string
	TripType        string
	PreferredDate   string
	AlternateDate   string
	AnglerCount     string
	YouthCount      string
	Heroes          bool
	Experience      string
	Lodging         string
	LodgingOther    string
	ClientNotes     string
	ReferredBy      string
	ClientName      string
	ClientEmail     string
	ClientPhone     string
	Status          string
	GuideID         *int64
	PaymentMethod   *string
	MatNotes        string
	StatusUpdatedAt *time.Time
}

// BookingEvent represents a row in the booking_events table.
type BookingEvent struct {
	ID               int64
	BookingRequestID int64
	CreatedAt        time.Time
	EventType        string
	Detail           string
}

// InsertBooking saves a new booking request and returns the new row ID.
func (db *DB) InsertBooking(b *BookingRequest) (int64, error) {
	heroes := 0
	if b.Heroes {
		heroes = 1
	}

	result, err := db.Exec(`
		INSERT INTO booking_requests (
			ip_address, trip_type, preferred_date, alternate_date,
			angler_count, youth_count, heroes,
			experience, lodging, lodging_other, client_notes, referred_by,
			client_name, client_email, client_phone
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		b.IPAddress, b.TripType, b.PreferredDate, b.AlternateDate,
		b.AnglerCount, b.YouthCount, heroes,
		b.Experience, b.Lodging, b.LodgingOther, b.ClientNotes, b.ReferredBy,
		b.ClientName, b.ClientEmail, b.ClientPhone,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// SetEmailedAt updates the emailed_at timestamp for a booking.
func (db *DB) SetEmailedAt(id int64, t time.Time) error {
	_, err := db.Exec("UPDATE booking_requests SET emailed_at = ? WHERE id = ?", t, id)
	return err
}

// InsertEvent appends a booking_events row.
func (db *DB) InsertEvent(bookingID int64, eventType, detail string) error {
	_, err := db.Exec(
		"INSERT INTO booking_events (booking_request_id, event_type, detail) VALUES (?, ?, ?)",
		bookingID, eventType, detail,
	)
	return err
}
