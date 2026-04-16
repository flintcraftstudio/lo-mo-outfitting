package meta

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// Client sends server-side events to Meta's Conversions API.
type Client struct {
	pixelID     string
	accessToken string
}

// NewClient creates a Meta CAPI client.
func NewClient(pixelID, accessToken string) *Client {
	return &Client{
		pixelID:     pixelID,
		accessToken: accessToken,
	}
}

// ContactEvent holds the data needed to send a Contact event.
type ContactEvent struct {
	EventID   string // shared with browser pixel for deduplication
	SourceURL string // page URL where the form was submitted
	Email     string
	Phone     string
	FirstName string
	LastName  string
	ClientIP  string
	UserAgent string
	FBC       string // _fbc cookie value
	FBP       string // _fbp cookie value
}

type capiRequest struct {
	Data []eventData `json:"data"`
}

type eventData struct {
	EventName  string   `json:"event_name"`
	EventTime  int64    `json:"event_time"`
	EventID    string   `json:"event_id,omitempty"`
	SourceURL  string   `json:"event_source_url"`
	ActionSrc  string   `json:"action_source"`
	UserData   userData `json:"user_data"`
}

type userData struct {
	Email     []string `json:"em,omitempty"`
	Phone     []string `json:"ph,omitempty"`
	FirstName []string `json:"fn,omitempty"`
	LastName  []string `json:"ln,omitempty"`
	ClientIP  string   `json:"client_ip_address,omitempty"`
	UserAgent string   `json:"client_user_agent,omitempty"`
	FBC       string   `json:"fbc,omitempty"`
	FBP       string   `json:"fbp,omitempty"`
}

// SendContact sends a Contact event to the Conversions API.
// This should be called in a goroutine so it doesn't block the response.
func (c *Client) SendContact(evt ContactEvent) {
	ud := userData{
		ClientIP:  evt.ClientIP,
		UserAgent: evt.UserAgent,
		FBC:       evt.FBC,
		FBP:       evt.FBP,
	}

	if v := hashNormalized(evt.Email); v != "" {
		ud.Email = []string{v}
	}
	if v := hashPhone(evt.Phone); v != "" {
		ud.Phone = []string{v}
	}

	// Split "client_name" into first/last
	if evt.FirstName != "" {
		ud.FirstName = []string{hashNormalized(evt.FirstName)}
	}
	if evt.LastName != "" {
		ud.LastName = []string{hashNormalized(evt.LastName)}
	}

	payload := capiRequest{
		Data: []eventData{{
			EventName: "Contact",
			EventTime: time.Now().Unix(),
			EventID:   evt.EventID,
			SourceURL: evt.SourceURL,
			ActionSrc: "website",
			UserData:  ud,
		}},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error("meta capi marshal error", "err", err)
		return
	}

	url := fmt.Sprintf("https://graph.facebook.com/v21.0/%s/events?access_token=%s", c.pixelID, c.accessToken)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		slog.Error("meta capi request failed", "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		slog.Error("meta capi error", "status", resp.StatusCode, "body", string(respBody))
		return
	}

	slog.Info("meta capi Contact event sent", "event_id", evt.EventID)
}

// hashNormalized lowercases, trims, and SHA-256 hashes a value per Meta's requirements.
func hashNormalized(val string) string {
	v := strings.TrimSpace(strings.ToLower(val))
	if v == "" {
		return ""
	}
	h := sha256.Sum256([]byte(v))
	return hex.EncodeToString(h[:])
}

// hashPhone strips non-digits, prepends country code if needed, then hashes.
func hashPhone(val string) string {
	var digits strings.Builder
	for _, r := range val {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	d := digits.String()
	if d == "" {
		return ""
	}
	// Assume US if 10 digits
	if len(d) == 10 {
		d = "1" + d
	}
	h := sha256.Sum256([]byte(d))
	return hex.EncodeToString(h[:])
}
