package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

// AdminSession represents a row in the admin_sessions table.
type AdminSession struct {
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewSessionToken generates a 32-byte random hex token.
func NewSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CreateSession inserts a new session token with the given expiry.
func (db *DB) CreateSession(token string, expiresAt time.Time) error {
	_, err := db.Exec(
		"INSERT INTO admin_sessions (token, expires_at) VALUES (?, ?)",
		token, expiresAt,
	)
	return err
}

// GetSession retrieves a session by token. Returns nil if not found or expired.
func (db *DB) GetSession(token string) (*AdminSession, error) {
	s := &AdminSession{}
	err := db.QueryRow(
		"SELECT token, created_at, expires_at FROM admin_sessions WHERE token = ? AND expires_at > ?",
		token, time.Now(),
	).Scan(&s.Token, &s.CreatedAt, &s.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteSession removes a session by token.
func (db *DB) DeleteSession(token string) error {
	_, err := db.Exec("DELETE FROM admin_sessions WHERE token = ?", token)
	return err
}

// CleanExpiredSessions removes all expired sessions. Returns the number removed.
func (db *DB) CleanExpiredSessions() (int64, error) {
	result, err := db.Exec("DELETE FROM admin_sessions WHERE expires_at <= ?", time.Now())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
