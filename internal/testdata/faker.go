// Package testdata builds fake editorial rows for tests and fixtures.
//
// Every value derives from one seeded gofakeit Faker, so a failing test
// reproduces with identical data on every run.
package testdata

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// Seed keeps generated fixtures identical across runs.
const Seed = uint64(20261103)

func New() *gofakeit.Faker {
	return gofakeit.New(Seed)
}

// Author returns a fresh fake columnist with a unique slug.
func Author(f *gofakeit.Faker, n int) models.Author {
	return models.Author{
		Slug:      fmt.Sprintf("autor-%d", n),
		Name:      f.Name(),
		Role:      f.JobTitle(),
		Bio:       f.Sentence(12),
		AvatarURL: f.URL(),
		Twitter:   "@" + f.Username(),
		GitHub:    f.Username(),
	}
}

// Article returns a fake essay with n unique slides; publishedAt spreads the
// feed backwards from base so ordering assertions stay readable.
func Article(f *gofakeit.Faker, n int, author models.Author, base time.Time) models.Article {
	return models.Article{
		Slug:            fmt.Sprintf("%s-ensaio-%d", author.Slug, n),
		Title:           f.Sentence(6),
		Subtitle:        f.Sentence(8),
		LeadParagraph:   f.Paragraph(2, 4, 8, " "),
		Category:        f.RandomString([]string{"IA & Modelos", "Engenharia de Sistemas", "Hardware & Chips"}),
		ReadTimeMinutes: f.Number(3, 12),
		PublishedAt:     base.AddDate(0, 0, -n),
		CoverImage:      "/static/img/test-cover.jpg",
		Author:          author,
		Tags:            []string{f.Word(), f.Word()},
		Blocks: []models.ContentBlock{
			{Type: "paragraph", Text: f.Sentence(10)},
			{Type: "takeaways", Items: []string{f.Sentence(4), f.Sentence(4)}},
		},
		Claps: f.Number(0, 500),
		Views: f.Number(0, 5000),
		Tier:  f.Number(1, 3),
	}
}

// Articles returns n fake essays for one author.
func Articles(f *gofakeit.Faker, n int, author models.Author, base time.Time) []models.Article {
	out := make([]models.Article, n)
	for i := range out {
		out[i] = Article(f, i+1, author, base)
	}
	return out
}

// Comment returns a fake reader reflection for an article.
func Comment(f *gofakeit.Faker, articleID int64) models.Comment {
	return models.Comment{
		ArticleID:  articleID,
		AuthorName: f.Name(),
		AuthorRole: f.JobTitle(),
		Content:    f.Sentence(14),
		CreatedAt:  time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC),
		Claps:      f.Number(0, 20),
	}
}
