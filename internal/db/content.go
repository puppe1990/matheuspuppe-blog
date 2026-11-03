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
			AvatarURL: "/static/img/author-matheus.jpg",
			GitHub:    "puppe1990",
		},
	}
}

// seedArticles is the published catalog: empty on purpose while the first
// edition is being written. Scheduled essays enter here with a future
// PublishedAt and appear on their day without another deploy.
func seedArticles() []models.Article {
	return nil
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
