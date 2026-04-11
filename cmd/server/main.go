package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/firefly-software-mt/standard-template/internal/config"
	"github.com/firefly-software-mt/standard-template/internal/database"
	"github.com/firefly-software-mt/standard-template/internal/handler"
	adminhandler "github.com/firefly-software-mt/standard-template/internal/handler/admin"
	"github.com/firefly-software-mt/standard-template/internal/mail"
	"github.com/firefly-software-mt/standard-template/internal/middleware"
	"github.com/firefly-software-mt/standard-template/internal/view"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := loadEnv(".env"); err != nil {
		slog.Error("env error", "err", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	// Tracking pixels
	view.GtagID = cfg.GtagID
	view.PixelID = cfg.PixelID
	if cfg.GtagID == "" {
		slog.Warn("GTAG_ID not set, Google Analytics disabled")
	}
	if cfg.PixelID == "" {
		slog.Warn("PIXEL_ID not set, Facebook Pixel disabled")
	}

	// Turnstile
	view.TurnstileSiteKey = cfg.TurnstileSiteKey
	if cfg.TurnstileSiteKey == "" || cfg.TurnstileSecretKey == "" {
		slog.Warn("TURNSTILE_SITE_KEY or TURNSTILE_SECRET_KEY not set, Turnstile disabled")
	}

	// Database
	db, err := database.Open(cfg.DatabasePath)
	if err != nil {
		slog.Error("database error", "err", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("database ready", "path", cfg.DatabasePath)

	// Mail client (nil if Postmark is not configured)
	var mailer *mail.Client
	if cfg.PostmarkToken != "" {
		mailer = mail.NewClient(cfg.PostmarkToken, cfg.PostmarkFrom, cfg.PostmarkTo)
		slog.Info("postmark configured")
	} else {
		slog.Info("postmark not configured, contact form emails disabled")
	}

	mux := http.NewServeMux()

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Pages
	mux.Handle("GET /{$}", handler.Home())
	mux.Handle("GET /about", handler.About())
	mux.Handle("GET /guides", handler.Guides())
	mux.Handle("GET /contact", handler.Contact())
	mux.Handle("POST /contact", handler.ContactSubmit(mailer, cfg.TurnstileSecretKey, db))
	mux.Handle("GET /policies", handler.Policies())
	mux.Handle("GET /reviews", handler.Reviews())
	mux.Handle("GET /store", handler.Store())

	// Admin auth routes (public)
	mux.Handle("GET /admin/login", adminhandler.LoginPage())
	mux.Handle("POST /admin/login", adminhandler.LoginSubmit(db, cfg.AdminPasswordHash))
	mux.Handle("POST /admin/logout", adminhandler.Logout(db))

	// Admin protected routes
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("GET /admin/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/inquiries", http.StatusSeeOther)
	})
	// Inquiry board
	adminMux.Handle("GET /admin/inquiries", adminhandler.Inquiries(db))
	adminMux.Handle("GET /admin/inquiries/new", adminhandler.NewBookingForm())
	adminMux.Handle("POST /admin/inquiries/new", adminhandler.NewBookingSubmit(db))
	adminMux.Handle("GET /admin/inquiries/{id}", adminhandler.InquiryDetail(db))
	adminMux.Handle("POST /admin/inquiries/{id}/status", adminhandler.StatusUpdate(db))
	adminMux.Handle("POST /admin/inquiries/{id}/guide", adminhandler.GuideAssign(db))
	adminMux.Handle("POST /admin/inquiries/{id}/note", adminhandler.NoteAdd(db))
	adminMux.Handle("POST /admin/inquiries/{id}/payment", adminhandler.PaymentMethod(db))

	// Upcoming trips + All bookings
	adminMux.Handle("GET /admin/upcoming", adminhandler.Upcoming(db))
	adminMux.Handle("GET /admin/all", adminhandler.AllBookings(db))

	protectedAdmin := middleware.AdminAuth(db)(adminMux)
	mux.Handle("GET /admin/", protectedAdmin)
	mux.Handle("POST /admin/", protectedAdmin)

	// Clean expired sessions on startup
	if n, err := db.CleanExpiredSessions(); err != nil {
		slog.Error("clean expired sessions", "err", err)
	} else if n > 0 {
		slog.Info("cleaned expired sessions", "count", n)
	}

	// 404 catch-all
	mux.Handle("GET /", handler.NotFound())

	srv := middleware.Logging(logger)(mux)

	slog.Info("server starting", "addr", cfg.Addr())
	fmt.Printf("listening on %s\n", cfg.Addr())
	if err := http.ListenAndServe(cfg.Addr(), srv); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}

// loadEnv reads a .env file and sets environment variables if not already set.
func loadEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, line := range splitLines(string(data)) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		key, val, ok := splitOnce(line, '=')
		if !ok {
			continue
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
	return nil
}

// splitLines splits a string into non-empty lines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// splitOnce splits a string on the first occurrence of sep.
func splitOnce(s string, sep byte) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}
