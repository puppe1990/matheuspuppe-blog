package db

import "github.com/puppe1990/matheuspuppe-blog/internal/models"

// amarraEssayEn é a edição inglesa do ensaio sobre o Amarra-cais.
func amarraEssayEn() englishEdition {
	return englishEdition{
		Slug:     "why-i-built-amarra-cais",
		Title:    "Why I built Amarra-cais",
		Subtitle: "A fork of Cais on September 5, 687 commits and 1.148 tests later: the bet that a mini app does not need a SPA, and what it costs.",
		Lead:     "Amarra-cais came from a specific annoyance: every mini app of mine carried a whole SPA to show a table and a form. On September 5, 2026 I imported Cais (Inertia + Svelte) as the starting point and began an HTML-first fork. Router, session, CSRF, job queue and SQLite stayed. The client that mounts screens left.",
		Tags:     []string{"Amarra-cais", "Go", "Architecture", "Frameworks", "HTML"},
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"The start was a chore commit: import Cais v0.11 and rename the module. The work that followed was cutting per layer: Vite, Svelte, Inertia and the duplicated validation.",
				"The bet: the handler calls view.Write, Drive morphs #amarra-main, and sort and filter are GET. The app does not build state on the client; it asks for the page again.",
				"The bill is paid by generated boredom: 23 kit components, g resource and g auth in the CLI, and a doctor that turns every bug into a check (1.148 Go test functions, 203 JS tests, 206 issues).",
			}},
			{Type: "heading2", Text: "What stayed and what left"},
			{Type: "paragraph", Text: "Cais already had what I did not want to rewrite: router, session, CSRF, job queue and SQLite in a single file, no CGO. Amarra kept that base and changed the layer above it. The project README sums it up in three sentences: the browser does not mount a SPA, handlers call view.Write, Drive morphs #amarra-main. There is no Vite, Svelte or Inertia in a generated app."},
			{
				Type:        "quote",
				Text:        "The browser does not mount a SPA. Handlers call `view.Write`. Drive morphs `#amarra-main`. There is no Vite, Svelte, or Inertia in generated apps.",
				AuthorQuote: "Amarra-cais README",
				SourceURL:   "https://github.com/puppe1990/amarra-cais/blob/93c44c3/README.md",
				SourceLabel: "README.md",
			},
			{Type: "paragraph", Text: "The price was clear from the start: migrating a Cais app is a UI rewrite, and that is written in the docs (docs/migrate-inertia.md). There was no automatic converter. The choice was between keeping two runtimes forever and paying for the rewrite once."},
			{Type: "heading2", Text: "The trade: a morph instead of a runtime"},
			{Type: "paragraph", Text: "Drive intercepts same-origin clicks and submits by default and replaces only #amarra-main. The rest of the HTML comes from the server: the layout delivers the shell, the page fills the content block. The docs are explicit about what that buys: first load, curl and crawlers get the whole page; navigation becomes a morph. There is no attribute to turn Drive on, it is already on; whoever needs a full visit adds data-amarra-skip."},
			{Type: "paragraph", Text: "Instead of Alpine, Stimulus or HTMX, the framework ships the kit (23 components: form, input, table, modal, tabs, filters, stat, tooltip) and hooks in a single 117 KB script. A generated app has no bundler. Sort and filter are GET (?q=&sort=), the server re-renders, and the component that needs state (dropdown, dialog, tabs) is a native hook, not a third-party plugin. Streaming fits the same contract: amarra/stream writes named operations (append, replace, morph, remove, toast) with an id target, and the live hub uses WebSocket for the same shape."},
			{Type: "heading2", Text: "Boring parts must be generated"},
			{Type: "paragraph", Text: "The CLI is Rails-style on purpose: new, install, dev, test, routes, db, doctor, pwa, link/unlink and the generators (resource, auth, job, stream, component). g resource emits the whole kit, not a sample CRUD with loose classes. The layout is predictable: generators mark points in the code (// cais:seeds, // cais:live-views) and the upgrade knows where to write."},
			{
				Type:     "code",
				Language: "go",
				Caption:  "A whole handler: the layout owns the shell, the page only fills content",
				Code: `view.Write(w, r, h.views, view.Page{
  Layout: "app",
  Name:   "dashboard",
  Data: map[string]any{
    "Title":     h.catalog.T("dashboard.title"),
    "Site":      meta.ForRequest(h.site, r),
    "CSRFToken": csrf.TokenFromRequest(r),
    "Flash":     flashMsg,
  },
}, h.cfg)`,
				SourceURL:   "https://github.com/puppe1990/amarra-cais/blob/93c44c3/website/src/content/docs/docs/how-to/pages-and-views.md#L30-L39",
				SourceLabel: "website/src/content/docs/docs/how-to/pages-and-views.md:30",
			},
			{Type: "heading2", Text: "Every bug became a check"},
			{Type: "paragraph", Text: "The doctor exists because the bugs happened. Each check was born from an incident with a number: install exiting 0 with broken CSS (#54), styles.css missing after a clone (#141), the vendored bundle drifting from go.mod (#320). The comment stays in the code with the issue number, so the next person knows why that error is fatal."},
			{
				Type:     "code",
				Language: "go",
				Caption:  "Install fails loud: exiting 0 here left people with an unstyled app (#54)",
				Code: `	// Build Tailwind so styles.css exists after clone/install (#141).
	if _, err := os.Stat(filepath.Join(dir, cssInput)); err == nil {
		_, _ = fmt.Fprintln(c.Out, "→ tailwind build (styles.css)")
		if err := runTailwindBuild(dir, false); err != nil {
			// Fail loudly: exiting 0 here left users on an unstyled app (#54).
			return fmt.Errorf("css build failed: %w — retry: amarra-cais css (NODE_ENV=production skips npm devDependencies)", err)
		}
	}`,
				SourceURL:   "https://github.com/puppe1990/amarra-cais/blob/93c44c3/internal/cli/commands.go#L48-L55",
				SourceLabel: "internal/cli/commands.go:48",
			},
			{Type: "paragraph", Text: "The testing culture follows: 1.148 Go test functions and 203 JS tests for roughly 27 thousand lines of code, with the test code itself above 29 thousand. The repo AGENTS.md opens by saying who it is written for, and it is not only people: the primary reader is often an agent, with cheap rg, expensive full reads and files truncated at two thousand lines. Greppable names, small modules and headless tests become a maintenance requirement when half the changes go through an agent."},
			{
				Type:        "quote",
				Text:        "Primary reader of this repo is often an LLM agent: cheap `rg`, expensive full reads, file truncation ~2k lines. Optimize for grep-unique names, small modules, and headless tests.",
				AuthorQuote: "Amarra-cais AGENTS.md",
				SourceURL:   "https://github.com/puppe1990/amarra-cais/blob/93c44c3/AGENTS.md#L3",
				SourceLabel: "AGENTS.md:3",
			},
			{Type: "heading2", Text: "The price, with name and number"},
			{Type: "paragraph", Text: "Not everything is a gain. The shell outside #amarra-main is not updated by the morph, and that bites: on this site, switching language left the nav and the toggle in the previous language until I filed an issue on the framework itself (#397). App JS is cache-first in the service worker, so a fix in web/static/js only lands with pwa --bump (#399). data-amarra-skip is a per-link decision: forgetting it becomes a silent bug."},
			{Type: "paragraph", Text: "There is also no rich client state, no choreographed transition between screens and no offline mutation queue of its own. When the app needs that, Amarra is not the tool. It was designed for the case where the screen is a projection of the database and the server decides."},
			{Type: "heading2", Text: "What the numbers say"},
			{Type: "paragraph", Text: "Two months after the fork: 687 commits, 206 issues, v0.27.0, 20 projects in this workspace depending on the module (from an EPUB translator to a rental ERP), 23 kit components and a 117 KB script in the browser. The scale is mine: each app with one SQLite, one server, one writer. Amarra was not designed for another."},
			{Type: "paragraph", Text: "This site is the most direct example. Sintaxe was born a SPA (React + Vite) and is Amarra today: the editorial catalog is Go code, scheduling is a published_at in the future, the English edition is one more column. The migration is in the public repository history, with the old SPA in the commit before."},
			{
				Type:        "paragraph",
				Text:        "The framework is open at",
				SourceURL:   "https://github.com/puppe1990/amarra-cais",
				SourceLabel: "github.com/puppe1990/amarra-cais",
			},
		},
	}
}
