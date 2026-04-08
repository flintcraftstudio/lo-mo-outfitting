package admin

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/firefly-software-mt/standard-template/internal/database"
	"github.com/firefly-software-mt/standard-template/internal/view/admin"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionCookieName = "lomo_session"
	sessionDuration   = 7 * 24 * time.Hour
)

// LoginPage renders the admin login form.
func LoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := admin.LoginPage("").Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// LoginSubmit validates the password, creates a session, and redirects to the admin.
// Checks the DB for the password hash first, falls back to the config value.
func LoginSubmit(db *database.DB, configPasswordHash string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		password := r.FormValue("password")

		// Read hash from DB, fall back to config env var
		passwordHash := configPasswordHash
		if dbHash, err := db.GetSetting("admin_password_hash"); err == nil && dbHash != "" {
			passwordHash = dbHash
		}

		if passwordHash == "" {
			slog.Error("no admin password configured")
			if err := admin.LoginPage("Admin password not configured.").Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
			if err := admin.LoginPage("Incorrect password.").Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		token, err := database.NewSessionToken()
		if err != nil {
			slog.Error("generate session token", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().Add(sessionDuration)
		if err := db.CreateSession(token, expiresAt); err != nil {
			slog.Error("create session", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    token,
			Path:     "/admin",
			MaxAge:   int(sessionDuration.Seconds()),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})

		http.Redirect(w, r, "/admin/inquiries", http.StatusSeeOther)
	}
}

// Logout clears the session and redirects to login.
func Logout(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(sessionCookieName); err == nil {
			if err := db.DeleteSession(cookie.Value); err != nil {
				slog.Error("delete session", "err", err)
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    "",
			Path:     "/admin",
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}
}
