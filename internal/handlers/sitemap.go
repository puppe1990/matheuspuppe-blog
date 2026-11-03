package handlers

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

// SitemapHandler serves /sitemap.xml. Crawler GETs must not set cookies —
// Search Console rejects a sitemap that arrives with Set-Cookie (#342).
type SitemapHandler struct {
	store store.Store
	site  meta.Site
	cfg   cais.Config
}

func NewSitemapHandler(st store.Store, site meta.Site, cfg cais.Config) *SitemapHandler {
	return &SitemapHandler{store: st, site: site, cfg: cfg}
}

type urlEntry struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type urlSet struct {
	XMLName xml.Name   `xml:"urlset"`
	Xmlns   string     `xml:"xmlns,attr"`
	URLs    []urlEntry `xml:"url"`
}

func (h *SitemapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	articles, err := h.store.ListArticles(r.Context())
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	base := strings.TrimSuffix(h.site.AppURL, "/")

	set := urlSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for _, path := range []string{"/", "/sobre", "/salvos", "/en", "/en/sobre", "/en/salvos"} {
		set.URLs = append(set.URLs, urlEntry{Loc: base + path})
	}
	for _, article := range articles {
		set.URLs = append(set.URLs, urlEntry{
			Loc:     base + "/artigo/" + article.Slug,
			LastMod: article.PublishedAt.Format("2006-01-02"),
		})
		if article.SlugEN != "" {
			set.URLs = append(set.URLs, urlEntry{
				Loc:     base + "/en/artigo/" + article.SlugEN,
				LastMod: article.PublishedAt.Format("2006-01-02"),
			})
		}
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write([]byte(xml.Header))
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	_ = encoder.Encode(set)
}

// RobotsHandler serves /robots.txt with the sitemap pointer; no cookies here
// either, same crawler rule as the sitemap (#342).
func RobotsHandler(site meta.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		base := strings.TrimSuffix(site.AppURL, "/")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("User-agent: *\nAllow: /\nSitemap: " + base + "/sitemap.xml\n"))
	}
}

// FeedHandler serves the RSS 2.0 feed with the ten newest essays.
type FeedHandler struct {
	store store.Store
	site  meta.Site
	cfg   cais.Config
}

func NewFeedHandler(st store.Store, site meta.Site, cfg cais.Config) *FeedHandler {
	return &FeedHandler{store: st, site: site, cfg: cfg}
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Items       []rssItem `xml:"item"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

func (h *FeedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	articles, err := h.store.ListArticles(r.Context())
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	base := strings.TrimSuffix(h.site.AppURL, "/")

	channel := rssChannel{
		Title:       "Sintaxe — Jornal de Tecnologia",
		Link:        base,
		Description: "Publicação editorial de tecnologia e engenharia de software.",
		Language:    "pt-BR",
	}
	for i, article := range articles {
		if i == 10 {
			break
		}
		link := base + "/artigo/" + article.Slug
		channel.Items = append(channel.Items, rssItem{
			Title:       article.Title,
			Link:        link,
			GUID:        link,
			PubDate:     article.PublishedAt.Format(time.RFC1123Z),
			Description: article.Subtitle,
		})
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	_, _ = w.Write([]byte(xml.Header))
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	_ = encoder.Encode(rssFeed{Version: "2.0", Channel: channel})
}
