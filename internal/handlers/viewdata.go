package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/boot"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/matheuspuppe-blog/internal/editorial"
	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// bookmarkCookie stores the reader's saved library. HttpOnly: the set is only
// ever rendered server-side, and the value is bounded by editorial.maxBookmarks
// so it stays inside the ~4 KB cookie limit.
const bookmarkCookie = "sintaxe_salvos"

const bookmarkMaxAge = 365 * 24 * 60 * 60

// localeOf reads the locale resolved by i18n.LocaleMiddleware; /en paths are
// forced to English by pathLocale in internal/app.
func localeOf(r *http.Request) string {
	if cat := i18n.CatalogFromRequest(r); cat != nil && cat.Locale() == editorial.LocaleEN {
		return editorial.LocaleEN
	}
	return editorial.LocalePT
}

// localePath prefixes a root-relative path with /en when the request is
// English ("/" becomes "/en", "/?q=go" becomes "/en?q=go").
func localePath(locale, path string) string {
	if locale != editorial.LocaleEN {
		return path
	}
	if path == "/" {
		return "/en"
	}
	if strings.HasPrefix(path, "/?") {
		return "/en" + strings.TrimPrefix(path, "/")
	}
	return "/en" + path
}

func articlePath(slug, locale string) string {
	return localePath(locale, "/artigo/"+slug)
}

func authorPath(slug, locale string) string {
	return localePath(locale, "/autor/"+slug)
}

// currentPath is the reader-visible URI. The lang=en that pathLocale injects
// for /en routes is an implementation detail, so it never reaches the page
// (`next` fields, hreflang self link).
func currentPath(r *http.Request) string {
	if !strings.HasPrefix(r.URL.Path, "/en") {
		return r.URL.RequestURI()
	}
	query := r.URL.Query()
	query.Del("lang")
	path := r.URL.Path
	if encoded := query.Encode(); encoded != "" {
		return path + "?" + encoded
	}
	return path
}

// alternatePath maps the current path to its counterpart language.
func alternatePath(path, locale string) string {
	if locale == editorial.LocaleEN {
		trimmed := strings.TrimPrefix(path, "/en")
		if trimmed == "" {
			return "/"
		}
		return trimmed
	}
	if path == "/" {
		return "/en"
	}
	return "/en" + path
}

// FilterLink is one pressable filter chip (category or reading time).
type FilterLink struct {
	Label  string
	Href   string
	Count  int
	Active bool
}

// categoryLinksFor builds the segmented caderno bar with live counts; the
// active search and reading-time filters survive the category switch.
func categoryLinksFor(articles []models.Article, search, activeCategory string, maxMinutes int, locale string) []FilterLink {
	counts := editorial.CategoryCounts(articles)
	links := make([]FilterLink, 0, len(editorial.Categories))
	for _, category := range editorial.Categories {
		links = append(links, FilterLink{
			Label:  editorial.CategoryLabel(category, locale),
			Href:   localePath(locale, editorial.CategoryHref(category, search, maxMinutes)),
			Count:  counts[category],
			Active: category == activeCategory,
		})
	}
	return links
}

// timeLinksFor renders the reading-time chips with the request catalog.
func timeLinksFor(search, activeCategory string, maxMinutes int, locale string, tr func(key string) string) []FilterLink {
	options := []struct {
		label   string
		minutes int
	}{
		{tr("filters.all"), 0},
		{tr("filters.le5"), 5},
		{tr("filters.le8"), 8},
	}
	links := make([]FilterLink, 0, len(options))
	for _, option := range options {
		links = append(links, FilterLink{
			Label:  option.label,
			Href:   localePath(locale, editorial.CategoryHref(activeCategory, search, option.minutes)),
			Active: maxMinutes == option.minutes,
		})
	}
	return links
}

// frameworkVersion is the installed Cais version behind the brand label. A dev
// build reports a Go pseudo-version, so the label keeps the semver core and the
// tooltip carries the full string.
func frameworkVersion() (label, full string) {
	full = boot.CaisVersion()
	label = full
	if i := strings.Index(label, "-0."); i > 0 {
		label = label[:i]
	}
	return label, full
}

func amarraData(r *http.Request, site meta.Site, extra map[string]any) map[string]any {
	s := meta.ForRequest(site, r)
	label, full := frameworkVersion()
	locale := localeOf(r)
	data := map[string]any{
		"Site":                 s,
		"CSRFToken":            s.CSRFToken,
		"Flash":                s.Flash,
		"FrameworkVersion":     label,
		"FrameworkVersionFull": full,
		"Lang":                 locale,
		"AltURL":               alternatePath(r.URL.Path, locale),
		// CurrentPath feeds the `next` field of bookmark/newsletter forms so a
		// Drive redirect returns the reader to the exact filtered list view.
		"CurrentPath": currentPath(r),
		"SavedCount":  len(bookmarksFromRequest(r)),
	}
	if cat := i18n.CatalogFromRequest(r); cat != nil {
		data["Locale"] = cat.Locale()
	}
	for k, v := range extra {
		data[k] = v
	}
	return data
}

// writeView names the layout at the call site so a future second layout (a
// standalone landing) does not render inside the journal chrome (#66).
func writeView(w http.ResponseWriter, r *http.Request, views *view.Renderer, cfg cais.Config, layout, name string, data any, status int) {
	view.Write(w, r, views, view.Page{Layout: layout, Name: name, Data: data, Status: status}, cfg)
}

// ArticleView decorates a row with display-ready fields, so templates never
// call date or number formatters. Slug and URL follow the request locale; the
// canonical Portuguese slug is kept as the bookmark identity.
type ArticleView struct {
	models.Article
	DateLabel     string
	ViewsLabel    string
	ShareText     string
	Saved         bool
	URL           string
	AuthorURL     string
	ClapURL       string
	CommentURL    string
	CanonicalSlug string
	CategoryLabel string
}

func newArticleView(article models.Article, saved []string, locale string) ArticleView {
	canonical := article.Slug
	localized := editorial.Localize(article, locale)
	url := articlePath(localized.Slug, locale)
	return ArticleView{
		Article:       localized,
		DateLabel:     editorial.FormatDateIn(localized.PublishedAt, locale),
		ViewsLabel:    editorial.FormatThousands(localized.Views),
		ShareText:     editorial.ShareText(localized.Title, localized.Subtitle),
		Saved:         containsSlug(saved, canonical),
		URL:           url,
		AuthorURL:     authorPath(localized.Author.Slug, locale),
		ClapURL:       url + "/aplausos",
		CommentURL:    url + "/comentarios",
		CanonicalSlug: canonical,
		CategoryLabel: editorial.CategoryLabel(localized.Category, locale),
	}
}

func newArticleViews(articles []models.Article, saved []string, locale string) []ArticleView {
	views := make([]ArticleView, len(articles))
	for i, article := range articles {
		views[i] = newArticleView(article, saved, locale)
	}
	return views
}

// LeadArticle keeps a typed pointer for the hero slot; a nil pointer renders
// no hero while the feed is empty.
func leadArticleView(article *models.Article, saved []string, locale string) *ArticleView {
	if article == nil {
		return nil
	}
	view := newArticleView(*article, saved, locale)
	return &view
}

// CommentView decorates a comment with its formatted date.
type CommentView struct {
	models.Comment
	DateLabel string
}

func newCommentViews(comments []models.Comment, locale string) []CommentView {
	views := make([]CommentView, len(comments))
	for i, comment := range comments {
		views[i] = CommentView{Comment: comment, DateLabel: editorial.FormatDateIn(comment.CreatedAt, locale)}
	}
	return views
}

func containsSlug(slugs []string, slug string) bool {
	for _, existing := range slugs {
		if existing == slug {
			return true
		}
	}
	return false
}

// bookmarksFromRequest returns the saved slugs, or nil when nothing is saved.
func bookmarksFromRequest(r *http.Request) []string {
	cookie, err := r.Cookie(bookmarkCookie)
	if err != nil {
		return nil
	}
	return editorial.ParseBookmarkCookie(cookie.Value)
}

func writeBookmarkCookie(w http.ResponseWriter, slugs []string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     bookmarkCookie,
		Value:    editorial.EncodeBookmarks(slugs),
		Path:     "/",
		MaxAge:   bookmarkMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
}

// sanitizeNext only accepts same-site absolute paths so a crafted `next`
// field cannot redirect readers to another host.
func sanitizeNext(next string) string {
	if next == "" || !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		return "/"
	}
	return next
}

// intParam reads an integer query parameter with a fallback for absent or
// malformed values; negative values fall back too.
func intParam(r *http.Request, name string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil || value < 0 {
		return fallback
	}
	return value
}
