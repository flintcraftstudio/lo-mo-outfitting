# standard-template

Static brochure websites offered at Firefly's "standard" tier. Built with Go, templ, and Tailwind CSS.

## Tech Stack

- **Go** (stdlib `net/http`) — server, routing, handlers
- **templ** — type-safe HTML templating
- **Tailwind CSS** (standalone CLI) — utility-first styling
- **HTMX** — form interactions without page reloads
- **Alpine.js** — lightweight client-side interactivity
- **Mage** — build task runner
- **Postmark** — transactional email for contact forms
- **Cloudflare Turnstile** — bot protection

## Getting Started

```bash
mage InstallTailwind   # download Tailwind CLI (once)
mage Dev               # terminal 1: watch CSS
go run ./cmd/server    # terminal 2: run server on :8080
```

Production build:

```bash
mage Build             # CSS + templ generate + go build
./bin/server
```

Docker:

```bash
docker compose up
```

## Project Structure

```
cmd/server/main.go              # entry point, routing, config
internal/
  config/config.go              # env var loading
  handler/                      # HTTP handlers (home, contact)
  middleware/                   # request logging
  mail/postmark.go              # Postmark API client
  view/                         # templ templates
    layout.templ                # base HTML wrapper
    home.templ                  # homepage
    contact.templ               # contact form
    nav.templ, footer.templ     # shared partials
    shared.go                   # constants (SiteName, tracking IDs)
tailwind/
  tailwind.config.js            # color palette, fonts, content paths
  input.css                     # font imports, Tailwind directives
web/static/
  css/site.css                  # compiled Tailwind output
  js/                           # HTMX, Alpine.js, custom scripts
```

## Environment Variables

All optional with graceful degradation. See `.env.example`.

| Variable | Purpose |
|---|---|
| `PORT` | Server port (default: 8080) |
| `POSTMARK_SERVER_TOKEN` | Postmark API key for contact form emails |
| `POSTMARK_FROM` | Sender email address |
| `POSTMARK_TO` | Recipient email address |
| `GTAG_ID` | Google Analytics tag |
| `PIXEL_ID` | Facebook Pixel ID |
| `TURNSTILE_SITE_KEY` | Cloudflare Turnstile site key |
| `TURNSTILE_SECRET_KEY` | Cloudflare Turnstile secret key |

---

## Claude Code Skills

This repo includes Claude Code skills invoked with `/skill-name` in conversation. Each skill is a structured prompt that guides Claude through a specific workflow.

### `/two-variation-site`

Build a brochure website presenting two distinct brand/design directions for a client to compare side-by-side.

**When to use:** Starting a new client project where you want to present two visual directions (e.g., "warm" vs "bold") from a single codebase.

**Inputs required before code is written:**

1. **Business details** — name, address, phone, email, hours, social links, tagline
2. **Brand guide for each variation** — color palette (hex values + roles), typography stack (families, sizes, weights for headline/body/accent/UI), voice/tone, layout personality
3. **Page copy** — approved text for each section of each page
4. **Images** — hero images, logos, team photos (or placeholders to use)
5. **Variation names/slugs** — evocative short names for each direction (e.g., "warm" / "bold", "classic" / "modern")

**What it produces:**

- Split-panel landing page at `/` comparing both directions
- Complete variation A site at `/[slug-a]/`
- Complete variation B site at `/[slug-b]/`
- Namespaced Tailwind color palettes (`va-*`, `vb-*`)
- Separate templ templates per variation with shared business data
- Scoped navigation (links stay within each variation's URL space)

**Example:**

```
/two-variation-site Henderson Bakery
```

Claude will ask for any missing inputs before writing code, present a structured summary for confirmation, then build the full site.

---

### `/qc`

Run a pre-deploy quality control check against the standard-tier site checklist.

**When to use:** The site is nearing completion and needs a final review before deployment. This is a read-only audit — it reports issues but does not fix them.

**Inputs required:**

1. **Client name** — used in the report header

**What it checks (7 sections):**

| Section | What it verifies |
|---|---|
| Project Structure | Required files exist (`main.go`, handlers, templates, Dockerfile, etc.) |
| Build Verification | `mage build` succeeds, `go vet` and `golangci-lint` pass |
| Functionality | Routing, contact form validation, Postmark integration, config safety |
| SEO | Unique titles, meta descriptions, heading hierarchy, Open Graph tags, robots.txt, sitemap |
| Accessibility | Semantic HTML, image alt text, form labels, keyboard navigation, skip link |
| Security | Security headers, CSRF protection, no hardcoded secrets, no localhost references |
| Deployment Readiness | `.env.example` complete, Docker config clean, all placeholder copy replaced |

**Output:** A structured report with tiered findings:

- **FAIL** — must fix before deploy
- **WARN** — should fix, does not block deploy
- **PASS** — requirement met

Final status: `READY`, `READY WITH WARNINGS`, or `NOT READY`.

**Example:**

```
/qc Henderson Bakery
```
