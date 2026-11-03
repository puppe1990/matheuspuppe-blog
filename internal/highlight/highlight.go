// Package highlight renders code blocks for the reader. Chroma runs on the
// server with classes (no inline styles, no client script), so the palette
// lives in the stylesheet and the CSP stays untouched.
package highlight

import (
	"bytes"
	"html/template"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
)

// Code tokenises code in language and returns spans ready to sit inside the
// reader's <pre class="amarra-code">. Unknown languages and lexer failures
// degrade to escaped plain text, never to raw markup.
func Code(code, language string) template.HTML {
	if code == "" {
		return ""
	}
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return template.HTML(template.HTMLEscapeString(code))
	}
	formatter := chromahtml.New(chromahtml.WithClasses(true), chromahtml.PreventSurroundingPre(true))
	var buffer bytes.Buffer
	if err := formatter.Format(&buffer, chromaFallbackStyle, iterator); err != nil {
		return template.HTML(template.HTMLEscapeString(code))
	}
	return template.HTML(buffer.String())
}

// chromaFallbackStyle only feeds the formatter's colouring decisions; with
// WithClasses(true) no colour leaves this package. The stylesheet owns the
// palette (see input.css).
var chromaFallbackStyle = chroma.MustNewStyle("sintaxe", chroma.StyleEntries{
	chroma.Keyword:       "bold #f59e0b",
	chroma.Comment:       "italic #6f6a60",
	chroma.LiteralString: "#93a87c",
	chroma.LiteralNumber: "#cd9e5c",
	chroma.NameFunction:  "#efe9de",
	chroma.NameBuiltin:   "#efe9de",
	chroma.Operator:      "#b9b2a7",
	chroma.Punctuation:   "#8f8a80",
	chroma.Background:    "#efe9de",
	chroma.Error:         "#e5484d",
})
