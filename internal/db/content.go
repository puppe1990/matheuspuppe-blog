package db

import "github.com/puppe1990/matheuspuppe-blog/internal/models"

// seedAuthors plants the only byline of the journal. The catalog is signed by
// Matheus Puppe; essays are appended here as they are written, and scheduling
// is just a PublishedAt in the future (see store.publishedNow).
func seedAuthors() []models.Author {
	return []models.Author{
		{
			Slug:      "matheus-puppe",
			Name:      "Matheus Puppe",
			Role:      "Engenheiro de software & fundador da Gestão Bem",
			Bio:       "Construo produtos de software de ponta a ponta — do framework Amarra-cais ao Purple Stock e ao Hora Solar — e escrevo sobre sistemas, IA e o que aprendo operando tudo em produção.",
			RoleEN:    "Software engineer & founder of Gestão Bem",
			BioEN:     "I build software products end to end — from the Amarra-cais framework to Purple Stock and Hora Solar — and write about systems, AI and what I learn running them in production.",
			AvatarURL: "/static/img/author-matheus.jpg",
			GitHub:    "puppe1990",
		},
	}
}

// seedArticles is the published catalog. The first essay is live; the next
// ones enter here as they are written, and scheduling is just a PublishedAt in
// the future.
func seedArticles() []models.Article {
	return []models.Article{primeiroEnsaio()}
}

// seedComment pairs a starter comment with its article slug so RunSeeds can
// resolve the article id without hardcoding AUTOINCREMENT values.
type seedComment struct {
	ArticleSlug string
	Comment     models.Comment
}

func seedComments() []seedComment {
	return nil
}

// authorRef points an article at its columnist; RunSeeds resolves the slug to
// the row id before the upsert.
func authorRef(slug string) models.Author {
	return models.Author{Slug: slug}
}
