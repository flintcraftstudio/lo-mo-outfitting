# Lo Mo Outfitting — SEO Page Structure & Template Spec

**For coding agent use.**
**Stack:** Go · html/template · htmx · Alpine.js · Tailwind CSS (Teal & Bone, DS-002)

---

## Overview

Eight SEO-targeted interior pages across three URL categories. Two Go templates
cover all eight pages — Template A for river pages, Template B for location, fish,
and experience pages. Each page is driven by a Go struct passed at render time,
so the template logic is written once and the content is data-driven.

No CMS. Page data is defined in Go source as a map or slice of structs, compiled
into the binary. Adding a new page = adding a new struct entry and rebuilding.

---

## URL Structure

```
/rivers/missouri-river      → Template A
/rivers/bitterroot          → Template A

/locations/craig-montana    → Template B
/locations/helena-montana   → Template B

/fish/brown-trout           → Template B
/fish/rainbow-trout         → Template B

/experience/beginners       → Template B
/experience/dry-fly         → Template B
```

---

## Go Router

Register all eight routes explicitly. No wildcard routing — each URL is discrete.

```go
mux.HandleFunc("GET /rivers/{slug}",     handlers.RiverPage)
mux.HandleFunc("GET /locations/{slug}",  handlers.LocationPage)
mux.HandleFunc("GET /fish/{slug}",       handlers.FishPage)
mux.HandleFunc("GET /experience/{slug}", handlers.ExperiencePage)
```

Each handler resolves the slug to a page data struct, returns 404 if not found,
and renders the appropriate template.

---

## Go Page Data Structs

### Shared types

```go
type GuideRef struct {
    Name    string
    License string
    Bio     string   // 2–3 sentences, page-specific angle
    Photo   string   // path to /static/img/guides/{filename}.webp
    Initials string  // fallback avatar
}

type TripCard struct {
    Eyebrow string
    Title   string
    Price   string
    Details string
    URL     string
}

type RelatedPage struct {
    Label string
    URL   string
    Desc  string
}

type SEOMeta struct {
    Title       string  // <title> tag — 50–60 chars
    Description string  // meta description — 150–160 chars
    Canonical   string  // full canonical URL
}
```

### Template A — River page

```go
type RiverPage struct {
    SEO         SEOMeta
    HeroPhoto   string   // /static/img/river-{slug}-hero.webp
    HeroAlt     string
    Eyebrow     string   // e.g. "Guided float trips · Craig, Montana"
    Headline    string   // e.g. "The Missouri River"
    Subhead     string

    // Quick facts bar (4 items)
    Facts       [4]struct {
        Number string
        Label  string
    }

    // About section
    AboutBody   []string // slice of paragraphs

    // Seasonal calendar
    Seasons     []struct {
        Name    string   // e.g. "Spring · March–May"
        Body    string
        Hot     bool     // true = teal accent on this season
    }

    // Techniques section
    TechIntro   string
    Techniques  []struct {
        Name string
        Body string
    }

    // Featured guides (2–3)
    Guides      []GuideRef

    // Related pages
    Related     []RelatedPage

    // CTA
    CTAHeadline string
    CTASubhead  string
}
```

### Template B — Location / Fish / Experience page

```go
type ContentPage struct {
    SEO         SEOMeta
    HeroPhoto   string
    HeroAlt     string
    Eyebrow     string
    Headline    string
    Subhead     string

    // Page intro
    IntroBody   []string // paragraphs

    // Why Lo Mo section
    WhyHeadline string
    WhyPoints   []struct {
        Title string
        Body  string
    }

    // Trip options (2–3 cards, subset of full trip list)
    Trips       []TripCard

    // Featured guides (1–3)
    Guides      []GuideRef

    // Practical info (flexible key-value pairs)
    PracticalHeadline string
    PracticalItems    []struct {
        Label string
        Body  string
    }

    // Related pages
    Related     []RelatedPage

    // CTA
    CTAHeadline string
    CTASubhead  string
}
```

---

## Template A — River Page Sections

### Section 1: Hero

Same hero pattern as homepage. Dark overlay on river-specific photo.
Eyebrow, headline, subhead, single CTA button.

```
Background:  River Deep (#1A2E30) + photo overlay at 65% opacity
Eyebrow:     teal, uppercase, tracking
Headline:    Playfair Display 700, bone white, 3rem desktop
Subhead:     Source Sans 3, bone/70%, max-w-hero-sub
CTA:         "Book a Float Trip" → /booking
```

### Section 2: River Quick Facts Bar

4-column stats bar. Same component as homepage stats bar.
River-specific numbers — not the generic Lo Mo stats.

```
Background:  River Deep, border-top border-white/6
Numbers:     Playfair Display 700 text-stat text-teal
Labels:      Source Sans 3 600 9px uppercase text-ivory/40
```

### Section 3: About This River

Bone background. Eyebrow + headline left-aligned. Body copy in 2-column prose
grid on desktop, single column on mobile. 3–4 paragraphs covering:
- What makes this river distinctive
- Where it is and who fishes it
- Lo Mo's specific connection to it

```
Background:  Bone (#F4F1EB)
Eyebrow:     teal uppercase
Headline:    Playfair Display 700 text-river-deep
Body:        Source Sans 3 300 text-stone leading-relaxed
Grid:        md:grid-cols-2 gap-10
```

### Section 4: Seasonal Fishing Calendar

White background. Horizontal row of season cards — 4 seasons, each with name,
body copy, and a "hot season" teal accent for peak months.
On mobile: vertical stack.

```
Background:  White
Card bg:     Bone for standard seasons, River Mid for hot season
Season name: Playfair Display 700 — bone on dark card, river-deep on bone card
Body:        Source Sans 3 300 11px
Hot badge:   "Best fishing" in teal, shown only on hot season card(s)
```

### Section 5: What to Expect (Techniques)

Bone background. Eyebrow + headline. 3-column grid of technique cards.
Each card: technique name as heading, 2–3 sentence description.
Keeps language approachable — not a gear manual.

```
Background:  Bone
Card bg:     White, border border-[#E0DBD2], rounded-card
Heading:     Playfair Display 700 text-h3 text-river-deep
Body:        Source Sans 3 300 text-stone
```

### Section 6: Featured Guides

River Deep background. 2–3 guide cards using the same guide card component
from the homepage. Bio copy is page-specific — angle it toward this river,
not a generic bio.

```
Background:  River Deep
Cards:       River Mid bg, 2–3 col grid
Photo:       3:4 aspect, full-bleed, lazy loaded
License:     teal badge, bottom-left of photo
Name:        Playfair Display 700 text-ivory
Bio:         Source Sans 3 300 text-ivory/60
```

### Section 7: Related Pages

Bone background. "Also worth reading" eyebrow. 3-column card grid linking to
related pages — other rivers, fish species, experience types. Each card:
linked title + 1-sentence description.

```
Background:  Bone
Cards:       White bg, border, hover border-teal transition
Title:       Playfair Display 700 text-river-deep, linked
Desc:        Source Sans 3 300 text-stone
```

### Sections 8–9: CTA Strip + Footer

Identical to homepage. Shared partials — `templates/partials/cta_strip.html`
and `templates/partials/footer.html`. Pass phone and email from site config.

---

## Template B — Location / Fish / Experience Page Sections

### Section 1: Hero

Same hero pattern. Photo is location/fish/experience specific.
CTA button copy adapts to page type:
- Location pages: "Book a Trip from [City]"
- Fish pages: "Book a [Fish] Float"
- Experience pages: "Book Your Trip"

### Section 2: Page Intro + Context

Bone background. Left-aligned eyebrow, headline, 2–3 paragraph intro.
This is the most SEO-critical section — the intro body contains the target
keyword naturally in the first paragraph, not stuffed.

```
Background:  Bone
Layout:      Single column, max-w-prose on body copy
```

### Section 3: Why Lo Mo for This

White background. Eyebrow + headline. 3-column grid of "why" points —
each point has a title and 2–3 sentence explanation. Content is specific
to the page's angle (location proximity, fish expertise, teaching approach).
Not generic "we're the best" claims — specific and verifiable.

```
Background:  White
Card bg:     Bone, no border, rounded-card
Title:       Playfair Display 700 text-h3 text-river-deep
Body:        Source Sans 3 300 text-stone
```

### Section 4: Relevant Trip Options

Bone background. 2–3 trip cards from the full trip list — whichever options
are most relevant to this page's audience. Use the same trip card component
as the homepage. Add a "See all trip options" text link below the grid.

```
Background:  Bone
Cards:       Same as homepage trip card component
Link:        "See all trip options →" → /trips, teal text link
```

### Section 5: Featured Guide(s)

River Deep background. 1–3 guide cards. Guide selection and bio angle
are page-specific:
- Location pages: guides who know the local logistics for that city
- Fish pages: guides who specialize in that species
- Experience pages: guides known for teaching or technical expertise

### Section 6: Practical Info

Bone background. Eyebrow + headline. Flexible key-value list covering
logistics specific to this page:
- Location pages: drive times, where to meet, lodging suggestions
- Fish pages: best months, recommended techniques, what's provided
- Experience pages: what to expect, how to prepare, what's included

```
Background:  Bone
Layout:      2-column definition list, dl/dt/dd or grid
Label (dt):  Source Sans 3 600 text-river-deep
Value (dd):  Source Sans 3 300 text-stone
```

### Section 7: Related Pages

White background. Same related pages component as Template A.

### Sections 8–9: CTA Strip + Footer

Identical shared partials.

---

## SEO Requirements Per Page

Every page must have:

```html
<title>{{ .SEO.Title }}</title>
<meta name="description" content="{{ .SEO.Description }}">
<link rel="canonical" href="{{ .SEO.Canonical }}">

<!-- Open Graph -->
<meta property="og:title" content="{{ .SEO.Title }}">
<meta property="og:description" content="{{ .SEO.Description }}">
<meta property="og:url" content="{{ .SEO.Canonical }}">
<meta property="og:type" content="website">
<meta property="og:image" content="{{ .HeroPhoto | absURL }}">

<!-- Structured data — LocalBusiness on every page -->
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "LocalBusiness",
  "name": "Lo Mo Outfitting",
  "url": "https://lomooutfitting.com",
  "telephone": "[PHONE]",
  "email": "matt@lomooutfitting.com",
  "address": {
    "@type": "PostalAddress",
    "addressLocality": "Craig",
    "addressRegion": "MT",
    "addressCountry": "US"
  },
  "description": "Guided fly fishing on the Missouri River out of Craig, Montana."
}
</script>
```

---

## Page-by-Page Content Spec

---

### /rivers/missouri-river

**SEO**
- Title: `Missouri River Fly Fishing Guide · Lo Mo Outfitting · Craig, MT`
- Description: `Guided fly fishing on the Missouri River near Craig, Montana. Born-here guides, drift boat floats, all gear included. Full day from $675.`
- Canonical: `https://lomooutfitting.com/rivers/missouri-river`

**Hero**
- Eyebrow: `Guided Float Trips · Craig, Montana`
- Headline: `The Missouri River`
- Subhead: `Thirty-five miles of blue-ribbon tailwater below Holter Dam. Wild rainbows and browns averaging 14–22 inches. This is the river Matt Mohar grew up on.`
- Photo: `/static/img/river-missouri-hero.webp`
- CTA: `Book a Float Trip`

**Facts bar**
- 35mi / Blue-Ribbon Water
- 5K+ / Trout per Mile
- Year-Round / Fishable
- Apr–Oct / Peak Season

**About body** (2 paragraphs)
- P1: What the Mo is — tailwater below Holter Dam, consistent flows, prolific hatches, 35 miles from dam to Cascade. The standard facts stated specifically.
- P2: Lo Mo's relationship to it — Matt has fished it for 20+ years. Guides live here. Not visiting guides who come up from Bozeman for spring.

**Seasons**
- Spring (Mar–May): BWO hatches, post-spawn rainbows feeding hard, high water by mid-May. Hot = true.
- Summer (Jun–Aug): PMD and caddis blanket hatches, dry fly paradise Jun–Jul, trico mornings in Aug.
- Fall (Sep–Nov): BWO returns, browns going aggressive, streamer season, fewest crowds. Hot = true.
- Winter (Dec–Feb): Midge fishing, nymphing deep. Fishable for the dedicated.

**Techniques**
- Nymphing: The most consistent producer. Sowbugs, midge larvae, PMD nymphs in the deep runs. Dead drift.
- Dry Fly: The Mo's reputation. PMDs, caddis, tricos, BWOs. Technical — fish are selective. Best Jun–Jul and Sep–Oct.
- Streamers: Big browns chase big flies. Cloudy days, low light. Best fall and early spring.

**Guides**: Ria French, Colter Day, Rain Keating
- Bio angles: Ria — TU conservation, knows the hatches; Colter — lifetime local; Rain — came for the river, stayed.

**Related pages**
- /fish/brown-trout — "Chasing Missouri browns"
- /fish/rainbow-trout — "The Mo's wild rainbows"
- /experience/dry-fly — "Dry fly fishing on the Missouri"
- /locations/craig-montana — "Fishing out of Craig"

**CTA**
- Headline: `Ready to fish the Missouri?`
- Subhead: `Matt will match you with the right guide for your dates and experience level.`

---

### /rivers/bitterroot

**SEO**
- Title: `Bitterroot River Fly Fishing Guide · Lo Mo Outfitting`
- Description: `Guided fly fishing on the Bitterroot River with Dave Buck — a Conner, Montana local who guides the Root year-round. Cutthroat, rainbow, and brown trout.`
- Canonical: `https://lomooutfitting.com/rivers/bitterroot`

**Hero**
- Eyebrow: `Guided Float Trips · Bitterroot Valley, Montana`
- Headline: `The Bitterroot River`
- Subhead: `Eighty miles of freestone river through the Bitterroot Valley. Westslope cutthroat, brown trout, and dry fly fishing that runs from the Skwala hatch in March through BWOs in October.`
- Photo: `/static/img/river-bitterroot-hero.webp`
- CTA: `Book a Bitterroot Float`

**Facts bar**
- 80mi / Of Fishable Water
- 3 Species / Cutthroat · Rainbow · Brown
- Mar–Oct / Prime Season
- Skwala / First Hatch of Spring

**About body**
- P1: What the Bitterroot is — freestone river, Sapphire and Bitterroot Mountains, flows north through the valley to the Clark Fork near Missoula.
- P2: Dave Buck's connection — grew up fishing central Oregon, guided western Washington, moved to Conner. The Bitterroot is his home water. He guides it on his days off.

**Seasons**
- Spring (Mar–May): Skwala stonefly hatch starts in March — first good dry fly fishing of the year in Montana. Hot = true.
- Summer (Jun–Aug): Salmonfly, PMD, caddis. Runoff in May-Jun, then excellent conditions through August.
- Fall (Sep–Oct): Trico, Hecuba, mahogany duns, BWO. Browns aggressive. Best streamer fishing. Hot = true.
- Winter (Nov–Feb): Nymphing in deep winter water. Streamer fishing for the dedicated.

**Techniques**
- Dry Fly: The Bitterroot's reputation. Prolific hatches, fish that eat on top. Can be demanding — selective fish in clear water.
- Nymphing: Heavy and deep. Stonefly and mayfly nymphs. Essential when fish aren't looking up.
- Streamers: Small sparkle minnows and kreelex most of the season. Big articulated patterns in spring and fall.

**Guides**: Dave Buck
- Bio angle: Grew up fishing Oregon, moved to Conner specifically for the Bitterroot. This is his home river — he guides it when he isn't guiding it.

**Related pages**
- /rivers/missouri-river — "Also fish the Missouri"
- /fish/brown-trout — "Chasing Bitterroot browns"
- /locations/helena-montana — "Driving from Helena"

**CTA**
- Headline: `Book a day on the Bitterroot.`
- Subhead: `Dave Buck guides this river year-round. Contact Matt to check availability.`

---

### /locations/craig-montana

**SEO**
- Title: `Fly Fishing Guide in Craig, Montana · Lo Mo Outfitting`
- Description: `Local fly fishing guides based out of Craig, MT. Drift boat floats on the Missouri River from Holter Dam to Cascade. Full day from $675, gear included.`
- Canonical: `https://lomooutfitting.com/locations/craig-montana`

**Hero**
- Eyebrow: `Craig, Montana · Missouri River`
- Headline: `Fishing out of Craig`
- Subhead: `Craig sits on the Missouri River between Helena and Great Falls. It's a small town built around fly fishing. Lo Mo guides are based here — not visiting.`
- Photo: `/static/img/location-craig-hero.webp`
- CTA: `Book a Craig Float`

**Intro body**
- P1: Craig is a town of fewer than 100 people on the banks of the Missouri River, 45 minutes north of Helena. It has four fly shops, several lodges, and nothing else. If you're coming to fish the Missouri, Craig is your base.
- P2: Lo Mo Outfitting operates out of Craig. Matt Mohar is a Helena native who has fished this stretch of river for over 20 years. The guides live here or close to it. You won't be handed off to a guide who drove three hours to meet you.

**Why Lo Mo**
- Local crew: Guides who live in and around Craig, not seasonal guides from Bozeman or Missoula.
- Right guide for you: Matt matches each angler with the right guide based on experience level, not whoever's available.
- No logistics gap: When your lodging is five minutes from the put-in, the day starts right.

**Trip options**: Full Day Single Boat, Full Day Multiple Boats, Heroes Program

**Guides**: Matt Mohar, Colter Day, Rain Keating
- Bio angles: Matt — Helena native, built Lo Mo around Craig; Colter — lifetime Montana local; Rain — moved here specifically for the Mo.

**Practical info**
- Getting here: Craig is on I-15, Exit 234, between Helena and Great Falls. About 45 minutes from Helena, 55 minutes from Great Falls.
- Where to stay: Craig has several lodges and cabins on or near the river. Wolf Creek (15 min south) has additional options. Helena is the closest city with full hotel options.
- Where to meet: Matt will confirm your meeting location before your trip — typically a fishing access site near Craig.
- What to bring: Layers appropriate for weather, polarized sunglasses, hat, closed-toe shoes. All fishing gear is provided.

**Related pages**
- /rivers/missouri-river — "About the Missouri River"
- /locations/helena-montana — "Coming from Helena"
- /experience/beginners — "First time fly fishing"

**CTA**
- Headline: `Book a Craig float.`
- Subhead: `Matt will confirm your dates and send a deposit link. No card required to request.`

---

### /locations/helena-montana

**SEO**
- Title: `Fly Fishing Guide near Helena, Montana · Lo Mo Outfitting`
- Description: `The Missouri River is 45 minutes from Helena. Lo Mo Outfitting runs guided drift boat trips out of Craig, MT — full day from $675, all gear included.`
- Canonical: `https://lomooutfitting.com/locations/helena-montana`

**Hero**
- Eyebrow: `45 Minutes from Helena · Missouri River`
- Headline: `Helena's closest world-class fishery`
- Subhead: `Matt Mohar grew up in Helena. The Missouri River is 45 minutes north. Lo Mo Outfitting is the guide service built by a Helena local, for anglers coming from Helena.`
- Photo: `/static/img/location-helena-hero.webp`
- CTA: `Book a Trip from Helena`

**Intro body**
- P1: If you're in Helena and want to fly fish, the Missouri River below Holter Dam is your answer. It's 35 miles of blue-ribbon tailwater — one of the most consistent trout fisheries in North America — and it's 45 minutes from downtown Helena on I-15.
- P2: Matt Mohar started Lo Mo Outfitting as a Helena native who has fished this river his whole life. He knows the water, he knows the guides, and he knows the drive.

**Why Lo Mo**
- Born here: Matt is from Helena. Not a guide service that relocated to Montana.
- Short drive, world-class water: 45 minutes on I-15, direct route with no navigation required.
- Early starts available: Helena-based anglers can be on the water by 8am without an overnight stay.

**Trip options**: Full Day Single Boat, Half Day Single Boat, Heroes Program

**Guides**: Matt Mohar
- Bio angle: Helena native, 20+ years on the Missouri. This is home water.

**Practical info**
- Drive time: 45 minutes Helena to Craig via I-15 North. Exit 234.
- Early start: Typical full day launches 7:30–9am depending on season. Matt confirms start time before your trip.
- Overnight not required: A full day trip from Helena is a day trip. Shuttle returns you to your vehicle.
- Hotels in Helena: If coming from out of state and using Helena as a base, the city has full hotel options within 45 minutes of the fishing.

**Related pages**
- /rivers/missouri-river — "About the Missouri River"
- /locations/craig-montana — "Fishing out of Craig"
- /experience/beginners — "First time on a fly rod"

**CTA**
- Headline: `Book a trip from Helena.`
- Subhead: `45 minutes to the put-in. Matt handles the rest.`

---

### /fish/brown-trout

**SEO**
- Title: `Guided Brown Trout Fly Fishing Montana · Lo Mo Outfitting`
- Description: `Chase Missouri River brown trout with Lo Mo Outfitting. Wild browns averaging 16–22 inches, with 2-footers a real possibility. Drift boat floats out of Craig, MT.`
- Canonical: `https://lomooutfitting.com/fish/brown-trout`

**Hero**
- Eyebrow: `Missouri River · Brown Trout`
- Headline: `Chasing Missouri browns`
- Subhead: `Wild brown trout averaging 16–22 inches. Some pushing 24. The Missouri River has one of the highest densities of large wild trout in the country — and the browns are the ones that keep guides up at night.`
- Photo: `/static/img/fish-brown-trout-hero.webp` (guide or angler holding a large brown)
- CTA: `Book a Brown Trout Float`

**Intro body**
- P1: The Missouri River brown trout are not stocked. They're wild, they're educated, and they're large. The average fish runs 16–22 inches. A two-footer is not a trophy story — it's a normal September on the Mo.
- P2: Browns favor structure — undercut banks, logjams, drop-offs. They eat streamers aggressively in low light and become technical dry fly targets during hatches. Fishing for them specifically is a different day than fishing for numbers.

**Why Lo Mo**
- Guides who target browns: Not every guide on the Missouri focuses on big fish. Lo Mo's guides know where the large browns hold and how to approach them.
- Streamer expertise: Sam Botz and Rain Keating specialize in streamer fishing for aggressive browns in fall.
- The right season: October browns on the Missouri are something specific. Matt will tell you the truth about when to come.

**Trip options**: Full Day Single Boat, Full Day Multiple Boats

**Guides**: Ria French, Sam Botz, Rain Keating
- Bio angles: Ria — "chasing brown trout" is literally in her existing bio; Sam — streamer obsessive, mice at night; Rain — calls the Mo home, knows where the big browns live.

**Practical info**
- Best months: September and October for streamer fishing and aggressive pre-spawn browns. June–July for browns rising to PMD and caddis hatches.
- Techniques: Streamers in low light and overcast conditions. Dry flies during hatches. Nymphing produces throughout.
- Flies: Articulated streamers (fall), sparkle minnows, PMD and caddis patterns (summer), sowbugs and midge nymphs year-round.
- What's provided: All flies, rods, reels, terminal tackle. Shoreline lunch on full day trips.

**Related pages**
- /rivers/missouri-river — "About the Missouri River"
- /fish/rainbow-trout — "Also on the Mo: wild rainbows"
- /experience/dry-fly — "Matching the hatch for browns"

**CTA**
- Headline: `Come for the browns.`
- Subhead: `Matt will put you with a guide who knows where they live.`

---

### /fish/rainbow-trout

**SEO**
- Title: `Guided Rainbow Trout Fly Fishing Montana · Lo Mo Outfitting`
- Description: `Wild Missouri River rainbow trout — 14–20 inches, 5,000+ per mile. Guided drift boat trips out of Craig, MT with Lo Mo Outfitting. From $675.`
- Canonical: `https://lomooutfitting.com/fish/rainbow-trout`

**Hero**
- Eyebrow: `Missouri River · Rainbow Trout`
- Headline: `Wild Missouri rainbows`
- Subhead: `The Missouri River runs 3:1 rainbows to browns. Wild fish, 14–20 inches, in water that has never seen a hatchery truck. This is what 5,000 trout per mile looks like.`
- Photo: `/static/img/fish-rainbow-trout-hero.webp`
- CTA: `Book a Rainbow Float`

**Intro body**
- P1: Rainbows are the Mo's most numerous wild trout. They outnumber browns roughly three to one and average 14–20 inches. They're the fish that make first-time visitors question why they waited so long to come to Montana.
- P2: In spring, post-spawn rainbows feed voraciously on nymphs. In summer, they rise to PMD and caddis hatches. In fall, they nymph hard before the water cools. There is no bad season for Missouri River rainbows.

**Why Lo Mo**
- Year-round fishery: Rainbows are catchable in every month. Lo Mo guides know how to fish them in all conditions.
- Right technique for the day: Nymphing for numbers, dry flies for the experience, streamers when conditions call for it.
- Great for beginners: Rainbow trout on the nymph rig is an achievable goal for a first-time fly angler. Guides teach as they fish.

**Trip options**: Full Day Single Boat, Half Day Single Boat, Early Season Full Day

**Guides**: Colter Day, Andrew Osborn, Dylan Huseby
- Bio angles: Colter — all skill levels, patient guide; Andrew — newer to guiding, good teacher; Dylan — loves teaching the drift, passion for moving water.

**Practical info**
- Best months: April–June for numbers of post-spawn fish. June–July for dry fly rainbows. Year-round for consistent nymphing.
- Techniques: Dead-drift nymphing is the most reliable. Dry fly during PMD, caddis, trico, BWO hatches.
- Flies: Sowbugs, midge larvae, PMD nymphs. Caddis and PMD dries in summer. Trico spinner falls in August mornings.
- Skill level: Rainbows on a nymph rig are an appropriate goal for complete beginners through advanced anglers.

**Related pages**
- /rivers/missouri-river — "About the Missouri River"
- /fish/brown-trout — "Also on the Mo: wild browns"
- /experience/beginners — "First time fly fishing"

**CTA**
- Headline: `Five thousand trout per mile.`
- Subhead: `Book a float and find out what that means.`

---

### /experience/beginners

**SEO**
- Title: `Beginner Fly Fishing Guide Montana · Lo Mo Outfitting`
- Description: `Never held a fly rod? Lo Mo Outfitting's guides teach while they fish. Drift boat float trips on the Missouri River near Craig, MT. All gear provided.`
- Canonical: `https://lomooutfitting.com/experience/beginners`

**Hero**
- Eyebrow: `First Time on the Water`
- Headline: `Your first fly fishing trip`
- Subhead: `No experience needed. Lo Mo's guides have put complete beginners on wild Missouri River trout on their first cast. All gear is provided. The teaching is part of the job.`
- Photo: `/static/img/experience-beginners-hero.webp`
- CTA: `Book Your First Trip`

**Intro body**
- P1: The Missouri River below Holter Dam is one of the best places in the country to learn fly fishing. The fish are abundant, the water is manageable from a drift boat, and the guides have spent years teaching the cast and the drift to people who have never fished before.
- P2: Lo Mo's booking form asks about experience level so Matt can match you with the right guide for your party. A group of beginners gets a guide who prioritizes teaching. An experienced angler gets a guide who will push them toward harder water.

**Why Lo Mo**
- Guides who teach: Not all guides have patience for beginners. Lo Mo's guides know that a first-timer who catches their first trout will book again next year.
- All gear provided: Rods, reels, flies, terminal tackle. The only things you need to bring are appropriate clothing and a fishing license.
- The drift boat does the work: You're not wading unfamiliar water. You're fishing from a stable platform while the guide positions you on fish.

**Trip options**: Full Day Single Boat, Half Day Single Boat

**Guides**: Colter Day, Andrew Osborn, Ria French
- Bio angles: Colter — helps anglers of all skill levels learn and progress; Andrew — remembers what it's like to be a beginner; Ria — 15 years teaching as a paramedic transferred to teaching on the water.

**Practical info**
- What's included: Fly rod and reel, flies, all terminal tackle, shuttle transportation. Full day includes shoreline lunch and cold drinks.
- What to bring: Layers for weather changes, polarized sunglasses (provided if you don't have them — ask Matt), hat, closed-toe shoes or water shoes. Waders are not required from a drift boat.
- Montana fishing license: Required for all anglers over 12. Purchase at fwp.mt.gov before your trip date.
- Youth anglers: Kids 12 and under fish free on the license. Children 5 and under cannot be accommodated for safety reasons.
- Managing expectations: The Mo is a prolific fishery but fly fishing has a learning curve. Most beginners catch fish. The guide's job is to make sure you do.

**Related pages**
- /rivers/missouri-river — "About the Missouri River"
- /experience/dry-fly — "Ready for something technical"
- /fish/rainbow-trout — "The fish you're likely to catch"

**CTA**
- Headline: `Everyone starts somewhere.`
- Subhead: `Book a trip and tell Matt it's your first time. He'll handle the rest.`

---

### /experience/dry-fly

**SEO**
- Title: `Dry Fly Fishing Missouri River Guide · Lo Mo Outfitting`
- Description: `Technical dry fly fishing on the Missouri River near Craig, MT. PMD, caddis, and trico hatches with Lo Mo Outfitting's guides. Drift boat floats from $675.`
- Canonical: `https://lomooutfitting.com/experience/dry-fly`

**Hero**
- Eyebrow: `Missouri River · Technical Dry Fly`
- Headline: `Matching the hatch on the Mo`
- Subhead: `The Missouri River is one of the finest dry fly fisheries in the American West. Blanket PMD and caddis hatches, trico spinner falls, BWO emergences. Wild trout rising to flies you tie on yourself. This is what the Mo is for.`
- Photo: `/static/img/experience-dry-fly-hero.webp`
- CTA: `Book a Dry Fly Float`

**Intro body**
- P1: The Missouri River below Holter Dam is a dry fly angler's paradise for most of the year. Consistent flows and fertile water produce prolific hatches of PMDs, caddis, tricos, BWOs, and more. Fish rise reliably — the challenge is presentation, not finding fish.
- P2: Dry fly fishing on the Mo is technical. The trout are educated and selective. Getting a good drift over a rising fish requires the right fly, the right angle, and the right presentation. Lo Mo's guides know this water hatch by hatch and will keep you on fish.

**Why Lo Mo**
- Hatch knowledge: Guides who fish the Mo year-round know when the PMDs come off, where the trico spinner falls happen, and which runs fish best in which conditions.
- Positioning: A drift boat guide's primary job on a dry fly day is positioning — getting you the right angle on rising fish. This takes local knowledge.
- Technical experience: The Mo rewards anglers who refine their presentation. Lo Mo's guides will tell you what's wrong with your drift and how to fix it.

**Trip options**: Full Day Single Boat, Full Day Multiple Boats

**Guides**: Sam Botz, Rain Keating, Dave Buck
- Bio angles: Sam — trico mornings, mice at night, technical dry fly obsessive; Rain — his days off on the Mo look exactly like his days on; Dave — Bitterroot's dry fly reputation, brings that expertise to the Mo.

**Practical info**
- Best months: June–July for PMD and caddis blanket hatches. August for trico spinner falls (early mornings). September–October for BWO emergences.
- Key hatches: PMD (Pale Morning Dun) · Caddis · Trico · Blue-Winged Olive (BWO) · Brown Drake (late June evenings)
- What makes it technical: Missouri River fish see thousands of fly presentations. Long, accurate drifts with appropriate patterns. Drag-free floats. Accurate tippet sizing. Guides will help with all of it.
- Still worth fishing on nymph days: If hatches aren't happening, guides switch to nymphing. You will still catch fish.

**Related pages**
- /rivers/missouri-river — "About the Missouri River"
- /fish/brown-trout — "Browns rising to PMDs"
- /fish/rainbow-trout — "Rainbows during caddis hatches"

**CTA**
- Headline: `The hatches are real. The fish are selective.`
- Subhead: `Book a dry fly float and come ready to refine your presentation.`

---

## Go Template Implementation Notes

- Both templates extend `templates/base.html` which provides the `<head>`,
  nav, and closing tags.
- SEO meta tags are in a `templates/partials/seo.html` partial — rendered
  once per page from the `SEOMeta` struct.
- Structured data JSON-LD is in `templates/partials/structured_data.html`.
  Phone number and address come from site config, not page data.
- The guide card component is `templates/partials/guide_card.html`.
  Accepts a `GuideRef` struct. Used on homepage, guides page, and all SEO pages.
- The trip card component is `templates/partials/trip_card.html`.
  Used on homepage and Template B pages.
- The related pages component is `templates/partials/related_pages.html`.
  Accepts `[]RelatedPage`.
- CTA strip and footer are `templates/partials/cta_strip.html` and
  `templates/partials/footer.html`. Phone and email from site config.
- All page data structs live in `internal/pages/data.go` as package-level vars.
  Handlers look up by slug from a `map[string]RiverPage` or `map[string]ContentPage`.
- Image paths follow the convention `/static/img/{category}-{slug}-{section}.webp`.
  All images are placeholders until Matt provides assets. Use the teal-initial
  avatar pattern from the guide card as the photo fallback pattern.

---

*Lo Mo Outfitting · SEO Page Spec · Firefly Software · fireflysoftware.dev*