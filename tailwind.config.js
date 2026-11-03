const { execSync } = require("node:child_process");

// #372: the kit HTML lives inside the Go module, not in web/templates, so
// Tailwind never sees its utility classes and a fresh app renders unstyled
// components. Resolve the module the same way Go does; a local replace
// (amarra-cais link) keeps working.
function shippedKitContent() {
  try {
    const out = execSync("go list -m -f {{.Dir}} github.com/puppe1990/amarra-cais", {
      encoding: "utf8",
      stdio: ["ignore", "pipe", "ignore"],
    });
    const dir = out.trim().replace(/\\/g, "/");
    return dir ? [dir + "/pkg/amarra/view/components/**/*.html"] : [];
  } catch {
    return [];
  }
}

/** @type {import("tailwindcss").Config} */
module.exports = {
  content: ["./web/templates/**/*.html", ...shippedKitContent()],
  safelist: [
    "cais-password-wrap",
    "cais-password-toggle",
    "cais-select-search",
    "cais-select-search-native",
    "cais-select-search-trigger",
    "cais-select-search-panel",
    "cais-select-search-input",
    "cais-select-search-list",
    "cais-select-search-option",
    "cais-select-search-label",
    "cais-select-search-chevron",
    "is-selected",
    "is-highlighted",
    "is-hidden",
    /* JS-only state classes: focus mode and the toast never appear in scanned
       HTML (sintaxe.js toggles them), so without the safelist Tailwind prunes
       their rules and the chrome never hides. */
    "amarra-focus",
    "amarra-toast",
    /* The runtime simulator picks one palette per stack in Go, so the class
       names never appear in scanned HTML. */
    "bg-amber-500",
    "bg-cyan-500",
    "bg-emerald-500",
    "bg-indigo-500",
    "border-amber-500/40",
    "border-cyan-500/40",
    "border-emerald-500/40",
    "border-indigo-500/40",
  ],
  theme: {
    extend: {
      colors: {
        /* CSS variables (input.css) so html.light can flip the palette (#258). */
        ink: "rgb(var(--amarra-ink) / <alpha-value>)",
        foam: "rgb(var(--amarra-foam) / <alpha-value>)",
        copper: "rgb(var(--amarra-copper) / <alpha-value>)",
        tide: "rgb(var(--amarra-tide) / <alpha-value>)",
        coral: "rgb(var(--amarra-coral) / <alpha-value>)",
        danger: "rgb(var(--amarra-danger) / <alpha-value>)",
        success: "rgb(var(--amarra-success) / <alpha-value>)",
      },
      fontFamily: {
        sans: ["Plus Jakarta Sans", "ui-sans-serif", "system-ui", "sans-serif"],
        serif: ["Fraunces", "Iowan Old Style", "Palatino Linotype", "Palatino", "Georgia", "serif"],
        display: [
          "Fraunces",
          "Iowan Old Style",
          "Palatino Linotype",
          "Palatino",
          "Georgia",
          "serif",
        ],
        mono: ["JetBrains Mono", "ui-monospace", "SFMono-Regular", "Menlo", "monospace"],
      },
      boxShadow: {
        "2xs": "0 1px 2px 0 rgb(0 0 0 / 0.05)",
        xs: "0 1px 2px 0 rgb(0 0 0 / 0.05)",
      },
    },
  },
  plugins: [],
};
