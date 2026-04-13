package view

import "time"

// SiteName is the display name used in templates.
const SiteName = "Lo Mo Outfitting"

// SiteURL is the canonical base URL (no trailing slash).
const SiteURL = "https://lomooutfitting.com"

// DefaultOGImage is the fallback Open Graph image path.
const DefaultOGImage = "/static/img/hero-river.jpg"

// Tracking IDs and Turnstile site key, set once at startup from config.
var (
	PixelID          string
	GtagID           string
	TurnstileSiteKey string
)

// Year returns the current year for copyright notices.
func Year() int {
	return time.Now().Year()
}
