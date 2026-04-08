//go:build mage

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/firefly-software-mt/standard-template/internal/database"
	"github.com/magefile/mage/sh"
	"golang.org/x/crypto/bcrypt"
)

const tailwindVersion = "v3.4.17"

// InstallTailwind downloads the Tailwind standalone CLI for the current platform
func InstallTailwind() error {
	binary := tailwindBinaryPath()
	if _, err := os.Stat(binary); err == nil {
		fmt.Println("Tailwind already installed, skipping.")
		return nil
	}

	url := tailwindDownloadURL()
	fmt.Printf("Downloading Tailwind %s from %s\n", tailwindVersion, url)

	if err := sh.Run("curl", "-sLo", binary, url); err != nil {
		return err
	}
	return sh.Run("chmod", "+x", binary)
}

// BuildCSS compiles Tailwind CSS
func BuildCSS() error {
	return sh.Run(
		tailwindBinaryPath(),
		"-c", "./tailwind/tailwind.config.js",
		"-i", "./tailwind/input.css",
		"-o", "./web/static/css/site.css",
		"--minify",
	)
}

// GenerateTempl runs templ generate
func GenerateTempl() error {
	return sh.Run("templ", "generate")
}

// BuildGo compiles the Go binary
func BuildGo() error {
	if err := GenerateTempl(); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "./bin/server", "./cmd/server")
}

// Build runs a full production build
func Build() error {
	if err := BuildCSS(); err != nil {
		return err
	}
	return BuildGo()
}

// Dev runs Tailwind in watch mode (run `go run ./cmd/server` in a second terminal)
func Dev() error {
	return sh.Run(
		tailwindBinaryPath(),
		"-c", "./tailwind/tailwind.config.js",
		"-i", "./tailwind/input.css",
		"-o", "./web/static/css/site.css",
		"--watch",
	)
}

// Run starts the server
func Run() error {
	return sh.Run("./bin/server")
}

// Seed sets the admin password. Usage: mage seed <password>
func Seed(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/lomo.db"
	}

	db, err := database.Open(dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	if err := db.SetSetting("admin_password_hash", string(hash)); err != nil {
		return fmt.Errorf("save password hash: %w", err)
	}

	fmt.Println("Admin password set successfully.")
	return nil
}

func tailwindBinaryPath() string {
	if runtime.GOOS == "windows" {
		return "./tailwind/tailwindcss.exe"
	}
	return "./tailwind/tailwindcss"
}

func tailwindDownloadURL() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	osName := map[string]string{
		"darwin":  "macos",
		"linux":   "linux",
		"windows": "windows",
	}[os]

	archName := map[string]string{
		"amd64": "x64",
		"arm64": "arm64",
	}[arch]

	ext := ""
	if os == "windows" {
		ext = ".exe"
	}

	return fmt.Sprintf(
		"https://github.com/tailwindlabs/tailwindcss/releases/download/%s/tailwindcss-%s-%s%s",
		tailwindVersion, osName, archName, ext,
	)
}