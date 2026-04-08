package middleware

import (
	"net/http"

	"github.com/firefly-software-mt/standard-template/internal/database"
)

const sessionCookieName = "lomo_session"

// AdminAuth returns middleware that checks for a valid session cookie.
// On failure it clears the cookie and redirects to /admin/login.
func AdminAuth(db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil {
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}

			session, err := db.GetSession(cookie.Value)
			if err != nil || session == nil {
				// Clear invalid/expired cookie
				http.SetCookie(w, &http.Cookie{
					Name:     sessionCookieName,
					Value:    "",
					Path:     "/admin",
					MaxAge:   -1,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				})
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
