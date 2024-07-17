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
}
