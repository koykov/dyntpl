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
	switch {
	case format == DocgenFormatMarkdown && compact == true:
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
	case format == DocgenFormatMarkdown && compact == false:
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
