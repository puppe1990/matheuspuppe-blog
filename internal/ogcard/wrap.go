package ogcard

import (
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func wrapLines(face font.Face, text string, maxPx, maxLines int) []string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return nil
	}
	max := fixed.I(maxPx)
	var lines []string
	current := words[0]
	for _, word := range words[1:] {
		trial := current + " " + word
		if font.MeasureString(face, trial) <= max {
			current = trial
			continue
		}
		lines = append(lines, current)
		current = word
	}
	lines = append(lines, current)
	if maxLines > 0 && len(lines) > maxLines {
		lines = lines[:maxLines]
		lines[maxLines-1] = ellipsize(face, lines[maxLines-1], maxPx)
	}
	return lines
}

func ellipsize(face font.Face, s string, maxPx int) string {
	max := fixed.I(maxPx)
	if font.MeasureString(face, s+"…") <= max {
		if font.MeasureString(face, s) <= max {
			return s
		}
		return s + "…"
	}
	runes := []rune(s)
	for len(runes) > 0 {
		for len(runes) > 0 && runes[len(runes)-1] == ' ' {
			runes = runes[:len(runes)-1]
		}
		cand := string(runes) + "…"
		if font.MeasureString(face, cand) <= max {
			return cand
		}
		runes = runes[:len(runes)-1]
	}
	return "…"
}

func measure(face font.Face, s string) int {
	return font.MeasureString(face, s).Round()
}
