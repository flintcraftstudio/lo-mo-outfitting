package database

import (
	"testing"
	"time"
)

func mustOpen(t *testing.T) *DB {
	t.Helper()
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open(:memory:) failed: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestOpenAndMigrate(t *testing.T) {
	db := mustOpen(t)

	// Verify tables exist by querying them
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM booking_requests").Scan(&count); err != nil {
		t.Fatalf("booking_requests table missing: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM guides").Scan(&count); err != nil {
		t.Fatalf("guides table missing: %v", err)
	}
	if count != 7 {
		t.Errorf("expected 7 seeded guides, got %d", count)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM booking_events").Scan(&count); err != nil {
		t.Fatalf("booking_events table missing: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM admin_sessions").Scan(&count); err != nil {
		t.Fatalf("admin_sessions table missing: %v", err)
	}

	// Verify schema_version
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_version").Scan(&count); err != nil {
		t.Fatalf("schema_version table missing: %v", err)
	}
	if count != 4 {
		t.Errorf("expected 4 schema versions, got %d", count)
	}
}

func TestMigrateIdempotent(t *testing.T) {
	db := mustOpen(t)

	// Running migrate again should be a no-op
	if err := db.migrate(); err != nil {
		t.Fatalf("second migrate failed: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_version").Scan(&count); err != nil {
		t.Fatalf("schema_version query failed: %v", err)
	}
	if count != 4 {
		t.Errorf("expected 4 schema versions after re-migrate, got %d", count)
	}
}

func TestInsertBooking(t *testing.T) {
	db := mustOpen(t)

	b := &BookingRequest{
		IPAddress:     "127.0.0.1",
		TripType:      "full_day_single",
		PreferredDate: "2026-07-15",
		AlternateDate: "2026-07-16",
		AnglerCount:   "2",
		YouthCount:    "0",
		Heroes:        false,
		Experience:    "comfortable",
		Lodging:       "craig",
		ClientName:    "John Doe",
		ClientEmail:   "john@example.com",
		ClientPhone:   "555-1234",
		ClientNotes:   "Looking forward to it",
		ReferredBy:    "A friend",
	}

	id, err := db.InsertBooking(b)
	if err != nil {
		t.Fatalf("InsertBooking failed: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero id")
	}

	// Verify the row
	var name, status string
	if err := db.QueryRow("SELECT client_name, status FROM booking_requests WHERE id = ?", id).Scan(&name, &status); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if name != "John Doe" {
		t.Errorf("expected name 'John Doe', got %q", name)
	}
	if status != "new" {
		t.Errorf("expected status 'new', got %q", status)
	}
}

func TestInsertBookingHeroes(t *testing.T) {
	db := mustOpen(t)

	b := &BookingRequest{
		TripType:      "heroes",
		PreferredDate: "2026-07-15",
		AnglerCount:   "1",
		Heroes:        true,
		Experience:    "some",
		Lodging:       "helena",
		ClientName:    "Jane Doe",
		ClientEmail:   "jane@example.com",
		ClientPhone:   "555-5678",
	}

	id, err := db.InsertBooking(b)
	if err != nil {
		t.Fatalf("InsertBooking failed: %v", err)
	}

	var heroes int
	if err := db.QueryRow("SELECT heroes FROM booking_requests WHERE id = ?", id).Scan(&heroes); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if heroes != 1 {
		t.Errorf("expected heroes=1, got %d", heroes)
	}
}

func TestSetEmailedAt(t *testing.T) {
	db := mustOpen(t)

	b := &BookingRequest{
		TripType:      "half_day_single",
		PreferredDate: "2026-08-01",
		AnglerCount:   "1",
		Experience:    "never",
		Lodging:       "not_sure",
		ClientName:    "Test User",
		ClientEmail:   "test@example.com",
		ClientPhone:   "555-0000",
	}

	id, err := db.InsertBooking(b)
	if err != nil {
		t.Fatalf("InsertBooking failed: %v", err)
	}

	now := time.Now()
	if err := db.SetEmailedAt(id, now); err != nil {
		t.Fatalf("SetEmailedAt failed: %v", err)
	}

	var emailedAt *string
	if err := db.QueryRow("SELECT emailed_at FROM booking_requests WHERE id = ?", id).Scan(&emailedAt); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if emailedAt == nil {
		t.Error("expected emailed_at to be set, got nil")
	}
}

func TestInsertEvent(t *testing.T) {
	db := mustOpen(t)

	b := &BookingRequest{
		TripType:      "full_day_single",
		PreferredDate: "2026-09-01",
		AnglerCount:   "3",
		Experience:    "advanced",
		Lodging:       "craig",
		ClientName:    "Event Test",
		ClientEmail:   "event@example.com",
		ClientPhone:   "555-9999",
	}

	id, err := db.InsertBooking(b)
	if err != nil {
		t.Fatalf("InsertBooking failed: %v", err)
	}

	if err := db.InsertEvent(id, "submitted", "Booking request submitted via website"); err != nil {
		t.Fatalf("InsertEvent failed: %v", err)
	}
	if err := db.InsertEvent(id, "email_sent", "Notification email sent"); err != nil {
		t.Fatalf("InsertEvent failed: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM booking_events WHERE booking_request_id = ?", id).Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 events, got %d", count)
	}

	// Verify event details
	var eventType, detail string
	if err := db.QueryRow("SELECT event_type, detail FROM booking_events WHERE booking_request_id = ? ORDER BY id LIMIT 1", id).Scan(&eventType, &detail); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if eventType != "submitted" {
		t.Errorf("expected event_type 'submitted', got %q", eventType)
	}
	if detail != "Booking request submitted via website" {
		t.Errorf("unexpected detail: %q", detail)
	}
}
