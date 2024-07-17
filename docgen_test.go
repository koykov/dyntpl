package dyntpl

import "testing"

const docgenVerbose = true

func TestDocgen(t *testing.T) {
	t.Run("markdown", func(t *testing.T) {
		docs, _ := Docgen(DocgenFormatMarkdown)
		if docgenVerbose {
			println(string(docs))
		}
	})
	t.Run("html", func(t *testing.T) {
		docs, _ := Docgen(DocgenFormatHTML)
		if docgenVerbose {
			println(string(docs))
		}
	})
	t.Run("json", func(t *testing.T) {
		docs, _ := Docgen(DocgenFormatJSON)
		if docgenVerbose {
			println(string(docs))
		}
	})
}
