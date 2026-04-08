package database

// Guide represents a row in the guides table.
type Guide struct {
	ID      int64
	Name    string
	License string
	Active  bool
}

// ListActiveGuides returns all guides where active = 1.
func (db *DB) ListActiveGuides() ([]Guide, error) {
	rows, err := db.Query("SELECT id, name, COALESCE(license, ''), active FROM guides WHERE active = 1 ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guides []Guide
	for rows.Next() {
		var g Guide
		var active int
		if err := rows.Scan(&g.ID, &g.Name, &g.License, &active); err != nil {
			return nil, err
		}
		g.Active = active == 1
		guides = append(guides, g)
	}
	return guides, rows.Err()
}

// GetGuide returns a single guide by ID.
func (db *DB) GetGuide(id int64) (*Guide, error) {
	var g Guide
	var active int
	err := db.QueryRow("SELECT id, name, COALESCE(license, ''), active FROM guides WHERE id = ?", id).
		Scan(&g.ID, &g.Name, &g.License, &active)
	if err != nil {
		return nil, err
	}
	g.Active = active == 1
	return &g, nil
}
