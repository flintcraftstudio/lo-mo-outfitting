package database

import "database/sql"

// GetSetting retrieves a value from admin_settings by key. Returns empty string if not found.
func (db *DB) GetSetting(key string) (string, error) {
	var value string
	err := db.QueryRow("SELECT value FROM admin_settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

// SetSetting upserts a value in admin_settings.
func (db *DB) SetSetting(key, value string) error {
	_, err := db.Exec(
		"INSERT INTO admin_settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value",
		key, value,
	)
	return err
}
