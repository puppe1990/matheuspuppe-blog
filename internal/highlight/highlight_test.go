package highlight

import (
	"strings"
	"testing"
)

func TestCode_GoGetsClassesWithoutPreWrapper(t *testing.T) {
	got := string(Code("func main() {}", "go"))
	if !strings.Contains(got, `class="kd">func`) {
		t.Fatalf("keyword span missing: %q", got)
	}
	if !strings.Contains(got, `class="nf">main`) {
		t.Fatalf("function name span missing: %q", got)
	}
	if strings.Contains(got, "<pre") {
		t.Fatalf("formatter added its own <pre>: %q", got)
	}
}

func TestCode_UnknownLanguageFallsBackEscaped(t *testing.T) {
	got := string(Code("<script>alert(1)</script>", "nao-existe"))
	if strings.Contains(got, "<script>") {
		t.Fatalf("raw markup leaked into the page: %q", got)
	}
	if !strings.Contains(got, "alert(1)") {
		t.Fatalf("code text lost in the fallback: %q", got)
	}
}

func TestCode_PlainTextLanguageStaysReadable(t *testing.T) {
	got := string(Code("docker ps -a", "shell"))
	if !strings.Contains(got, "docker") || !strings.Contains(got, "ps") {
		t.Fatalf("shell snippet lost text: %q", got)
	}
}

func TestCode_EmptyCodeIsEmpty(t *testing.T) {
	if got := string(Code("", "go")); got != "" {
		t.Fatalf("empty code rendered as %q", got)
	}
}
