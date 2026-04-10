# Lo Mo Outfitting — UI & Typography Implementation Guide
**For coding agent use. Stack: Go · html/template · htmx · Alpine.js · Tailwind CSS v3**
**Design system: Teal & Bone (DS-002)**

---

## 1. Tailwind Configuration

Extend `tailwind.config.js` with these values. Do not use arbitrary values in templates — use these tokens only.

```js
// tailwind.config.js
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./static/**/*.js"],
  theme: {
    extend: {
      colors: {
        // Palette — Teal & Bone
        "river-deep":  "#1A2E30",   // primary dark — nav, hero, footer
        "river-mid":   "#2A4A4E",   // surface dark — guide cards, about
        "teal":        "#3E9E98",   // accent — CTAs, stats, links, stars
        "teal-dark":   "#2C7A75",   // CTA hover, CTA strip bg
        "bone":        "#F4F1EB",   // alt section bg, review cards
        "stone":       "#7A8E8F",   // secondary text, captions
        "slate-text":  "#1A2E30",   // body text on light bg (same as river-deep)
        "ivory":       "#F4F1EB",   // alias — text on dark sections
      },
      fontFamily: {
        display: ["Playfair Display", "Georgia", "serif"],
        body:    ["Source Sans 3", "system-ui", "sans-serif"],
      },
      fontSize: {
        // Named scale — use these, not arbitrary rem values
        "eyebrow":  ["0.6875rem", { lineHeight: "1", letterSpacing: "0.15em" }], // 11px
        "label":    ["0.5625rem", { lineHeight: "1", letterSpacing: "0.12em" }], // 9px
        "body-sm":  ["0.875rem",  { lineHeight: "1.65" }],                       // 14px
        "body":     ["1rem",      { lineHeight: "1.65" }],                       // 16px
        "price":    ["1.4rem",    { lineHeight: "1.2"  }],
        "stat":     ["2rem",      { lineHeight: "1.1"  }],
        "h3":       ["1.15rem",   { lineHeight: "1.3"  }],
        "h2":       ["2rem",      { lineHeight: "1.2"  }],
        "h1":       ["3rem",      { lineHeight: "1.1"  }],
        "h1-mob":   ["2rem",      { lineHeight: "1.1"  }],
      },
      spacing: {
        "section": "3.5rem",   // --sp-section: vertical padding on all sections
      },
      borderRadius: {
        "card": "3px",
        "btn":  "2px",
      },
      maxWidth: {
        "site": "1280px",      // max-w-site — outer container
        "prose": "520px",      // max-w-prose — section subheads
        "hero-sub": "480px",   // max-w-hero-sub — hero subheadline
      },
    },
  },
  plugins: [],
}
```

---

## 2. Google Fonts

Add to the `<head>` of the base layout template. Load both weights needed, nothing extra.

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,700;1,400&family=Source+Sans+3:wght@300;400;600;700&display=swap" rel="stylesheet">
```

**Weights in use:**
- Playfair Display 700 — headlines, stats, prices, nav wordmark
- Playfair Display 400 italic — review quotes, hero sub on selected sections
- Source Sans 3 300 — body copy, footer secondary text
- Source Sans 3 400 — general body
- Source Sans 3 600 — nav links, card details
- Source Sans 3 700 — eyebrow labels, CTAs, badges, link arrows

---

## 3. CSS Custom Properties

Add to `static/css/site.css` or a `<style>` block in the base layout. These are the authoritative values — Tailwind config mirrors them.

```css
:root {
  --river-deep:  #1A2E30;
  --river-mid:   #2A4A4E;
  --teal:        #3E9E98;
  --teal-dark:   #2C7A75;
  --bone:        #F4F1EB;
  --stone:       #7A8E8F;
  --border:      #E0DBD2;
  --ivory:       #F4F1EB;

  --font-display: "Playfair Display", Georgia, serif;
  --font-body:    "Source Sans 3", system-ui, sans-serif;

  --sp-xs:      0.25rem;
  --sp-sm:      0.5rem;
  --sp-md:      1rem;
  --sp-lg:      1.5rem;
  --sp-xl:      2rem;
  --sp-section: 3.5rem;
}

/* Tap-to-call — all phone number links */
a[href^="tel"] {
  color: var(--teal);
  font-weight: 600;
  text-decoration: none;
}

/* Smooth scroll */
html { scroll-behavior: smooth; }

/* Image defaults */
img {
  display: block;
  max-width: 100%;
  height: auto;
}
```

---

## 4. Typography Rules

### Display (Playfair Display)

| Element              | Classes                                                              |
|----------------------|----------------------------------------------------------------------|
| Hero headline        | `font-display font-bold text-h1 md:text-h1 text-ivory leading-tight` |
| Section headline     | `font-display font-bold text-h2 text-river-deep leading-snug`        |
| Section headline inv | `font-display font-bold text-h2 text-ivory leading-snug`             |
| Stat number          | `font-display font-bold text-stat text-teal`                         |
| Price                | `font-display font-bold text-price text-river-deep`                  |
| Review quote         | `font-display italic text-body text-river-deep leading-relaxed`      |
| Pull quote           | `font-display italic text-h3 text-ivory`                             |
| Nav wordmark         | `font-display font-bold text-base text-ivory`                        |

### Body (Source Sans 3)

| Element              | Classes                                                                    |
|----------------------|----------------------------------------------------------------------------|
| Eyebrow label        | `font-body font-bold text-eyebrow text-teal uppercase tracking-[0.15em]`   |
| Stat label           | `font-body font-semibold text-label uppercase tracking-[0.1em] text-ivory/40` |
| Nav links            | `font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-ivory/60` |
| Nav phone            | `font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-teal` |
| Body copy            | `font-body font-light text-body text-stone leading-relaxed`                |
| Body copy inv        | `font-body font-light text-body text-ivory/70 leading-relaxed`             |
| Card eyebrow tag     | `font-body font-bold text-label uppercase tracking-[0.12em] text-teal`     |
| Card detail          | `font-body font-light text-body-sm text-stone leading-relaxed`             |
| Footer col title     | `font-body font-medium text-label uppercase tracking-[0.12em] text-ivory/30` |
| Footer link          | `font-body text-body-sm text-ivory/60 hover:text-ivory transition-colors`  |
| Legal / caption      | `font-body font-light text-xs text-ivory/30`                               |
| Button label         | `font-body font-bold text-eyebrow uppercase tracking-[0.1em]`              |
| Badge label          | `font-body font-bold text-[0.5625rem] uppercase tracking-[0.08em]`         |

---

## 5. Section Wrapper Pattern

Every section uses this outer wrapper. Background class changes per section — see rhythm table below.

```html
<section class="w-full py-section">
  <div class="max-w-site mx-auto px-6">
    <!-- section content -->
  </div>
</section>
```

**Section rhythm — background classes in sequence:**

| #  | Section        | Outer class                           |
|----|----------------|---------------------------------------|
| 01 | Nav            | `bg-river-deep` (sticky, not a section) |
| 02 | Hero           | `bg-river-deep` (+ photo overlay)     |
| 03 | Stats bar      | `bg-river-deep border-t border-white/[0.06]` |
| 04 | Trip options   | `bg-bone`                             |
| 05 | Guides preview | `bg-river-deep`                       |
| 06 | Heroes program | `bg-bone`                             |
| 07 | Reviews        | `bg-white`                            |
| 08 | About Matt     | `bg-river-mid`                        |
| 09 | CTA strip      | `bg-teal-dark`                        |
| 10 | Footer         | `bg-river-deep`                       |

---

## 6. Component Patterns

### 6.1 Navigation

```html
<nav
  class="w-full bg-river-deep sticky top-0 z-50"
  x-data="{ open: false }"
>
  <div class="max-w-site mx-auto px-6 h-[62px] flex items-center justify-between">

    <!-- Wordmark -->
    <a href="/" class="font-display font-bold text-base text-ivory">
      Lo Mo Outfitting
    </a>

    <!-- Desktop links -->
    <div class="hidden md:flex items-center gap-6">
      <a href="/trips"           class="font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-ivory/60 hover:text-ivory transition-colors">Trips</a>
      <a href="/meet-our-guides" class="font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-ivory/60 hover:text-ivory transition-colors">Guides</a>
      <a href="/read-our-reviews"class="font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-ivory/60 hover:text-ivory transition-colors">Reviews</a>
      <a href="/lo-mo-store"     class="font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-ivory/60 hover:text-ivory transition-colors">Store</a>
      <a href="tel:+1406XXXXXXX" class="font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-teal">[PHONE]</a>
      <a href="/booking"         class="font-body font-bold text-eyebrow uppercase tracking-[0.08em] bg-teal hover:bg-teal-dark text-white px-4 py-[7px] rounded-btn transition-colors">Book Now</a>
    </div>

    <!-- Mobile hamburger -->
    <button
      class="md:hidden text-ivory/70"
      @click="open = !open"
      aria-label="Toggle menu"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path x-show="!open" stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16"/>
        <path x-show="open"  stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
      </svg>
    </button>
  </div>

  <!-- Mobile drawer -->
  <div
    x-show="open"
    x-transition
    class="md:hidden bg-river-deep border-t border-white/10 px-6 pb-6 flex flex-col gap-5"
  >
    <a href="/trips"            class="font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-ivory/70 pt-5">Trips</a>
    <a href="/meet-our-guides"  class="font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-ivory/70">Guides</a>
    <a href="/read-our-reviews" class="font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-ivory/70">Reviews</a>
    <a href="/lo-mo-store"      class="font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-ivory/70">Store</a>
    <a href="tel:+1406XXXXXXX"  class="font-body font-semibold text-eyebrow uppercase tracking-[0.1em] text-teal">[PHONE]</a>
    <a href="/booking"          class="font-body font-bold text-eyebrow uppercase tracking-[0.08em] bg-teal text-white px-4 py-3 rounded-btn text-center">Book Now</a>
  </div>
</nav>
```

---

### 6.2 Hero

```html
<section class="relative w-full bg-river-deep overflow-hidden min-h-[480px] flex flex-col justify-end">

  <!-- Background photo with overlay -->
  <div class="absolute inset-0">
    <img
      src="/static/img/hero-river.webp"
      alt="Missouri River near Craig, Montana"
      class="w-full h-full object-cover"
      loading="eager"
    >
    <div class="absolute inset-0 bg-river-deep/65"></div>
  </div>

  <!-- Content -->
  <div class="relative max-w-site mx-auto px-6 pb-16 pt-24">
    <p class="font-body font-bold text-eyebrow uppercase tracking-[0.15em] text-teal mb-4">
      Missouri River &middot; Craig, Montana
    </p>
    <h1 class="font-display font-bold text-h1-mob md:text-h1 text-ivory leading-tight mb-5 max-w-xl">
      The Missouri.<br>
      Local guides.<br>
      No pretense.
    </h1>
    <p class="font-body font-light text-body text-ivory/70 leading-relaxed mb-8 max-w-hero-sub">
      Born and raised on this river. Matt Mohar and a crew of seven guides
      who actually live here &mdash; not a booking platform in waders.
    </p>
    <div class="flex flex-wrap gap-4">
      <a href="/booking"
         class="font-body font-bold text-eyebrow uppercase tracking-[0.1em] bg-teal hover:bg-teal-dark text-white px-7 py-3 rounded-btn transition-colors">
        Book a Trip
      </a>
      <a href="/meet-our-guides"
         class="font-body font-semibold text-eyebrow uppercase tracking-[0.08em] border border-ivory/35 text-ivory hover:bg-white/10 px-6 py-3 rounded-btn transition-colors">
        Meet Our Guides
      </a>
    </div>
  </div>
</section>
```

---

### 6.3 Stats Bar

```html
<div class="w-full bg-river-deep border-t border-white/[0.06]">
  <div class="max-w-site mx-auto px-6 py-7 grid grid-cols-2 md:grid-cols-4 gap-6">

    <!-- Repeat this block × 4 -->
    <div class="text-center">
      <span class="font-display font-bold text-stat text-teal block">20+</span>
      <span class="font-body font-semibold text-label uppercase tracking-[0.1em] text-ivory/40 mt-1 block">Years on the Mo</span>
    </div>

    <div class="text-center">
      <span class="font-display font-bold text-stat text-teal block">7</span>
      <span class="font-body font-semibold text-label uppercase tracking-[0.1em] text-ivory/40 mt-1 block">Licensed Guides</span>
    </div>

    <div class="text-center">
      <span class="font-display font-bold text-stat text-teal block">35mi</span>
      <span class="font-body font-semibold text-label uppercase tracking-[0.1em] text-ivory/40 mt-1 block">Blue-Ribbon Water</span>
    </div>

    <div class="text-center">
      <span class="font-display font-bold text-stat text-teal block">5K+</span>
      <span class="font-body font-semibold text-label uppercase tracking-[0.1em] text-ivory/40 mt-1 block">Trout per Mile</span>
    </div>

  </div>
</div>
```

---

### 6.4 Section Eyebrow + Headline Pattern

Used at the top of every content section. Copy the exact class strings.

```html
<!-- Eyebrow label -->
<p class="font-body font-bold text-eyebrow uppercase tracking-[0.15em] text-teal mb-3">
  Float Trip Options
</p>

<!-- Headline — light background -->
<h2 class="font-display font-bold text-h2 text-river-deep leading-snug mb-3">
  Choose Your Day on the Water
</h2>

<!-- Headline — dark background -->
<h2 class="font-display font-bold text-h2 text-ivory leading-snug mb-3">
  Guides Who Live on the Water
</h2>

<!-- Section subhead — light -->
<p class="font-body font-light text-body text-stone leading-relaxed max-w-prose mb-10">
  All trips include fly fishing gear and shuttle transportation.
</p>

<!-- Section subhead — dark -->
<p class="font-body font-light text-body text-ivory/65 leading-relaxed max-w-prose mb-10">
  Seven licensed guides. Most of them fish this river on their days off.
</p>
```

---

### 6.5 Trip Card

```html
<!-- 2×2 grid wrapper -->
<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">

  <!-- Featured card (amber left border) -->
  <div class="bg-white border border-[#E0DBD2] border-l-[3px] border-l-teal rounded-card p-6">
    <span class="font-body font-bold text-[0.5625rem] uppercase tracking-[0.08em] bg-teal text-white px-2 py-[3px] rounded-[2px] inline-block mb-2">
      Most Popular
    </span>
    <h3 class="font-display font-bold text-h3 text-river-deep mb-1">Full Day — Single Boat</h3>
    <p class="font-display font-bold text-price text-river-deep mb-1">$675</p>
    <p class="font-body font-light text-body-sm text-stone leading-relaxed mb-4">
      1&ndash;2 anglers &middot; 7&ndash;8 hours<br>
      Gear, lunch, drinks, shuttle included
    </p>
    <a href="/booking" class="font-body font-bold text-label uppercase tracking-[0.08em] text-teal hover:text-teal-dark transition-colors">
      Book This Trip &rarr;
    </a>
  </div>

  <!-- Standard card -->
  <div class="bg-white border border-[#E0DBD2] rounded-card p-6">
    <span class="font-body font-bold text-label uppercase tracking-[0.12em] text-teal block mb-2">
      Half Day
    </span>
    <h3 class="font-display font-bold text-h3 text-river-deep mb-1">Half Day — Single Boat</h3>
    <p class="font-display font-bold text-price text-river-deep mb-1">$575</p>
    <p class="font-body font-light text-body-sm text-stone leading-relaxed mb-4">
      1&ndash;2 anglers &middot; 4&ndash;5 hours<br>
      Gear and shuttle included &middot; No lunch
    </p>
    <a href="/booking" class="font-body font-bold text-label uppercase tracking-[0.08em] text-teal hover:text-teal-dark transition-colors">
      Book This Trip &rarr;
    </a>
  </div>

</div>
```

---

### 6.6 Guide Card

```html
<!-- 3-column grid wrapper -->
<div class="grid grid-cols-1 md:grid-cols-3 gap-5">

  <div class="bg-river-mid rounded-card overflow-hidden">

    <!-- Photo slot — replace div with img when photo is available -->
    <div class="aspect-[3/4] bg-[#243d40] flex items-center justify-center relative">
      <!-- When photo available: -->
      <!-- <img src="/static/img/guides/ria-french.webp" alt="Ria French" class="w-full h-full object-cover" loading="lazy"> -->

      <!-- Avatar fallback (remove when photo available) -->
      <div class="w-16 h-16 rounded-full bg-teal/20 flex items-center justify-center">
        <span class="font-display font-bold text-xl text-teal">RF</span>
      </div>

      <!-- License badge — always visible -->
      <div class="absolute bottom-3 left-3 bg-river-deep/85 px-2 py-1 rounded-[2px]">
        <span class="font-body font-bold text-[0.5rem] uppercase tracking-[0.08em] text-teal">GUD-LIC-37359</span>
      </div>
    </div>

    <!-- Info -->
    <div class="p-4">
      <h3 class="font-display font-bold text-base text-ivory mb-2">Ria French</h3>
      <p class="font-body font-light text-body-sm text-ivory/60 leading-relaxed">
        Great Falls native. Former paramedic turned full-time guide.
        On the Trout Unlimited board for the Missouri River Flyfishers
        chapter &mdash; she fishes this river when she&rsquo;s not working it.
      </p>
    </div>

  </div>

</div>

<!-- Section CTA -->
<div class="mt-8 text-center">
  <a href="/meet-our-guides"
     class="font-body font-semibold text-eyebrow uppercase tracking-[0.08em] border border-ivory/35 text-ivory hover:bg-white/10 px-6 py-3 rounded-btn transition-colors inline-block">
    Meet All Seven Guides &rarr;
  </a>
</div>
```

---

### 6.7 Review Card

```html
<!-- 3-column grid wrapper -->
<div class="grid grid-cols-1 md:grid-cols-3 gap-4">

  <div class="bg-bone rounded-card p-6">
    <p class="text-teal text-sm tracking-[0.15em] mb-3">&#9733;&#9733;&#9733;&#9733;&#9733;</p>
    <blockquote class="font-display italic text-body text-river-deep leading-relaxed mb-3">
      &ldquo;[Review text from Matt&rsquo;s Google or Facebook]&rdquo;
    </blockquote>
    <p class="font-body font-bold text-label uppercase tracking-[0.06em] text-stone">
      &mdash; [Name] &middot; Google
    </p>
  </div>

</div>

<!-- Section CTA -->
<div class="mt-8 text-center">
  <a href="/read-our-reviews"
     class="font-body font-bold text-label uppercase tracking-[0.08em] text-teal hover:text-teal-dark transition-colors">
    Read All Reviews &rarr;
  </a>
</div>
```

---

### 6.8 Heroes Section (2-column)

```html
<section class="w-full py-section bg-bone">
  <div class="max-w-site mx-auto px-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-12 items-start">

      <!-- Left col: copy -->
      <div>
        <!-- Badge -->
        <div class="inline-flex items-center gap-2 border border-teal/30 bg-teal/10 rounded-[3px] px-3 py-2 mb-5">
          <span class="font-body font-bold text-label uppercase tracking-[0.1em] text-teal">Montana Heroes Program</span>
        </div>

        <p class="font-body font-bold text-eyebrow uppercase tracking-[0.15em] text-teal mb-3">
          Military &middot; Fire &middot; EMS &middot; Law Enforcement
        </p>
        <h2 class="font-display font-bold text-h2 text-river-deep leading-snug mb-4">
          A Full Day, at a Different Rate.
        </h2>
        <p class="font-body font-light text-body text-stone leading-relaxed mb-5">
          Some people don&rsquo;t take enough days off. Matt built this program for the ones
          who serve &mdash; active and retired military, National Guard, firefighters, search
          and rescue, law enforcement, and EMTs. A full day on the Missouri, gear and
          lunch included, at $500.
        </p>

        <!-- Price -->
        <p class="font-display font-bold text-[3rem] leading-none text-teal mb-1">$500</p>
        <p class="font-body font-light text-body-sm text-stone mb-5">full day trip, gear and lunch included</p>

        <!-- Eligibility list -->
        <ul class="mb-7 space-y-1">
          <li class="font-body font-light text-body-sm text-stone before:content-['—'] before:mr-2 before:text-teal">Active and retired military</li>
          <li class="font-body font-light text-body-sm text-stone before:content-['—'] before:mr-2 before:text-teal">National Guard</li>
          <li class="font-body font-light text-body-sm text-stone before:content-['—'] before:mr-2 before:text-teal">Firefighters and search and rescue</li>
          <li class="font-body font-light text-body-sm text-stone before:content-['—'] before:mr-2 before:text-teal">Law enforcement and EMT/Paramedics</li>
        </ul>

        <a href="/booking?rate=heroes"
           class="font-body font-bold text-eyebrow uppercase tracking-[0.1em] bg-teal hover:bg-teal-dark text-white px-7 py-3 rounded-btn transition-colors inline-block">
          Book the Heroes Rate
        </a>
      </div>

      <!-- Right col: pull quote panel -->
      <div class="bg-river-mid rounded-card p-8 flex flex-col justify-center min-h-[260px]">
        <blockquote class="font-display italic text-[1.35rem] text-ivory leading-snug mb-4">
          &ldquo;Taking people fishing is my passion.&rdquo;
        </blockquote>
        <p class="font-body font-light text-body-sm text-ivory/50 uppercase tracking-[0.08em]">
          &mdash; Matt Mohar, Owner
        </p>
      </div>

    </div>
  </div>
</section>
```

---

### 6.9 About Matt (2-column)

```html
<section class="w-full py-section bg-river-mid">
  <div class="max-w-site mx-auto px-6">
    <div class="grid grid-cols-1 md:grid-cols-[1fr_1.6fr] gap-10 items-start">

      <!-- Photo -->
      <div class="aspect-[3/4] bg-[#243d40] rounded-card overflow-hidden flex items-center justify-center">
        <!-- When photo available: -->
        <!-- <img src="/static/img/matt-mohar.webp" alt="Matt Mohar" class="w-full h-full object-cover" loading="lazy"> -->
        <div class="w-20 h-20 rounded-full bg-teal/20 flex items-center justify-center">
          <span class="font-display font-bold text-2xl text-teal">M</span>
        </div>
      </div>

      <!-- Copy -->
      <div>
        <p class="font-body font-bold text-eyebrow uppercase tracking-[0.15em] text-teal mb-3">The Outfitter</p>
        <h2 class="font-display font-bold text-h2 text-ivory leading-snug mb-1">Matt Mohar</h2>
        <p class="font-body font-bold text-label uppercase tracking-[0.1em] text-teal mb-6">Owner &middot; Lo Mo Outfitting</p>
        <p class="font-body font-light text-body text-ivory/70 leading-relaxed mb-4">
          Born and raised in Helena, Montana. I&rsquo;ve been fly fishing the Missouri River
          for more than 20 years and guiding on it for close to ten. This is home water
          &mdash; not a destination I fly into for the season.
        </p>
        <p class="font-body font-light text-body text-ivory/70 leading-relaxed">
          I built Lo Mo around guides who feel the same way. You won&rsquo;t get someone
          reading the water for the first time. You&rsquo;ll fish with people who are out
          here on their days off.
        </p>
      </div>

    </div>
  </div>
</section>
```

---

### 6.10 CTA Strip

```html
<div class="w-full bg-teal-dark py-10">
  <div class="max-w-site mx-auto px-6 flex flex-col md:flex-row items-start md:items-center justify-between gap-6">

    <div>
      <h2 class="font-display font-bold text-[1.5rem] text-white leading-snug mb-1">Ready to book your trip?</h2>
      <p class="font-body font-light text-body-sm text-white/70">Email to start. Call when you&rsquo;re ready to talk dates.</p>
    </div>

    <div class="flex flex-wrap items-center gap-4">
      <a href="tel:+1406XXXXXXX"
         class="font-display font-bold text-[1.2rem] text-white bg-white/15 px-4 py-2 rounded-[3px]">
        [PHONE]
      </a>
      <a href="mailto:matt@lomooutfitting.com"
         class="font-body font-light text-body-sm text-white/75 underline">
        matt@lomooutfitting.com
      </a>
      <a href="/booking"
         class="font-body font-bold text-eyebrow uppercase tracking-[0.1em] bg-river-deep hover:bg-river-mid text-ivory px-6 py-3 rounded-btn transition-colors">
        Book Online &rarr;
      </a>
    </div>

  </div>
</div>
```

---

### 6.11 Footer

```html
<footer class="w-full bg-river-deep">

  <div class="max-w-site mx-auto px-6 pt-12 pb-6">

    <!-- 3-column grid -->
    <div class="grid grid-cols-1 md:grid-cols-[1.5fr_1fr_1fr] gap-8 mb-10">

      <!-- Col 1: Brand -->
      <div>
        <p class="font-display font-bold text-base text-ivory mb-2">Lo Mo Outfitting LLC</p>
        <p class="font-body font-light text-body-sm text-ivory/40 leading-relaxed max-w-[200px]">
          Missouri River fly fishing out of Craig, Montana.
          Local guides. Real fish. No pretense.
        </p>
      </div>

      <!-- Col 2: Nav -->
      <div>
        <p class="font-body font-medium text-label uppercase tracking-[0.12em] text-ivory/30 mb-4">Navigation</p>
        <ul class="space-y-2">
          <li><a href="/trips"            class="font-body text-body-sm text-ivory/60 hover:text-ivory transition-colors">Trip Options</a></li>
          <li><a href="/meet-our-guides"  class="font-body text-body-sm text-ivory/60 hover:text-ivory transition-colors">Meet Our Guides</a></li>
          <li><a href="#heroes"           class="font-body text-body-sm text-ivory/60 hover:text-ivory transition-colors">Heroes Program</a></li>
          <li><a href="/read-our-reviews" class="font-body text-body-sm text-ivory/60 hover:text-ivory transition-colors">Read Our Reviews</a></li>
          <li><a href="/lo-mo-store"      class="font-body text-body-sm text-ivory/60 hover:text-ivory transition-colors">Lo Mo Store</a></li>
        </ul>
      </div>

      <!-- Col 3: Contact -->
      <div>
        <p class="font-body font-medium text-label uppercase tracking-[0.12em] text-ivory/30 mb-4">Contact</p>
        <ul class="space-y-2">
          <li><a href="tel:+1406XXXXXXX"              class="font-body text-body-sm text-ivory/60 hover:text-ivory transition-colors">[PHONE]</a></li>
          <li><a href="mailto:matt@lomooutfitting.com" class="font-body text-body-sm text-ivory/60 hover:text-ivory transition-colors">matt@lomooutfitting.com</a></li>
          <li><span                                    class="font-body text-body-sm text-ivory/60">Craig, Montana</span></li>
        </ul>
      </div>

    </div>

    <!-- Bottom bar -->
    <div class="border-t border-white/[0.08] pt-5 flex flex-col md:flex-row items-start md:items-center justify-between gap-2">
      <p class="font-body font-light text-xs text-ivory/30">&copy; 2025 Lo Mo Outfitting LLC &middot; All rights reserved.</p>
      <p class="font-body font-light text-xs text-ivory/30">Montana Outfitter License [OTF-XXXXX]</p>
      <p class="font-body font-light text-xs text-ivory/30">Site by <a href="https://fireflysoftware.dev" class="hover:text-ivory/50 transition-colors">Firefly Software</a></p>
    </div>

  </div>

</footer>
```

---

## 7. Buttons — All Variants

```html
<!-- Primary (teal) -->
<a href="#" class="font-body font-bold text-eyebrow uppercase tracking-[0.1em] bg-teal hover:bg-teal-dark text-white px-7 py-3 rounded-btn transition-colors inline-block">
  Book a Trip
</a>

<!-- Outline — light background -->
<a href="#" class="font-body font-semibold text-eyebrow uppercase tracking-[0.08em] border border-river-deep text-river-deep hover:bg-river-deep/5 px-6 py-3 rounded-btn transition-colors inline-block">
  Meet Our Guides
</a>

<!-- Outline — dark background -->
<a href="#" class="font-body font-semibold text-eyebrow uppercase tracking-[0.08em] border border-ivory/35 text-ivory hover:bg-white/10 px-6 py-3 rounded-btn transition-colors inline-block">
  Meet Our Guides
</a>

<!-- Dark solid (used in CTA strip) -->
<a href="#" class="font-body font-bold text-eyebrow uppercase tracking-[0.1em] bg-river-deep hover:bg-river-mid text-ivory px-6 py-3 rounded-btn transition-colors inline-block">
  Book Online &rarr;
</a>
```

---

## 8. Accessibility & Implementation Notes

- **Tap-to-call:** all phone links use `href="tel:+1406XXXXXXX"` — full E.164 format. Display text may be formatted `(406) XXX-XXXX`.
- **Images:** all non-hero images use `loading="lazy"`. Hero image uses `loading="eager"`. All images require descriptive `alt` text.
- **Photo formats:** serve WebP with JPEG fallback. Use `srcset` for responsive sizes.
- **Nav:** `sticky top-0 z-50`. Alpine.js controls the mobile drawer with `x-data="{ open: false }"` and `@click="open = !open"`.
- **License number:** store the outfitter license value in Go site config (`config.OutfitterLicense`), not hardcoded in templates.
- **Phone number:** store in Go site config (`config.PhoneDisplay` and `config.PhoneE164`), not hardcoded.
- **Placeholder brackets:** `[PHONE]`, `[OTF-XXXXX]`, and all review text are pending client verification. Do not ship with any bracket remaining.
- **Teal discipline:** teal appears in one role per section — either a CTA button or a set of stat numbers, not both. Do not add teal text links in a section that already has a teal CTA button.
- **Border on featured trip card:** the featured card uses `border-l-[3px] border-l-teal` — the left border overrides the default `border border-[#E0DBD2]`. Ensure both classes are present; Tailwind applies the more specific left border.

---

## 9. Placeholder Values — Do Not Ship

| Location            | Placeholder            | Source                     |
|---------------------|------------------------|----------------------------|
| Nav phone           | `[PHONE]`              | Matt Mohar                 |
| CTA strip phone     | `[PHONE]`              | Matt Mohar                 |
| Footer phone        | `[PHONE]`              | Matt Mohar                 |
| Footer license      | `[OTF-XXXXX]`          | Matt Mohar                 |
| Guide licenses      | `[VERIFY]` (6 guides)  | Matt Mohar                 |
| Review card 1–3     | `[Review text]`        | Matt's Google / Facebook   |
| Hero photo          | `hero-river.webp`      | Matt's photo library       |
| Matt's portrait     | avatar fallback        | Matt's photo library       |

---

*Lo Mo Outfitting · UI & Typography Guide · Firefly Software · fireflysoftware.dev*
*DS-002 Teal & Bone · Go · htmx · Alpine.js · Tailwind CSS v3*