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
	data := map[string]any{
		"Site":                 s,
		"CSRFToken":            s.CSRFToken,
		"Flash":                s.Flash,
		"FrameworkVersion":     label,
		"FrameworkVersionFull": full,
		// CurrentPath feeds the `next` field of bookmark/newsletter forms so a
		// Drive redirect returns the reader to the exact filtered list view.
		"CurrentPath": r.URL.RequestURI(),
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
// call date or number formatters.
type ArticleView struct {
	models.Article
	DateLabel  string
	ViewsLabel string
	ShareText  string
	Saved      bool
}

func newArticleView(article models.Article, saved []string) ArticleView {
	return ArticleView{
		Article:    article,
		DateLabel:  editorial.FormatDate(article.PublishedAt),
		ViewsLabel: editorial.FormatThousands(article.Views),
		ShareText:  editorial.ShareText(article.Title, article.Subtitle),
		Saved:      containsSlug(saved, article.Slug),
	}
}

func newArticleViews(articles []models.Article, saved []string) []ArticleView {
	views := make([]ArticleView, len(articles))
	for i, article := range articles {
		views[i] = newArticleView(article, saved)
	}
	return views
}

// LeadArticle keeps a typed pointer for the hero slot; a nil pointer renders
// no hero while the feed is empty.
func leadArticleView(article *models.Article, saved []string) *ArticleView {
	if article == nil {
		return nil
	}
	view := newArticleView(*article, saved)
	return &view
}

// CommentView decorates a comment with its formatted date.
type CommentView struct {
	models.Comment
	DateLabel string
}

func newCommentViews(comments []models.Comment) []CommentView {
	views := make([]CommentView, len(comments))
	for i, comment := range comments {
		views[i] = CommentView{Comment: comment, DateLabel: editorial.FormatDate(comment.CreatedAt)}
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
