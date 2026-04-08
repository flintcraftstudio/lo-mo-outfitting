package database

import (
	"testing"
	"time"
)

func TestNewSessionToken(t *testing.T) {
	token, err := NewSessionToken()
	if err != nil {
		t.Fatalf("NewSessionToken failed: %v", err)
	}
	if len(token) != 64 { // 32 bytes = 64 hex chars
		t.Errorf("expected 64 char token, got %d", len(token))
	}

	// Ensure uniqueness
	token2, _ := NewSessionToken()
	if token == token2 {
		t.Error("two tokens should not be equal")
	}
}

func TestCreateAndGetSession(t *testing.T) {
	db := mustOpen(t)

	token := "test-session-token-abc123"
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	if err := db.CreateSession(token, expiresAt); err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	session, err := db.GetSession(token)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	if session == nil {
		t.Fatal("expected session, got nil")
	}
	if session.Token != token {
		t.Errorf("expected token %q, got %q", token, session.Token)
	}
}

func TestGetSessionExpired(t *testing.T) {
	db := mustOpen(t)

	token := "expired-session-token"
	expiresAt := time.Now().Add(-1 * time.Hour) // already expired

	if err := db.CreateSession(token, expiresAt); err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	session, err := db.GetSession(token)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	if session != nil {
		t.Error("expected nil for expired session")
	}
}

func TestGetSessionNotFound(t *testing.T) {
	db := mustOpen(t)

	session, err := db.GetSession("nonexistent-token")
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	if session != nil {
		t.Error("expected nil for missing session")
	}
}

func TestDeleteSession(t *testing.T) {
	db := mustOpen(t)

	token := "delete-me-token"
	if err := db.CreateSession(token, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	if err := db.DeleteSession(token); err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}

	session, err := db.GetSession(token)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	if session != nil {
		t.Error("expected nil after deletion")
	}
}

func TestCleanExpiredSessions(t *testing.T) {
	db := mustOpen(t)

	// Create one valid and two expired sessions
	db.CreateSession("valid-token", time.Now().Add(time.Hour))
	db.CreateSession("expired-1", time.Now().Add(-time.Hour))
	db.CreateSession("expired-2", time.Now().Add(-2*time.Hour))

	n, err := db.CleanExpiredSessions()
	if err != nil {
		t.Fatalf("CleanExpiredSessions failed: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 cleaned, got %d", n)
	}

	// Valid session should still exist
	session, _ := db.GetSession("valid-token")
	if session == nil {
		t.Error("valid session should still exist")
	}
}
