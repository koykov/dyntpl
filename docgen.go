package dyntpl

import (
	"bytes"
	"fmt"
	"io"
)

type DocgenFormat string

const (
	DocgenFormatMarkdown DocgenFormat = "markdown"
	DocgenFormatHTML                  = "html"
	DocgenFormatJSON                  = "json"
)

func Docgen(format DocgenFormat) ([]byte, error) {
	var buf bytes.Buffer
	err := WriteDocgen(&buf, format)
	return buf.Bytes(), err
}

func WriteDocgen(w io.Writer, format DocgenFormat) error {
	switch format {
	case DocgenFormatMarkdown:
		return writeDocgenMarkdown(w)
	case DocgenFormatHTML:
		return writeDocgenHTML(w)
	case DocgenFormatJSON:
		return writeDocgenJSON(w)
	}
	return fmt.Errorf("unknown format: %s", format)
}

func writeDocgenMarkdown(w io.Writer) error {
	_, _ = w.Write([]byte("# API\n\n"))

	_, _ = w.Write([]byte("## Modifiers\n\n"))
	for i := 0; i < len(modBuf); i++ {
		tuple := &modBuf[i]
		_, _ = w.Write([]byte("### "))
		_, _ = w.Write([]byte(tuple.id))
		_, _ = w.Write([]byte("\n"))
		if len(tuple.alias) > 0 {
			_, _ = w.Write([]byte("Alias: `" + tuple.alias + "`\n"))
		}
		_, _ = w.Write([]byte("\n"))
		if len(tuple.params) > 0 {
			_, _ = w.Write([]byte("Params:\n"))
			for j := 0; j < len(tuple.params); j++ {
				param := &tuple.params[j]
				_, _ = w.Write([]byte("* `" + param.param + "`"))
				if len(param.desc) > 0 {
					_, _ = w.Write([]byte(" " + param.desc))
				}
				_, _ = w.Write([]byte("\n"))
			}
			_, _ = w.Write([]byte("\n"))
		}

		if len(tuple.desc) > 0 {
			_, _ = w.Write([]byte(tuple.desc))
			_, _ = w.Write([]byte("\n\n"))
		}

		if len(tuple.example) > 0 {
			_, _ = w.Write([]byte("Example:\n```\n"))
			_, _ = w.Write([]byte(tuple.example))
			_, _ = w.Write([]byte("\n```\n\n"))
		}
	}

	return nil
}

func writeDocgenHTML(w io.Writer) error {
	return nil
}

func writeDocgenJSON(w io.Writer) error {
	return nil
}
