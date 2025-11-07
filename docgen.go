package dyntpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"
)

type DocgenFormat string

const (
	DocgenFormatMarkdown DocgenFormat = "markdown"
	DocgenFormatHTML     DocgenFormat = "html"
	DocgenFormatJSON     DocgenFormat = "json"
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
		_ = tuple.write(w, DocgenFormatMarkdown, false)
	}

	_, _ = w.Write([]byte("## Condition helpers\n\n"))
	for i := 0; i < len(condBuf); i++ {
		tuple := &condBuf[i]
		_ = tuple.write(w, DocgenFormatMarkdown, false)
	}

	_, _ = w.Write([]byte("## Condition-OK helpers\n\n"))
	for i := 0; i < len(condBuf); i++ {
		tuple := &condBuf[i]
		_ = tuple.write(w, DocgenFormatMarkdown, false)
	}

	_, _ = w.Write([]byte("## Empty checks\n\n"))
	for i := 0; i < len(emptyCheckBuf); i++ {
		tuple := &emptyCheckBuf[i]
		_ = tuple.write(w, DocgenFormatMarkdown, true)
	}
	_, _ = w.Write([]byte("\n"))

	_, _ = w.Write([]byte("## Global variables\n\n"))
	for i := 0; i < len(globBuf); i++ {
		tuple := &globBuf[i]
		_ = tuple.write(w, DocgenFormatMarkdown, true)
	}
	_, _ = w.Write([]byte("\n"))

	_, _ = w.Write([]byte("## Variable-inspector pairs\n\n"))
	for i := 0; i < len(varInsBuf); i++ {
		tuple := &varInsBuf[i]
		_ = tuple.write(w, DocgenFormatMarkdown, true)
	}
	_, _ = w.Write([]byte("\n"))

	return nil
}

func writeDocgenHTML(w io.Writer) error {
	_, _ = w.Write([]byte("<html><head><meta charset=\"utf-8\">"))
	_, _ = w.Write([]byte(`<style>`))
	_, _ = w.Write([]byte(`*{font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Noto Sans,Helvetica,Arial,sans-serif,Apple Color Emoji,Segoe UI Emoji;font-size:16px;line-height:1.5;word-wrap:break-word}body{padding:1em}blockquote,code{margin:0}blockquote{background-color:#e7f3fe;border-left:6px solid #2196F3;padding:10px 15px;}code,pre{border-radius:6px;font-family:monospace;font-size:12px}li,pre{margin-top:.5rem}h1{font-size:var(--h1-size,32px)}h1,h2,h3,h4{font-weight:var(--base-text-weight-semibold,600)}h2{font-size:var(--h2-size,24px)}h3{font-size:var(--h3-size,20px)}p{margin-top:0;margin-bottom:10px}code{padding:.2em .4em;white-space:break-spaces;background-color:rgba(99,110,123,.2)}pre{padding:10px 9px;background-color:rgb(45,51,59,.2)}`))
	_, _ = w.Write([]byte(`</style>`))
	_, _ = w.Write([]byte("</head><body><h1>API</h1>"))

	_, _ = w.Write([]byte(`<a name="mod"></a><h2>Modifiers</h2>`))
	for i := 0; i < len(modBuf); i++ {
		tuple := &modBuf[i]
		_ = tuple.write(w, DocgenFormatHTML, false)
	}

	_, _ = w.Write([]byte(`<a name="cond"></a><h2>Condition helpers</h2>`))
	for i := 0; i < len(condBuf); i++ {
		tuple := &condBuf[i]
		_ = tuple.write(w, DocgenFormatHTML, false)
	}

	_, _ = w.Write([]byte(`<a name="condOK"></a><h2>Condition-OK helpers</h2>`))
	for i := 0; i < len(condOkBuf); i++ {
		tuple := &condOkBuf[i]
		_ = tuple.write(w, DocgenFormatHTML, false)
	}

	_, _ = w.Write([]byte(`<a name="empty"></a><h2>Empty checks</h2>`))
	for i := 0; i < len(emptyCheckBuf); i++ {
		tuple := &emptyCheckBuf[i]
		_ = tuple.write(w, DocgenFormatHTML, true)
	}
	_, _ = w.Write([]byte("\n"))

	_, _ = w.Write([]byte(`<a name="global"></a><h2>Global variables</h2>`))
	for i := 0; i < len(globBuf); i++ {
		tuple := &globBuf[i]
		_ = tuple.write(w, DocgenFormatHTML, true)
	}
	_, _ = w.Write([]byte("\n"))

	_, _ = w.Write([]byte(`<a name="varins"></a><h2>Variable-inspector pairs</h2>`))
	for i := 0; i < len(varInsBuf); i++ {
		tuple := &varInsBuf[i]
		_ = tuple.write(w, DocgenFormatHTML, true)
	}
	_, _ = w.Write([]byte("\n"))

	_, _ = w.Write([]byte("</body></html>"))
	return nil
}

func writeDocgenJSON(w io.Writer) error {
	cpy := func(t *docgen) (r docgenJSON) {
		r.Name = t.name
		r.Alias = t.alias
		r.Type = t.typ
		r.Desc = t.desc
		r.Note = t.note
		r.Example = t.example
		r.Inspector = t.ins
		for i := 0; i < len(t.params); i++ {
			p := &t.params[i]
			r.Params = append(r.Params, docgenParamJSON{
				Name: p.name,
				Desc: p.desc,
			})
		}
		return
	}

	var x docgenContainerJSON
	for i := 0; i < len(modBuf); i++ {
		x.Modifiers = append(x.Modifiers, cpy(&modBuf[i].docgen))
	}
	for i := 0; i < len(condBuf); i++ {
		x.ConditionHelpers = append(x.ConditionHelpers, cpy(&condBuf[i].docgen))
	}
	for i := 0; i < len(condOkBuf); i++ {
		x.ConditionOKHelpers = append(x.ConditionOKHelpers, cpy(&condOkBuf[i].docgen))
	}
	for i := 0; i < len(emptyCheckBuf); i++ {
		x.EmptyCheckHelpers = append(x.EmptyCheckHelpers, cpy(&emptyCheckBuf[i].docgen))
	}
	for i := 0; i < len(globBuf); i++ {
		x.GlobalVariables = append(x.GlobalVariables, cpy(&globBuf[i].docgen))
	}
	for i := 0; i < len(varInsBuf); i++ {
		x.VariableInspectorPairs = append(x.VariableInspectorPairs, cpy(&varInsBuf[i].docgen))
	}
	b, err := json.Marshal(&x)
	_, _ = w.Write(b)

	return err
}

type docgenParam struct {
	name, desc string
}

type docgen struct {
	name, alias, typ, desc, note, example, ins string

	params []docgenParam
}

func (t *docgen) WithDescription(desc string) *docgen {
	t.desc = desc
	return t
}

func (t *docgen) WithType(typ string) *docgen {
	t.typ = typ
	return t
}

func (t *docgen) WithParam(param, desc string) *docgen {
	t.params = append(t.params, docgenParam{
		name: param,
		desc: desc,
	})
	return t
}

func (t *docgen) WithNote(note string) *docgen {
	t.note = note
	return t
}

func (t *docgen) WithExample(example string) *docgen {
	t.example = example
	return t
}

func (t *docgen) write(w io.Writer, format DocgenFormat, compact bool) error {
	htmlEscape := func(s string, nl2br bool) string {
		s = html.EscapeString(s)
		s = reDGHTMLCode.ReplaceAllString(s, "<code>$0</code>")
		if nl2br {
			s = strings.ReplaceAll(s, "\n", "<br/>")
		}
		s = strings.ReplaceAll(s, "`", "")
		return s
	}

	switch {
	case format == DocgenFormatMarkdown && compact:
		_, _ = w.Write([]byte("* `"))
		_, _ = w.Write([]byte(t.name))
		if len(t.typ) > 0 {
			_, _ = w.Write([]byte(" "))
			_, _ = w.Write([]byte(t.typ))
		}
		if len(t.ins) > 0 {
			_, _ = w.Write([]byte(" "))
			_, _ = w.Write([]byte(t.ins))
		}
		_, _ = w.Write([]byte("`"))
		if len(t.desc) > 0 {
			_, _ = w.Write([]byte(" "))
			_, _ = w.Write([]byte(t.desc))
		}
		_, _ = w.Write([]byte("\n"))
	case format == DocgenFormatMarkdown && !compact:
		_, _ = w.Write([]byte("### "))
		_, _ = w.Write([]byte(t.name))
		_, _ = w.Write([]byte("\n"))
		if len(t.alias) > 0 {
			_, _ = w.Write([]byte("Alias: `" + t.alias + "`\n"))
		}
		_, _ = w.Write([]byte("\n"))
		if len(t.params) > 0 {
			_, _ = w.Write([]byte("Params:\n"))
			for j := 0; j < len(t.params); j++ {
				param := &t.params[j]
				_, _ = w.Write([]byte("* `" + param.name + "`"))
				if len(param.desc) > 0 {
					_, _ = w.Write([]byte(" " + param.desc))
				}
				_, _ = w.Write([]byte("\n"))
			}
			_, _ = w.Write([]byte("\n"))
		}

		if len(t.desc) > 0 {
			_, _ = w.Write([]byte(t.desc))
			_, _ = w.Write([]byte("\n\n"))
		}

		if len(t.note) > 0 {
			_, _ = w.Write([]byte("> **_NOTE:_** "))
			_, _ = w.Write([]byte(t.note))
			_, _ = w.Write([]byte("\n\n"))
		}

		if len(t.example) > 0 {
			_, _ = w.Write([]byte("Example:\n```\n"))
			_, _ = w.Write([]byte(t.example))
			_, _ = w.Write([]byte("\n```\n\n"))
		}
	case format == DocgenFormatHTML && compact:
		_, _ = w.Write([]byte("<li><code>"))
		_, _ = w.Write([]byte(t.name))
		if len(t.typ) > 0 {
			_, _ = w.Write([]byte(" "))
			_, _ = w.Write([]byte(t.typ))
		}
		if len(t.ins) > 0 {
			_, _ = w.Write([]byte(" "))
			_, _ = w.Write([]byte(t.ins))
		}
		_, _ = w.Write([]byte("</code>"))
		if len(t.desc) > 0 {
			_, _ = w.Write([]byte(" "))
			_, _ = w.Write([]byte(htmlEscape(t.desc, false)))
		}
		_, _ = w.Write([]byte("</li>"))
	case format == DocgenFormatHTML && !compact:
		_, _ = w.Write([]byte("<h3>" + t.name + "</h3>"))
		if len(t.alias) > 0 {
			_, _ = w.Write([]byte("<block>Alias: <code>" + t.alias + "</code></block>"))
		}
		if len(t.params) > 0 {
			_, _ = w.Write([]byte("<div>Params:<ul>"))
			for j := 0; j < len(t.params); j++ {
				param := &t.params[j]
				_, _ = w.Write([]byte("<li><code>" + param.name + "</code>"))
				if len(param.desc) > 0 {
					_, _ = w.Write([]byte(" " + htmlEscape(param.desc, true)))
				}
				_, _ = w.Write([]byte("</li>"))
			}
			_, _ = w.Write([]byte("</ul></div>"))
		}

		if len(t.desc) > 0 {
			_, _ = w.Write([]byte("<p>"))
			_, _ = w.Write([]byte(htmlEscape(t.desc, true)))
			_, _ = w.Write([]byte("</p>"))
		}

		if len(t.note) > 0 {
			_, _ = w.Write([]byte("<blockquote>"))
			_, _ = w.Write([]byte(htmlEscape(t.note, false)))
			_, _ = w.Write([]byte("</blockquote>"))
		}

		if len(t.example) > 0 {
			_, _ = w.Write([]byte("<div>Example:<pre>"))
			_, _ = w.Write([]byte(htmlEscape(t.example, false)))
			_, _ = w.Write([]byte("</pre></div>"))
		}
	}

	return nil
}

var reDGHTMLCode = regexp.MustCompile("`([^`]+)`")

type (
	docgenJSON struct {
		Name      string `json:"name,omitempty"`
		Alias     string `json:"alias,omitempty"`
		Type      string `json:"type,omitempty"`
		Desc      string `json:"desc,omitempty"`
		Note      string `json:"note,omitempty"`
		Example   string `json:"example,omitempty"`
		Inspector string `json:"inspector,omitempty"`

		Params []docgenParamJSON `json:"params,omitempty"`
	}
	docgenParamJSON struct {
		Name string `json:"name,omitempty"`
		Desc string `json:"desc,omitempty"`
	}
	docgenContainerJSON struct {
		Modifiers              []docgenJSON `json:"modifiers,omitempty"`
		ConditionHelpers       []docgenJSON `json:"conditionHelpers,omitempty"`
		ConditionOKHelpers     []docgenJSON `json:"conditionOKHelpers,omitempty"`
		EmptyCheckHelpers      []docgenJSON `json:"emptyCheckHelpers,omitempty"`
		GlobalVariables        []docgenJSON `json:"globalVariables,omitempty"`
		VariableInspectorPairs []docgenJSON `json:"variableInspectorPairs,omitempty"`
	}
)
