package dyntpl

import (
	"bytes"
	"testing"
)

var (
	tplModDef       = []byte(`Cost is: {%= user.Cost|default(999.99) %} USD`)
	tplModDefStatic = []byte(`{% ctx defaultCost = 999.99 %}Cost is: {%= user.Cost|default(defaultCost) %} USD`)
	expectModDef    = []byte(`Cost is: 999.99 USD`)

	tplModJsonQ    = []byte(`{"id":"foo","name":"{%= userName|jsonQuote %}"}`)
	expectModJsonQ = []byte(`{"id":"foo","name":"Foo\"bar"}`)

	tplModHtmlE    = []byte(`<a href="https://golang.org/" title="{%= title|htmlEscape %}">{%= text|he %}</a>`)
	expectModHtmlE = []byte(`<a href="https://golang.org/" title="&lt;h1&gt;Go is an open source programming language that makes it easy to build &lt;strong&gt;simple&lt;strong&gt;, &lt;strong&gt;reliable&lt;/strong&gt;, and &lt;strong&gt;efficient&lt;/strong&gt; software.&lt;/h1&gt;">Visit &gt;</a>`)
)

func TestTplModDef(t *testing.T) {
	testBase(t, "tplModDef", expectModDef, "mod def tpl mismatch")
}

func TestTplModDefStatic(t *testing.T) {
	testBase(t, "tplModDefStatic", expectModDef, "mod def static tpl mismatch")
}

func BenchmarkTplModDef(b *testing.B) {
	benchBase(b, "tplModDef", expectModDef, "mod def tpl mismatch")
}

func BenchmarkTplModDefStatic(b *testing.B) {
	benchBase(b, "tplModDefStatic", expectModDef, "mod def static tpl mismatch")
}

func TestTplModJsonQuote(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.SetStatic("userName", `Foo"bar`)
	result, err := Render("tplModJsonQ", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectModJsonQ) {
		t.Error("json quote tpl mismatch")
	}
}

func TestTplModHtmlEscape(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.SetStatic("title", `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`)
	ctx.SetStatic("text", `Visit >`)
	result, err := Render("tplModHtmlE", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectModHtmlE) {
		t.Error("html escape tpl mismatch")
	}
}

func BenchmarkTplModJsonQuote(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		buf.Reset()
		ctx.SetStatic("userName", `Foo"bar`)
		err := RenderTo(&buf, "tplModJsonQ", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModJsonQ) {
			b.Error("json quote tpl mismatch")
		}
		CP.Put(ctx)
	}
}

func BenchmarkTplModHtmlEscape(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		buf.Reset()
		ctx.SetStatic("title", `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`)
		ctx.SetStatic("text", `Visit >`)
		err := RenderTo(&buf, "tplModHtmlE", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModHtmlE) {
			b.Error("html escape tpl mismatch")
		}
		CP.Put(ctx)
	}
}
