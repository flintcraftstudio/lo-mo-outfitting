# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Lo Mo Outfitting — a brochure website for a Montana fly fishing outfitter. Built on Firefly Software's standard-template (Go + templ + Tailwind CSS).

## Build & Dev Commands

```bash
mage Dev               # build CSS + templ generate + run server on :8080
mage InstallTailwind   # download Tailwind standalone CLI (one-time)
mage Build             # full production build: CSS + templ generate + go build
mage BuildCSS          # compile Tailwind only
mage GenerateTempl     # run templ generate only
mage BuildGo           # templ generate + go build (outputs ./bin/server)
docker compose up      # run via Docker
```

`mage Dev` is the recommended way to develop. It builds CSS, generates templ, and starts the server in one command.

## Architecture

- **Go stdlib `net/http`** for routing and serving (no framework)
- **templ** for type-safe HTML templates (`.templ` files in `internal/view/`)
- **Tailwind CSS v3** via standalone CLI (config at `tailwind/tailwind.config.js`, input at `tailwind/input.css`, output at `web/static/css/site.css`)
- **HTMX** for form interactions, **Alpine.js** for client-side interactivity (mobile nav toggle, etc.)
- **Mage** as the build task runner (`magefile.go`)

### Key directories

- `cmd/server/main.go` — entry point, routing, .env loading
- `internal/handler/` — HTTP handlers (home, contact, 404)
- `internal/view/` — templ templates (layout, pages, partials) + `shared.go` for constants
- `internal/config/` — env var loading into config struct
- `internal/mail/` — Postmark API client for contact form emails
- `internal/middleware/` — request logging
- `web/static/` — compiled CSS, JS (HTMX, Alpine), images

## Design System — "Teal & Bone" (DS-002)

The full UI guide is at `docs/lomo-ui-guide.md`. Key points:

- **Color palette**: `river-deep` (#1A2E30), `river-mid` (#2A4A4E), `teal` (#3E9E98), `teal-dark` (#2C7A75), `bone` (#F4F1EB), `stone` (#7A8E8F) — all defined as Tailwind theme extensions
- **Fonts**: Playfair Display (display/headlines), Source Sans 3 (body) — loaded via Google Fonts
- **Named font sizes**: use the custom scale (`text-eyebrow`, `text-h1`, `text-h2`, `text-stat`, etc.) instead of arbitrary rem values
- **Section wrapper pattern**: every section uses `<section class="w-full py-section"><div class="max-w-site mx-auto px-6">...</div></section>`
- **Teal discipline**: teal appears in one role per section (CTA button OR stat numbers, not both)

## Skills

- `/two-variation-site` — scaffold a two-direction brand comparison site
- `/qc` — pre-deploy quality control audit (read-only, does not fix issues)

## Environment

All env vars are optional with graceful degradation. Copy `.env.example` to `.env`. Key vars: `PORT`, `POSTMARK_SERVER_TOKEN`, `POSTMARK_FROM`, `POSTMARK_TO`, `GTAG_ID`, `PIXEL_ID`, `TURNSTILE_SITE_KEY`, `TURNSTILE_SECRET_KEY`.

## Design Context

### Users
Potential clients booking guided fly fishing trips on the Missouri River near Craig, Montana. Range from first-timers to experienced anglers. Many are out-of-state visitors; some are military/first responders eligible for the Heroes rate. Job to be done: find a credible, local outfitter and book with confidence.

### Brand Personality
**Three words:** Adventurous, Local, Honest

Direct and conversational voice. No marketing fluff. Sounds like a guide talking at the boat ramp. "No pretense" is the governing principle.

**Emotional goal:** Trust and calm confidence — competent, reliable, low-pressure. Credibility earned through specificity (license numbers, named guides, real bios) rather than polish.

### Aesthetic Direction
Earthy, restrained, grounded. Dark river tones with warm bone/ivory contrast. Late-afternoon light on the Missouri — warm but not bright, natural but not rustic-kitschy. Sharp corners (3px cards, 2px buttons). Light/dark sections alternate to create natural rhythm.

**Anti-references:** Not slick corporate outdoor brands. Not cheap cluttered booking platforms. Stay in the authentic middle.

### Design Principles
1. **Specificity over polish.** Real details (license numbers, named guides, concrete prices) earn more trust than visual refinement.
2. **Restraint is the aesthetic.** Tight palette, minimal radii, generous whitespace, no decorative elements.
3. **One accent per section.** Teal appears once per section in a single role — CTA or accent text, never both.
4. **Copy is UI.** The voice is a design element. Never use placeholder/lorem ipsum copy — use visible `[PLACEHOLDER]` markers.
5. **Mobile-first, scroll-friendly.** Single-column mobile stacks, tap-to-call, generous touch targets.

### Accessibility
WCAG 2.1 AA. Semantic HTML, proper heading hierarchy, keyboard navigable, descriptive alt text, `prefers-reduced-motion` respected.
