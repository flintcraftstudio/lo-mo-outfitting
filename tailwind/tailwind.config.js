/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/view/**/*.templ",
  ],
  theme: {
    extend: {
      colors: {
        "river-deep":  "#1A2E30",
        "river-mid":   "#2A4A4E",
        "teal":        "#3E9E98",
        "teal-dark":   "#2C7A75",
        "bone":        "#F4F1EB",
        "stone":       "#7A8E8F",
        "slate-text":  "#1A2E30",
        "ivory":       "#F4F1EB",
      },
      fontFamily: {
        display: ["Playfair Display", "Georgia", "serif"],
        body:    ["Source Sans 3", "system-ui", "sans-serif"],
      },
      fontSize: {
        "eyebrow":  ["0.6875rem", { lineHeight: "1", letterSpacing: "0.15em" }],
        "label":    ["0.5625rem", { lineHeight: "1", letterSpacing: "0.12em" }],
        "body-sm":  ["0.875rem",  { lineHeight: "1.65" }],
        "body":     ["1rem",      { lineHeight: "1.65" }],
        "price":    ["1.4rem",    { lineHeight: "1.2"  }],
        "stat":     ["2rem",      { lineHeight: "1.1"  }],
        "h3":       ["1.15rem",   { lineHeight: "1.3"  }],
        "h2":       ["2rem",      { lineHeight: "1.2"  }],
        "h1":       ["3rem",      { lineHeight: "1.1"  }],
        "h1-mob":   ["2rem",      { lineHeight: "1.1"  }],
      },
      spacing: {
        "section": "3.5rem",
      },
      borderRadius: {
        "card": "3px",
        "btn":  "2px",
      },
      maxWidth: {
        "site":     "1280px",
        "prose":    "520px",
        "hero-sub": "480px",
      },
    },
  },
  plugins: [],
}
