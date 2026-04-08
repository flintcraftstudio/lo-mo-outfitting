package database

import (
	"fmt"
	"strings"
	"time"
)

// BookingListItem is a lightweight struct for list display.
type BookingListItem struct {
	ID              int64
	CreatedAt       time.Time
	TripType        string
	PreferredDate   string
	AnglerCount     string
	YouthCount      string
	Heroes          bool
	Experience      string
	Lodging         string
	ClientName      string
	Status          string
	GuideName       string
	PaymentMethod   *string
	StatusUpdatedAt *time.Time
}

// CountByStatus returns a map of status → count for all booking requests.
func (db *DB) CountByStatus() (map[string]int, error) {
	rows, err := db.Query("SELECT status, COUNT(*) FROM booking_requests GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, rows.Err()
}

// ListOpenBookings returns booking list items filtered by status.
// If statusFilter is empty, returns all open (excludes complete and cancelled).
// "new" requests are floated to the top, then sorted by preferred_date ASC.
func (db *DB) ListOpenBookings(statusFilter string) ([]BookingListItem, error) {
	var where string
	var args []interface{}

	switch statusFilter {
	case "new", "contacted", "deposit_sent", "confirmed", "complete", "cancelled":
		where = "WHERE br.status = ?"
		args = append(args, statusFilter)
	default:
		// "all open" — exclude complete and cancelled
		where = "WHERE br.status NOT IN ('complete', 'cancelled')"
	}

	query := fmt.Sprintf(`
		SELECT br.id, br.created_at, br.trip_type, br.preferred_date, br.angler_count,
		       br.youth_count, br.heroes, br.experience, br.client_name, br.status,
		       COALESCE(g.name, ''), br.payment_method, br.status_updated_at
		FROM booking_requests br
		LEFT JOIN guides g ON br.guide_id = g.id
		%s
		ORDER BY
			CASE WHEN br.status = 'new' THEN 0 ELSE 1 END,
			br.preferred_date ASC
	`, where)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []BookingListItem
	for rows.Next() {
		var item BookingListItem
		var heroes int
		if err := rows.Scan(
			&item.ID, &item.CreatedAt, &item.TripType, &item.PreferredDate,
			&item.AnglerCount, &item.YouthCount, &heroes, &item.Experience,
			&item.ClientName, &item.Status, &item.GuideName,
			&item.PaymentMethod, &item.StatusUpdatedAt,
		); err != nil {
			return nil, err
		}
		item.Heroes = heroes == 1
		items = append(items, item)
	}
	return items, rows.Err()
}

// GetBookingDetail returns a full booking request by ID.
func (db *DB) GetBookingDetail(id int64) (*BookingRequest, error) {
	var b BookingRequest
	var heroes int
	err := db.QueryRow(`
		SELECT id, created_at, emailed_at, ip_address,
		       trip_type, preferred_date, COALESCE(alternate_date, ''), angler_count,
		       youth_count, heroes, experience, lodging, COALESCE(lodging_other, ''),
		       COALESCE(client_notes, ''), COALESCE(referred_by, ''),
		       client_name, client_email, client_phone,
		       status, guide_id, payment_method, COALESCE(mat_notes, ''), status_updated_at
		FROM booking_requests WHERE id = ?
	`, id).Scan(
		&b.ID, &b.CreatedAt, &b.EmailedAt, &b.IPAddress,
		&b.TripType, &b.PreferredDate, &b.AlternateDate, &b.AnglerCount,
		&b.YouthCount, &heroes, &b.Experience, &b.Lodging, &b.LodgingOther,
		&b.ClientNotes, &b.ReferredBy,
		&b.ClientName, &b.ClientEmail, &b.ClientPhone,
		&b.Status, &b.GuideID, &b.PaymentMethod, &b.MatNotes, &b.StatusUpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	b.Heroes = heroes == 1
	return &b, nil
}

// ListBookingEvents returns all events for a booking, oldest first.
func (db *DB) ListBookingEvents(bookingID int64) ([]BookingEvent, error) {
	rows, err := db.Query(`
		SELECT id, booking_request_id, created_at, event_type, COALESCE(detail, '')
		FROM booking_events
		WHERE booking_request_id = ?
		ORDER BY id ASC
	`, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []BookingEvent
	for rows.Next() {
		var e BookingEvent
		if err := rows.Scan(&e.ID, &e.BookingRequestID, &e.CreatedAt, &e.EventType, &e.Detail); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// UpdateStatus changes the status of a booking and records the transition.
func (db *DB) UpdateStatus(id int64, newStatus string) error {
	valid := map[string]bool{
		"new": true, "contacted": true, "deposit_sent": true,
		"confirmed": true, "complete": true, "cancelled": true,
	}
	if !valid[newStatus] {
		return fmt.Errorf("invalid status: %q", newStatus)
	}

	// Get current status for the event log
	var oldStatus string
	if err := db.QueryRow("SELECT status FROM booking_requests WHERE id = ?", id).Scan(&oldStatus); err != nil {
		return err
	}

	now := time.Now()
	if _, err := db.Exec(
		"UPDATE booking_requests SET status = ?, status_updated_at = ? WHERE id = ?",
		newStatus, now, id,
	); err != nil {
		return err
	}

	detail := fmt.Sprintf("Status changed: %s → %s", oldStatus, newStatus)
	return db.InsertEvent(id, "status_changed", detail)
}

// AssignGuide assigns a guide to a booking and logs the event.
func (db *DB) AssignGuide(id int64, guideID int64) error {
	guide, err := db.GetGuide(guideID)
	if err != nil {
		return fmt.Errorf("guide not found: %w", err)
	}

	if _, err := db.Exec(
		"UPDATE booking_requests SET guide_id = ? WHERE id = ?",
		guideID, id,
	); err != nil {
		return err
	}

	detail := fmt.Sprintf("Guide assigned: %s", guide.Name)
	return db.InsertEvent(id, "guide_assigned", detail)
}

// AddNote saves Matt's note and logs the event.
func (db *DB) AddNote(id int64, note string) error {
	note = strings.TrimSpace(note)
	if note == "" {
		return nil
	}

	if _, err := db.Exec(
		"UPDATE booking_requests SET mat_notes = ? WHERE id = ?",
		note, id,
	); err != nil {
		return err
	}

	detail := fmt.Sprintf("Note: %s", note)
	return db.InsertEvent(id, "note_added", detail)
}

// SetPaymentMethod sets the payment method and logs the event.
func (db *DB) SetPaymentMethod(id int64, method string) error {
	valid := map[string]bool{"cash": true, "venmo": true, "stripe": true, "other": true}
	if !valid[method] {
		return fmt.Errorf("invalid payment method: %q", method)
	}

	if _, err := db.Exec(
		"UPDATE booking_requests SET payment_method = ? WHERE id = ?",
		method, id,
	); err != nil {
		return err
	}

	detail := fmt.Sprintf("Payment method set: %s", method)
	return db.InsertEvent(id, "payment_method_set", detail)
}

// CheckGuideConflict returns the count of confirmed bookings for a guide on a given date,
// excluding the specified booking ID.
func (db *DB) CheckGuideConflict(guideID int64, date string, excludeBookingID int64) (int, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM booking_requests
		WHERE guide_id = ?
		AND preferred_date = ?
		AND status = 'confirmed'
		AND id != ?
	`, guideID, date, excludeBookingID).Scan(&count)
	return count, err
}

// UpcomingSummary holds aggregate metrics for the upcoming trips view.
type UpcomingSummary struct {
	TripCount          int
	GuidesScheduled    int
	TotalAnglers       int
	DepositOutstanding int
}

// ListUpcomingConfirmed returns confirmed bookings within the next windowDays, sorted by preferred_date ASC.
func (db *DB) ListUpcomingConfirmed(windowDays int) ([]BookingListItem, error) {
	today := time.Now().Format("2006-01-02")
	endDate := time.Now().AddDate(0, 0, windowDays).Format("2006-01-02")

	rows, err := db.Query(`
		SELECT br.id, br.created_at, br.trip_type, br.preferred_date, br.angler_count,
		       br.youth_count, br.heroes, br.experience, br.client_name, br.status,
		       COALESCE(g.name, ''), br.payment_method, br.status_updated_at,
		       br.lodging, COALESCE(br.lodging_other, '')
		FROM booking_requests br
		LEFT JOIN guides g ON br.guide_id = g.id
		WHERE br.status = 'confirmed'
		AND br.preferred_date >= ?
		AND br.preferred_date <= ?
		ORDER BY br.preferred_date ASC
	`, today, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []BookingListItem
	for rows.Next() {
		var item BookingListItem
		var heroes int
		var lodging, lodgingOther string
		if err := rows.Scan(
			&item.ID, &item.CreatedAt, &item.TripType, &item.PreferredDate,
			&item.AnglerCount, &item.YouthCount, &heroes, &item.Experience,
			&item.ClientName, &item.Status, &item.GuideName,
			&item.PaymentMethod, &item.StatusUpdatedAt,
			&lodging, &lodgingOther,
		); err != nil {
			return nil, err
		}
		item.Heroes = heroes == 1
		if lodging == "other" && lodgingOther != "" {
			item.Lodging = lodgingOther
		} else {
			item.Lodging = lodging
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// GetUpcomingSummary returns aggregate metrics for confirmed trips in the next windowDays.
func (db *DB) GetUpcomingSummary(windowDays int) (*UpcomingSummary, error) {
	today := time.Now().Format("2006-01-02")
	endDate := time.Now().AddDate(0, 0, windowDays).Format("2006-01-02")

	s := &UpcomingSummary{}

	err := db.QueryRow(`
		SELECT
			COUNT(*),
			COUNT(DISTINCT guide_id),
			COALESCE(SUM(CAST(
				CASE WHEN angler_count GLOB '[0-9]*' THEN angler_count ELSE '0' END
			AS INTEGER)), 0),
			SUM(CASE WHEN payment_method IS NULL THEN 1 ELSE 0 END)
		FROM booking_requests
		WHERE status = 'confirmed'
		AND preferred_date >= ?
		AND preferred_date <= ?
	`, today, endDate).Scan(&s.TripCount, &s.GuidesScheduled, &s.TotalAnglers, &s.DepositOutstanding)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// SearchBookings returns paginated booking list items matching a search query.
// Searches client_name, client_email, and client_phone. Returns items and total count.
func (db *DB) SearchBookings(query string, page, perPage int) ([]BookingListItem, int, error) {
	query = strings.TrimSpace(query)
	offset := (page - 1) * perPage

	var where string
	var args []interface{}

	if query != "" {
		where = `WHERE (br.client_name LIKE ? OR br.client_email LIKE ? OR br.client_phone LIKE ?)`
		like := "%" + query + "%"
		args = append(args, like, like, like)
	}

	// Get total count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM booking_requests br %s", where)
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get page of results
	dataQuery := fmt.Sprintf(`
		SELECT br.id, br.created_at, br.trip_type, br.preferred_date, br.angler_count,
		       br.youth_count, br.heroes, br.experience, br.client_name, br.status,
		       COALESCE(g.name, ''), br.payment_method, br.status_updated_at
		FROM booking_requests br
		LEFT JOIN guides g ON br.guide_id = g.id
		%s
		ORDER BY br.created_at DESC
		LIMIT ? OFFSET ?
	`, where)

	dataArgs := append(args, perPage, offset)
	rows, err := db.Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []BookingListItem
	for rows.Next() {
		var item BookingListItem
		var heroes int
		if err := rows.Scan(
			&item.ID, &item.CreatedAt, &item.TripType, &item.PreferredDate,
			&item.AnglerCount, &item.YouthCount, &heroes, &item.Experience,
			&item.ClientName, &item.Status, &item.GuideName,
			&item.PaymentMethod, &item.StatusUpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		item.Heroes = heroes == 1
		items = append(items, item)
	}
	return items, total, rows.Err()
}
