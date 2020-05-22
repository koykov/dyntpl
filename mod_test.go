package dyntpl

import (
	"bytes"
	"testing"
)

var (
	tplModDef       = []byte(`Cost is: {%= user.Cost|default(999.99) %} USD`)
	tplModDefStatic = []byte(`{% ctx defaultCost = 999.99 %}Cost is: {%= user.Cost|default(defaultCost) %} USD`)
	expectModDef    = []byte(`Cost is: 999.99 USD`)

	tplModJsonEscape      = []byte(`{"id":"foo","name":"{%= userName|jsonEscape %}"}`)
	tplModJsonEscapeShort = []byte(`{"id":"foo","name":"{%j= userName %}"}`)
	tplModJsonQuoteShort  = []byte(`{"id":"foo","name":{%q= userName %}}`)
	expectModJson         = []byte(`{"id":"foo","name":"Foo\"bar"}`)

	tplModHtmlEscape      = []byte(`<a href="https://golang.org/" title="{%= title|htmlEscape %}">{%= text|he %}</a>`)
	tplModHtmlEscapeShort = []byte(`<a href="https://golang.org/" title="{%h= title %}">{%h= text %}</a>`)
	expectModHtml         = []byte(`<a href="https://golang.org/" title="&lt;h1&gt;Go is an open source programming language that makes it easy to build &lt;strong&gt;simple&lt;strong&gt;, &lt;strong&gt;reliable&lt;/strong&gt;, and &lt;strong&gt;efficient&lt;/strong&gt; software.&lt;/h1&gt;">Visit &gt;</a>`)

	tplModIfThen        = []byte(`{%= allow|ifThen("You're allow to buy!") %}`)
	expectModIfThen     = []byte(`You're allow to buy!`)
	tplModIfThenElse    = []byte(`Welcome, {%= logged|ifThenElse(userName, "anonymous") %}!`)
	expectModIfThenElse = []byte(`Welcome, foobar!`)
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
	result, err := Render("tplModJsonQuoteShort", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectModJson) {
		t.Error("json quote tpl mismatch")
	}
}

func TestTplModJsonEscape(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.SetStatic("userName", `Foo"bar`)
	result, err := Render("tplModJsonEscapeShort", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectModJson) {
		t.Error("json escape tpl mismatch")
	}
}

func TestTplModHtmlEscape(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.SetStatic("title", `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`)
	ctx.SetStatic("text", `Visit >`)
	result, err := Render("tplModHtmlEscape", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectModHtml) {
		t.Error("html escape tpl mismatch")
	}
}

func TestTplModIfThen(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.SetStatic("allow", true)
	result, err := Render("tplModIfThen", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectModIfThen) {
		t.Error("ifThen tpl mismatch")
	}
}

func TestTplModIfThenElse(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.SetStatic("logged", true)
	ctx.SetStatic("userName", "foobar")
	result, err := Render("tplModIfThenElse", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectModIfThenElse) {
		t.Error("ifThenElse tpl mismatch")
	}
}

func BenchmarkTplModJsonQuote(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		buf.Reset()
		ctx.SetStatic("userName", `Foo"bar`)
		err := RenderTo(&buf, "tplModJsonQuoteShort", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModJson) {
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
		err := RenderTo(&buf, "tplModHtmlEscapeShort", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModHtml) {
			b.Error("html escape tpl mismatch")
		}
		CP.Put(ctx)
	}
}

func BenchmarkTplModIfThen(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		buf.Reset()
		ctx.SetStatic("allow", true)
		err := RenderTo(&buf, "tplModIfThen", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModIfThen) {
			b.Error("ifThen tpl mismatch")
		}
		CP.Put(ctx)
	}
}

func BenchmarkTplModIfThenElse(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		buf.Reset()
		ctx.SetStatic("logged", true)
		ctx.SetStatic("userName", "foobar")
		err := RenderTo(&buf, "tplModIfThenElse", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModIfThenElse) {
			b.Error("ifThenElse tpl mismatch")
		}
		CP.Put(ctx)
	}
}
