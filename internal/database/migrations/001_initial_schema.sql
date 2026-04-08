CREATE TABLE IF NOT EXISTS booking_requests (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    emailed_at          DATETIME,
    ip_address          TEXT,

    -- Trip
    trip_type           TEXT NOT NULL,
    preferred_date      DATE NOT NULL,
    alternate_date      DATE,
    angler_count        TEXT NOT NULL,
    youth_count         TEXT NOT NULL DEFAULT '0',
    heroes              INTEGER NOT NULL DEFAULT 0,

    -- Party
    experience          TEXT NOT NULL,
    lodging             TEXT NOT NULL,
    lodging_other       TEXT,
    client_notes        TEXT,
    referred_by         TEXT,

    -- Contact
    client_name         TEXT NOT NULL,
    client_email        TEXT NOT NULL,
    client_phone        TEXT NOT NULL,

    -- CMS state
    status              TEXT NOT NULL DEFAULT 'new',
    guide_id            INTEGER REFERENCES guides(id),
    payment_method      TEXT,
    mat_notes           TEXT,
    status_updated_at   DATETIME
);

CREATE TABLE IF NOT EXISTS guides (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    license     TEXT,
    active      INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS booking_events (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    booking_request_id  INTEGER NOT NULL REFERENCES booking_requests(id),
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    event_type          TEXT NOT NULL,
    detail              TEXT
);

CREATE TABLE IF NOT EXISTS admin_sessions (
    token       TEXT PRIMARY KEY,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at  DATETIME NOT NULL
);
