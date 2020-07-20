package dyntpl

import (
	"bytes"
	"testing"
)

var (
	tplModDef       = []byte(`Cost is: {%= user.Cost|default(999.99) %} USD`)
	tplModDefStatic = []byte(`{% ctx defaultCost = 999.99 %}Cost is: {%= user.Cost|default(defaultCost) %} USD`)
	expectModDef    = []byte(`Cost is: 999.99 USD`)

	tplModJsonEscape       = []byte(`{"id":"foo","name":"{%= userName|jsonEscape %}"}`)
	tplModJsonEscapeShort  = []byte(`{"id":"foo","name":"{%j= userName %}"}`)
	tplModJsonQuoteShort   = []byte(`{"id":"foo","name":{%q= userName %}}`)
	expectModJson          = []byte(`{"id":"foo","name":"Foo\"bar"}`)
	tplModJsonEscapeDbl    = []byte(`{"obj":"{% jsonquote %}{"inner":"{%j= valueWithQuotes %}"}{% endjsonquote %}"}`)
	expectModJsonEscapeDbl = []byte(`{"obj":"{\"inner\":\"He said: \\\"welcome friend\\\"\"}"}`)

	tplModUrlEncode    = []byte(`<a href="https://redir.com/{%u= url %}">go to >>></a>`)
	expectModUrlEncode = []byte(`<a href="https://redir.com/https%3A%2F%2Fgolang.org%2Fsrc%2Fnet%2Furl%2Furl.go%23L100">go to >>></a>`)

	tplModHtmlEscape      = []byte(`<a href="https://golang.org/" title="{%= title|htmlEscape %}">{%= text|he %}</a>`)
	tplModHtmlEscapeShort = []byte(`<a href="https://golang.org/" title="{%h= title %}">{%h= text %}</a>`)
	expectModHtml         = []byte(`<a href="https://golang.org/" title="&lt;h1&gt;Go is an open source programming language that makes it easy to build &lt;strong&gt;simple&lt;strong&gt;, &lt;strong&gt;reliable&lt;/strong&gt;, and &lt;strong&gt;efficient&lt;/strong&gt; software.&lt;/h1&gt;">Visit &gt;</a>`)

	tplModIfThen        = []byte(`{%= allow|ifThen("You're allow to buy!") %}`)
	expectModIfThen     = []byte(`You're allow to buy!`)
	tplModIfThenElse    = []byte(`Welcome, {%= logged|ifThenElse(userName, "anonymous") %}!`)
	expectModIfThenElse = []byte(`Welcome, foobar!`)

	tplModRound    = []byte(`Price 1: {%= f0|round %}; Price 2: {%= f1|roundPrec(3) %}; Price 3: {%= f2|ceil %}; Price 4: {%F.3= f3 %}; Price 5: {%= f4|floor %}; Price 6: {%f.3= f5 %}`)
	expectModRound = []byte(`Price 1: 7; Price 2: 3.141; Price 3: 12; Price 4: 56.688; Price 5: 67; Price 6: 20.214`)
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

func TestTplModJsonEscapeDbl(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.SetStatic("valueWithQuotes", `He said: "welcome friend"`)
	result, err := Render("tplModJsonEscapeDbl", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectModJsonEscapeDbl) {
		t.Error("json escape (dbl) tpl mismatch")
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

func TestTplModUrlEncode(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.SetStatic("url", `https://golang.org/src/net/url/url.go#L100`)
	result, err := Render("tplModUrlEncode", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectModUrlEncode) {
		t.Error("url encode tpl mismatch")
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

func TestTplModRoundPrec(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.SetStatic("f0", 7.243242)
	ctx.SetStatic("f1", 3.1415)
	ctx.SetStatic("f2", 11.39)
	ctx.SetStatic("f3", 56.68734)
	ctx.SetStatic("f4", 67.999)
	ctx.SetStatic("f5", 20.214999)
	result, err := Render("tplModRound", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectModRound) {
		t.Error("round tpl mismatch")
	}
}

func BenchmarkTplModJsonQuote(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("userName", `Foo"bar`)
		err := RenderTo(&buf, "tplModJsonQuoteShort", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModJson) {
			b.Error("json quote tpl mismatch")
		}
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplModJsonEscapeDbl(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("valueWithQuotes", `He said: "welcome friend"`)
		err := RenderTo(&buf, "tplModJsonEscapeDbl", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModJsonEscapeDbl) {
			b.Error("json escape (dbl) tpl mismatch")
		}
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplModHtmlEscape(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
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
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplModUrlEncode(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("url", `https://golang.org/src/net/url/url.go#L100`)
		err := RenderTo(&buf, "tplModUrlEncode", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModUrlEncode) {
			b.Error("url encode tpl mismatch")
		}
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplModIfThen(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("allow", true)
		err := RenderTo(&buf, "tplModIfThen", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModIfThen) {
			b.Error("ifThen tpl mismatch")
		}
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplModIfThenElse(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
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
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplModRound(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("f0", 7.243242)
		ctx.SetStatic("f1", 3.1415)
		ctx.SetStatic("f2", 11.39)
		ctx.SetStatic("f3", 56.68734)
		ctx.SetStatic("f4", 67.999)
		ctx.SetStatic("f5", 20.214999)
		err := RenderTo(&buf, "tplModRound", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModRound) {
			b.Error("round tpl mismatch")
		}
		ReleaseCtx(ctx)
	}
}
