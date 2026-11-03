package editorial

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// WordsPerMinute is the reading-pace constant behind ReadTimeMinutes.
const WordsPerMinute = 200

// maxBookmarks bounds the saved-articles cookie so a swollen value never
// exceeds the ~4 KB cookie limit browsers enforce.
const maxBookmarks = 100

// ReadingTime converts a word count into minutes, never less than one.
func ReadingTime(words int) int {
	if words <= 0 {
		return 1
	}
	minutes := (words + WordsPerMinute - 1) / WordsPerMinute
	if minutes < 1 {
		return 1
	}
	return minutes
}

// WordCount totals the words across every text-bearing block.
func WordCount(blocks []models.ContentBlock) int {
	count := 0
	for _, block := range blocks {
		count += len(strings.Fields(block.Text))
		for _, item := range block.Items {
			count += len(strings.Fields(item))
		}
	}
	return count
}

// ToggleBookmark adds or removes a slug, preserving insertion order.
func ToggleBookmark(slugs []string, slug string) []string {
	for i, existing := range slugs {
		if existing == slug {
			return append(slugs[:i:i], slugs[i+1:]...)
		}
	}
	if len(slugs) >= maxBookmarks {
		return slugs
	}
	return append(slugs, slug)
}

// AddClaps raises a clap counter with a minimum increment of one, mirroring the
// reader's batch-clap behavior.
func AddClaps(current, increment int) int {
	if increment < 1 {
		increment = 1
	}
	return current + increment
}

// RemainingReadingTime estimates the minutes left at a scroll percentage.
func RemainingReadingTime(readTimeMinutes, scrollPercent int) int {
	if scrollPercent < 0 {
		scrollPercent = 0
	}
	if scrollPercent > 100 {
		scrollPercent = 100
	}
	remaining := readTimeMinutes * (100 - scrollPercent)
	minutes := (remaining + 99) / 100
	if minutes < 0 {
		return 0
	}
	return minutes
}

// ShareText builds the Web Share API fallback text for one article.
func ShareText(title, subtitle string) string {
	if subtitle == "" {
		return title
	}
	return title + " — " + subtitle
}

var monthsPtBR = [...]string{
	"janeiro", "fevereiro", "março", "abril", "maio", "junho",
	"julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
}

// FormatDate renders an editorial date in pt-BR ("28 de março de 2026").
func FormatDate(t time.Time) string {
	return strconv.Itoa(t.Day()) + " de " + monthsPtBR[int(t.Month())-1] + " de " + strconv.Itoa(t.Year())
}

// FormatThousands groups an integer with pt-BR dots ("4.890").
func FormatThousands(n int) string {
	negative := n < 0
	if negative {
		n = -n
	}
	digits := strconv.Itoa(n)
	var out strings.Builder
	for i, digit := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			out.WriteByte('.')
		}
		out.WriteRune(digit)
	}
	if negative {
		return "-" + out.String()
	}
	return out.String()
}

// ParseBookmarkCookie reads the saved-articles cookie, trimming blanks and
// duplicates and skipping anything beyond the size cap.
func ParseBookmarkCookie(raw string) []string {
	out := make([]string, 0, 8)
	seen := map[string]bool{}
	for _, part := range strings.Split(raw, ",") {
		slug := strings.TrimSpace(part)
		if slug == "" || seen[slug] {
			continue
		}
		seen[slug] = true
		out = append(out, slug)
		if len(out) >= maxBookmarks {
			break
		}
	}
	return out
}

// EncodeBookmarks serializes the saved slugs for the cookie value.
func EncodeBookmarks(slugs []string) string {
	return strings.Join(slugs, ",")
}

// CategoryHref builds a front-page link that keeps the active search and
// reading-time filters while switching category.
func CategoryHref(category, search string, maxMinutes int) string {
	values := url.Values{}
	if category != "" && category != "Todos" {
		values.Set("categoria", category)
	}
	if strings.TrimSpace(search) != "" {
		values.Set("q", search)
	}
	if maxMinutes > 0 {
		values.Set("tempo", strconv.Itoa(maxMinutes))
	}
	if len(values) == 0 {
		return "/"
	}
	return "/?" + values.Encode()
}
