package app

import (
	"github.com/puppe1990/amarra-cais/pkg/amarra/live"
	"github.com/puppe1990/amarra-cais/pkg/cais"

	"github.com/puppe1990/matheuspuppe-blog/internal/handlers"
)

func registerRoutes(r *cais.Router, deps Deps, cfg cais.Config) {
	home := handlers.NewHomeHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)
	artigo := handlers.NewArticleHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)
	salvos := handlers.NewSavedHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)
	sobre := handlers.NewSobreHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)
	autor := handlers.NewAuthorHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)
	newsletter := handlers.NewNewsletterHandler(deps.Store, deps.Catalog, cfg)
	sitemap := handlers.NewSitemapHandler(deps.Store, deps.Site, cfg)
	feed := handlers.NewFeedHandler(deps.Store, deps.Site, cfg)
	og := handlers.NewOGHandler(deps.Store, deps.Site, cfg)
	notFound := handlers.NewNotFoundHandler(deps.Views, deps.Site, cfg)

	r.Get("/", home.ServeHTTP)
	r.Get("/en", home.ServeHTTP)
	r.Get("/artigo/{slug}", cais.StringParam("slug", artigo.Show))
	r.Post("/artigo/{slug}/comentarios", cais.StringParam("slug", artigo.CommentPost))
	r.Post("/artigo/{slug}/aplausos", cais.StringParam("slug", artigo.ClapPost))
	r.Get("/en/artigo/{slug}", cais.StringParam("slug", artigo.Show))
	r.Post("/en/artigo/{slug}/comentarios", cais.StringParam("slug", artigo.CommentPost))
	r.Post("/en/artigo/{slug}/aplausos", cais.StringParam("slug", artigo.ClapPost))
	r.Get("/salvos", salvos.ServeHTTP)
	r.Get("/en/salvos", salvos.ServeHTTP)
	r.Post("/salvos/{slug}", cais.StringParam("slug", salvos.Toggle))
	r.Post("/en/salvos/{slug}", cais.StringParam("slug", salvos.Toggle))
	r.Get("/sobre", sobre.ServeHTTP)
	r.Get("/en/sobre", sobre.ServeHTTP)
	r.Get("/autor/{slug}", cais.StringParam("slug", autor.Serve))
	r.Get("/en/autor/{slug}", cais.StringParam("slug", autor.Serve))
	r.Post("/newsletter", newsletter.Post)
	r.Post("/en/newsletter", newsletter.Post)
	r.Get("/sitemap.xml", sitemap.ServeHTTP)
	r.Get("/rss.xml", feed.ServeHTTP)
	r.Get("/robots.txt", handlers.RobotsHandler(deps.Site))
	r.Get("/og/{slug}", cais.StringParam("slug", og.Serve))
	r.NotFound(notFound.ServeHTTP)
}

func registerLiveViews(hub *live.Hub, deps Deps) {
	_ = hub
	_ = deps
	// cais:live-views
}
