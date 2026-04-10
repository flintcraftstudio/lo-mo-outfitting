# Lo Mo Outfitting — Booking CMS Spec

**Stack:** Go · html/template · htmx · Alpine.js · Tailwind CSS
**Auth:** Single-user, session-based. Matt only. No roles, no multi-user.
**Route prefix:** `/admin/*` — protected by session middleware.

---

## 1. Data Model

### `booking_requests`

Primary table. One row per form submission.

```sql
CREATE TABLE booking_requests (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    emailed_at          DATETIME,
    ip_address          TEXT,

    -- Trip
    trip_type           TEXT NOT NULL,
        -- values: full_day_single | half_day_single | early_season |
        --         multiple_boats | heroes
    preferred_date      DATE NOT NULL,
    alternate_date      DATE,
    angler_count        INTEGER NOT NULL,
    youth_count         INTEGER NOT NULL DEFAULT 0,
    heroes              INTEGER NOT NULL DEFAULT 0, -- boolean 0/1

    -- Party
    experience          TEXT NOT NULL,
        -- values: never | some | comfortable | advanced
    lodging             TEXT NOT NULL,
        -- values: craig | wolf_creek | helena | great_falls |
        --         not_sure | other
    lodging_other       TEXT,
    client_notes        TEXT,
    referred_by         TEXT,

    -- Contact
    client_name         TEXT NOT NULL,
    client_email        TEXT NOT NULL,
    client_phone        TEXT NOT NULL,

    -- CMS state
    status              TEXT NOT NULL DEFAULT 'new',
        -- values: new | contacted | deposit_sent | confirmed | complete | cancelled
    guide_id            INTEGER REFERENCES guides(id),
    payment_method      TEXT,
        -- values: cash | venmo | stripe | other
        -- set by Matt when marking a booking confirmed
        -- null until deposit is received
    mat_notes           TEXT,
    status_updated_at   DATETIME
);
```

### `guides`

Lookup table for guide assignment. Populated manually by Matt or seeded at deploy.

```sql
CREATE TABLE guides (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    license     TEXT,
    active      INTEGER NOT NULL DEFAULT 1
);
```

Seed data:
```sql
INSERT INTO guides (name, license) VALUES
    ('Ria French',    'GUD-LIC-37359'),
    ('Colter Day',    NULL),
    ('Andrew Osborn', NULL),
    ('Dylan Huseby',  NULL),
    ('Sam Botz',      NULL),
    ('Rain Keating',  NULL),
    ('Dave Buck',     NULL);
```

### `booking_events`

Append-only activity log. Never updated, never deleted.

```sql
CREATE TABLE booking_events (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    booking_request_id  INTEGER NOT NULL REFERENCES booking_requests(id),
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    event_type          TEXT NOT NULL,
        -- values: submitted | email_sent | email_failed | status_changed |
        --         guide_assigned | note_added | payment_method_set
    detail              TEXT
        -- human-readable string logged with the event
        -- e.g. "Status changed: new → contacted"
        -- e.g. "Guide assigned: Ria French"
        -- e.g. "Note: Called twice, no answer. Try again Fri."
);
```

---

## 2. Routes

All routes are under `/admin`. Every handler checks for a valid session cookie
before processing — redirect to `/admin/login` if missing or expired.

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/admin/login` | `HandleLoginPage` | Login form |
| POST | `/admin/login` | `HandleLoginSubmit` | Validate password, set session cookie |
| POST | `/admin/logout` | `HandleLogout` | Clear session cookie |
| GET | `/admin` | redirect → `/admin/inquiries` | |
| GET | `/admin/inquiries` | `HandleInquiries` | Inquiry board (default view) |
| GET | `/admin/inquiries/{id}` | `HandleInquiryDetail` | Detail pane via htmx swap |
| POST | `/admin/inquiries/{id}/status` | `HandleStatusUpdate` | Update status, log event |
| POST | `/admin/inquiries/{id}/guide` | `HandleGuideAssign` | Assign guide, log event |
| POST | `/admin/inquiries/{id}/note` | `HandleNoteAdd` | Save Matt's note, log event |
| POST | `/admin/inquiries/{id}/payment` | `HandlePaymentMethod` | Set payment method, log event |
| GET | `/admin/upcoming` | `HandleUpcoming` | Upcoming confirmed trips |
| GET | `/admin/all` | `HandleAllBookings` | Searchable full history |

---

## 3. Views

### 3.1 Login

Simple centered form. Password only — no username, Matt is the only user.
Session cookie expires after 7 days. No "remember me" toggle needed.

Fields:
- Password (required)

On success: redirect to `/admin/inquiries`.
On failure: inline error — "Incorrect password."

Store the hashed password in Go config (env var `ADMIN_PASSWORD_HASH`),
not in the database. Use `bcrypt`.

---

### 3.2 Inquiry Board (`/admin/inquiries`)

Two-panel layout on desktop. Single column on mobile (list only, tap to open detail).

**Left panel — inquiry list**

Top: count pills — one per status showing the count of open requests in that state.
Clicking a pill filters the list. Default filter: all open (excludes complete and cancelled).

Filter options:
- All open
- New
- Deposit pending
- Upcoming (confirmed, date in the future)
- Complete

Each row in the list shows:
- Client name
- Trip type + angler count + preferred date
- Status indicator (coloured dot + label)
- Tags for notable flags: Heroes, Youth, Multiple boats

Rows sorted by: `preferred_date ASC` within each status group,
with `new` requests floated to the top regardless of date.

**Right panel — detail view**

Loaded via `hx-get="/admin/inquiries/{id}"` `hx-target="#detail-pane"` on row click.
Default state (no row selected): empty with a prompt — "Select an inquiry to view details."

Detail pane sections, top to bottom:

1. **Header** — client name, submission timestamp, status select dropdown
2. **Trip** — all trip fields in a 2-column grid
3. **Party** — experience, lodging, referred by
4. **Client notes** — freetext from the booking form, displayed as a blockquote
5. **Contact** — email (mailto link) and phone (tel link)
6. **Payment** — payment method select (Cash · Venmo · Stripe · Other). Visible at
   all stages but only meaningful once a deposit is collected. Saves immediately on
   change via `hx-post="/admin/inquiries/{id}/payment"` `hx-trigger="change"`.
   Writes a `payment_method_set` event to the timeline.
7. **Matt's notes** — private textarea, saves on blur via `hx-post` + `hx-trigger="blur"`.
   Append-only in the event log but the field itself is editable (last value wins).
8. **Activity timeline** — all `booking_events` rows for this request, newest last
9. **Action bar** — contextual buttons based on current status (see below)

**Status select dropdown**

Changing the value triggers `hx-post="/admin/inquiries/{id}/status"` immediately
(Alpine.js `@change`). No confirm dialog. The event log records the transition.
The list panel count pills update via an out-of-band htmx swap
(`hx-swap-oob="true"` on the count pills fragment).

**Action bar — buttons by status**

| Current status | Buttons shown |
|----------------|---------------|
| new | Mark as contacted · Assign guide · Add note · Cancel inquiry |
| contacted | Deposit sent · Assign guide · Add note · Cancel inquiry |
| deposit_sent | Mark confirmed · Assign guide · Add note · Cancel inquiry |
| confirmed | Mark complete · Add note · Cancel inquiry |
| complete | Add note |
| cancelled | Add note |

"Mark as contacted" / "Deposit sent" / "Mark confirmed" / "Mark complete" are all
`hx-post` calls to `/admin/inquiries/{id}/status` with the target status value.
They write a `status_changed` event and re-render the detail pane.

"Assign guide" opens an inline dropdown (Alpine.js `x-show`) listing active guides.
Selecting one posts to `/admin/inquiries/{id}/guide`, writes a `guide_assigned` event,
and re-renders the detail pane. Shows a conflict warning if the selected guide
already has a confirmed booking on the same preferred date — warning only, not a block.

"Add note" opens an inline textarea (Alpine.js `x-show`). Submitting posts to
`/admin/inquiries/{id}/note`, writes a `note_added` event with the note text,
appends to the timeline, and clears the textarea.

"Cancel inquiry" shows an inline confirmation prompt before posting.

---

### 3.3 Upcoming Trips (`/admin/upcoming`)

Single-column list, no detail panel. Confirmed trips only, sorted by `preferred_date ASC`,
grouped by calendar month. Window: next 60 days from today.

**Summary bar** — 4 metric cards at the top:
- Trips this period (count of confirmed in the window)
- Guides scheduled (distinct guide_ids in confirmed trips this period)
- Total anglers (sum of angler_count for confirmed trips this period)
- Deposit outstanding (count of confirmed trips where status was set to confirmed
  but no payment event is recorded — see note below)

**Trip card** — each confirmed booking shows:
- Date block (day number + weekday abbreviation). Today's date gets a teal highlight.
- Client name
- Trip type + angler count + lodging location
- Assigned guide name(s) — for multiple-boat trips, list all assigned guides
- Tags: trip type, Heroes, Youth, deposit pending flag

Clicking a trip card navigates to `/admin/inquiries/{id}` — the full detail view.

**Deposit outstanding flag** — shown when `status = 'confirmed'` and
`payment_method IS NULL`. Matt sets the payment method when he collects the deposit,
so a null value on a confirmed booking means deposit has not yet been recorded.
This is simpler and more reliable than inferring from status transition history.

---

### 3.4 All Bookings (`/admin/all`)

Paginated table. 25 rows per page. All statuses including complete and cancelled.

Columns: Date submitted · Client name · Trip type · Preferred date · Status · Guide

Search bar at the top — filters by client name, email, or phone.
`hx-get="/admin/all"` `hx-trigger="input changed delay:300ms"` — live search,
replaces the table body only.

No inline detail panel. Clicking a row navigates to `/admin/inquiries/{id}`.

---

## 4. Auth & Session

Single password stored as a bcrypt hash in the environment:

```
ADMIN_PASSWORD_HASH=<bcrypt hash of Matt's password>
```

Session stored server-side in a SQLite table:

```sql
CREATE TABLE admin_sessions (
    token       TEXT PRIMARY KEY,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at  DATETIME NOT NULL
);
```

Session token is a 32-byte random hex string, stored in a cookie named `lomo_session`.
Cookie attributes: `HttpOnly`, `Secure`, `SameSite=Strict`, 7-day expiry.

Session middleware checks:
1. Cookie present
2. Token exists in `admin_sessions`
3. `expires_at` is in the future

On any failure: clear the cookie, redirect to `/admin/login`.

No password reset flow needed — Matt can update the env var and redeploy.

---

## 5. Postmark Email (outbound only)

On form submission the Go handler sends one email to Matt via Postmark.

**To:** `matt@lomooutfitting.com` (from config)
**From:** `bookings@lomooutfitting.com` (or same as above — confirm domain setup)
**Reply-To:** submitter's email address
**Subject:** `New booking request — [trip type] · [preferred date] · [client name]`

Body: plain text, structured summary of all form fields.
See `lomo-booking-form.md` for the exact template.

The handler writes `emailed_at` to the row on success. If Postmark returns an error,
`emailed_at` stays null — the row is still saved. Matt can see un-emailed submissions
in the inquiry board (they appear as normal new requests; the missing email is
visible in the event timeline as an `email_failed` event).

No inbound email parsing. No reply tracking. Matt replies manually from his email client.

---

## 6. Navigation

Admin nav is a simple top bar, always visible:

- **Lo Mo Bookings** (wordmark, links to `/admin/inquiries`)
- **Inquiries** (with new-request badge count)
- **Upcoming**
- **All bookings**
- **Log out** (right side)

The badge count on Inquiries shows the number of `new` status requests.
It updates on every page load — no websockets, no polling.

---

## 7. Implementation Notes

- All admin templates live in `templates/admin/` and extend a shared
  `admin_base.html` layout that includes the nav and session check.
- htmx is already a project dependency — no additional JS libraries needed
  for the admin beyond Alpine.js (also already present).
- The guide conflict check on assignment is a simple query:
  ```sql
  SELECT COUNT(*) FROM booking_requests
  WHERE guide_id = ?
  AND preferred_date = ?
  AND status = 'confirmed'
  AND id != ?
  ```
  Returns a warning fragment if count > 0. Rendered inline above the guide dropdown.
- The `/admin/all` search query uses `LIKE` on name, email, and phone — sufficient
  for a single-user tool with a small dataset. No full-text search needed.
- Pagination on `/admin/all`: use `LIMIT 25 OFFSET (page-1)*25`. Pass `?page=N`
  as a query param. htmx replaces the table body and pagination controls only.
- Matt's private notes field saves on blur (`hx-trigger="blur"`). This is the
  one field where the last-write value is authoritative. The event log records
  each save, so previous note versions are recoverable from `booking_events.detail`
  if needed.
- The `status_updated_at` column on `booking_requests` is updated on every
  status change. Used to surface stale inquiries — requests that have been in
  `contacted` or `deposit_sent` for more than 5 days could be highlighted in
  the list (amber row background). Implement as a view-layer concern, not
  a separate column.
- No soft delete. Cancelled inquiries stay in the database with `status = cancelled`.
  They are excluded from the "all open" filter and the upcoming view but visible
  in `/admin/all`.
- The payment method field (`cash | venmo | stripe | other`) is set by Matt at the
  point of deposit collection. It is not part of the status flow — Matt can set it
  independently at any stage. The deposit outstanding flag on the upcoming view
  checks `payment_method IS NULL` on confirmed bookings. The `HandlePaymentMethod`
  route updates the column and writes a `payment_method_set` event with the selected
  value as the detail string (e.g. "Payment method set: Venmo").

---

## 8. What This Is Not

- Not a calendar. Guides are assigned to requests; no availability grid is shown.
- Not a payment processor. Deposit links are sent by Matt outside this system.
- Not a client portal. Clients have no login, no booking status page.
- Not a notification system. No emails to guides, no SMS, no reminders.
- Not a reporting dashboard. The summary bar on the upcoming view is the extent
  of analytics.

All of the above are potential Phase 2 additions. None are in scope for this build.