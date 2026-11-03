package editorial

import (
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// Locale tags used by the journal: the root is pt-BR, /en is the English
// edition. Content and URLs follow the locale; chrome follows the catalog.
const (
	LocalePT = "pt"
	LocaleEN = "en"
)

// Localize returns the article in the requested locale, falling back to the
// Portuguese edition when the English one is missing. The English slug travels
// with the translation, so the URL changes with the language.
func Localize(article models.Article, locale string) models.Article {
	if locale != LocaleEN || !HasEnglish(article) {
		return article
	}
	article.Slug = article.SlugEN
	article.Title = article.TitleEN
	article.Subtitle = article.SubtitleEN
	article.LeadParagraph = article.LeadParagraphEN
	article.Tags = article.TagsEN
	article.Blocks = article.BlocksEN
	article.Author = LocalizeAuthor(article.Author, locale)
	return article
}

// LocalizeAuthor swaps role and bio for the English edition when present.
func LocalizeAuthor(author models.Author, locale string) models.Author {
	if locale != LocaleEN {
		return author
	}
	if author.RoleEN != "" {
		author.Role = author.RoleEN
	}
	if author.BioEN != "" {
		author.Bio = author.BioEN
	}
	return author
}

// HasEnglish reports whether the article carries a complete English edition.
func HasEnglish(article models.Article) bool {
	return article.SlugEN != "" && article.TitleEN != "" && len(article.BlocksEN) > 0
}

var categoryEN = map[string]string{
	"Todos":                    "All",
	"IA & Modelos":             "AI & Models",
	"Engenharia de Sistemas":   "Systems Engineering",
	"Hardware & Chips":         "Hardware & Chips",
	"Design & UX":              "Design & UX",
	"Segurança & Criptografia": "Security & Cryptography",
	"Opinião & Futuro":         "Opinion & Future",
}

// CategoryLabel translates a caderno name for display. The stored value stays
// in Portuguese (it is the filter key); only the label changes.
func CategoryLabel(category, locale string) string {
	if locale != LocaleEN {
		return category
	}
	if label, ok := categoryEN[category]; ok {
		return label
	}
	return category
}

// FormatDateIn renders the editorial date in the request locale.
func FormatDateIn(t time.Time, locale string) string {
	if locale == LocaleEN {
		return t.Format("January 2, 2006")
	}
	return FormatDate(t)
}
