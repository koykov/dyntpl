package dyntpl

import (
	"os"
	"testing"
)

func TestDocgen(t *testing.T) {
	generate := os.Getenv("DOCGEN_GENERATE") == "1"

	t.Run("markdown", func(t *testing.T) {
		docs, _ := Docgen(DocgenFormatMarkdown)
		if generate {
			_ = os.WriteFile("docgen.md", docs, 0644)
		}
	})
	t.Run("html", func(t *testing.T) {
		docs, _ := Docgen(DocgenFormatHTML)
		if generate {
			_ = os.WriteFile("docgen.html", docs, 0644)
		}
	})
	t.Run("json", func(t *testing.T) {
		docs, _ := Docgen(DocgenFormatJSON)
		if generate {
			_ = os.WriteFile("docgen.json", docs, 0644)
		}
	})
}
