package dyntpl

import (
	"bytes"
	"testing"
)

var (
	tplModDefault       = []byte(`Cost is: {%= user.Cost|default(999.99) %} USD`)
	tplModDefaultStatic = []byte(`{% ctx defaultCost = 999.99 %}Cost is: {%= user.Cost|default(defaultCost) %} USD`)
	expectModDefault    = []byte(`Cost is: 999.99 USD`)

	tplModDefault1    = []byte(`<span style="background-color:{%= user.ProfileColor|default("#fff") %}"></span>"`)
	expectModDefault1 = []byte(`<span style="background-color:#fff"></span>"`)

	tplModJSONEscape       = []byte(`{"id":"foo","name":"{%= userName|jsonEscape %}"}`)
	tplModJSONEscapeShort  = []byte(`{"id":"foo","name":"{%j= userName %}"}`)
	tplModJSONQuoteShort   = []byte(`{"id":"foo","name":{%q= userName %}}`)
	expectModJSON          = []byte(`{"id":"foo","name":"Foo\"bar"}`)
	tplModJSONEscapeDbl    = []byte(`{"obj":"{% jsonquote %}{"inner":"{%j= valueWithQuotes %}"}{% endjsonquote %}"}`)
	expectModJSONEscapeDbl = []byte(`{"obj":"{\"inner\":\"He said: \\\"welcome friend\\\"\"}"}`)

	tplModURLEncode     = []byte(`<a href="https://redir.com/{%u= url %}">go to >>></a>`)
	expectModURLEncode  = []byte(`<a href="https://redir.com/https%3A%2F%2Fgolang.org%2Fsrc%2Fnet%2Furl%2Furl.go%23L100">go to >>></a>`)
	tplModURLEncode2    = []byte(`<a href="https://redir.com/{%uu= url %}">go to >>></a>`)
	expectModURLEncode2 = []byte(`<a href="https://redir.com/https%253A%252F%252Fgolang.org%252Fsrc%252Fnet%252Furl%252Furl.go%2523L100">go to >>></a>`)
	tplModURLEncode3    = []byte(`<a href="https://redir.com/{%uuu= url %}">go to >>></a>`)
	expectModURLEncode3 = []byte(`<a href="https://redir.com/https%25253A%25252F%25252Fgolang.org%25252Fsrc%25252Fnet%25252Furl%25252Furl.go%252523L100">go to >>></a>`)

	tplModHtmlEscape      = []byte(`<a href="https://golang.org/" title="{%= title|htmlEscape %}">{%= text|he %}</a>`)
	tplModHtmlEscapeShort = []byte(`<a href="https://golang.org/" title="{%h= title %}">{%h= text %}</a>`)
	expectModHtml         = []byte(`<a href="https://golang.org/" title="&lt;h1&gt;Go is an open source programming language that makes it easy to build &lt;strong&gt;simple&lt;strong&gt;, &lt;strong&gt;reliable&lt;/strong&gt;, and &lt;strong&gt;efficient&lt;/strong&gt; software.&lt;/h1&gt;">Visit &gt;</a>`)

	tplModLinkEscape    = []byte(`<a href="{%l= link %}">`)
	expectModLinkEscape = []byte(`<a href="http://x.com/link-with-\"-and+space-symbol">`)

	tplModIfThen        = []byte(`{%= allow|ifThen("You're allow to buy!") %}`)
	expectModIfThen     = []byte(`You're allow to buy!`)
	tplModIfThenElse    = []byte(`Welcome, {%= logged|ifThenElse(userName, "anonymous") %}!`)
	expectModIfThenElse = []byte(`Welcome, foobar!`)

	tplModRound    = []byte(`Price 1: {%= f0|round %}; Price 2: {%= f1|roundPrec(3) %}; Price 3: {%= f2|ceil %}; Price 4: {%F.3= f3 %}; Price 5: {%= f4|floor %}; Price 6: {%f.3= f5 %}`)
	expectModRound = []byte(`Price 1: 7; Price 2: 3.141; Price 3: 12; Price 4: 56.688; Price 5: 67; Price 6: 20.214`)
)

type modStage struct {
	key  string
	args map[string]interface{}
	fn   func(t *testing.T, st *modStage)
}

func TestMod(t *testing.T) {
	loadStages()

	modStages := []modStage{
		{key: "modDefault"},
		{key: "modDefaultStatic"},
		{key: "modDefault1"},
		{key: "modJSONEscape", fn: testModWA, args: map[string]interface{}{"userName": `Foo"bar`}},
		{key: "modJSONEscapeShort", fn: testModWA, args: map[string]interface{}{"userName": `Foo"bar`}},
		{key: "modJSONEscapeDbl", fn: testModWA, args: map[string]interface{}{"valueWithQuotes": `He said: "welcome friend"`}},
		{key: "modJSONQuoteShort", fn: testModWA, args: map[string]interface{}{"userName": `Foo"bar`}},
		{key: "modHtmlEscape", fn: testModWA, args: map[string]interface{}{
			"title": `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`,
			"text":  `Visit >`,
		}},
		{key: "modHtmlEscapeShort", fn: testModWA, args: map[string]interface{}{
			"title": `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`,
			"text":  `Visit >`,
		}},
		{key: "modLinkEscape", fn: testModWA, args: map[string]interface{}{"link": `http://x.com/link-with-"-and space-symbol`}},
		{key: "modURLEncode", fn: testModWA, args: map[string]interface{}{"url": `https://golang.org/src/net/url/url.go#L100`}},
		{key: "modURLEncode2", fn: testModWA, args: map[string]interface{}{"url": `https://golang.org/src/net/url/url.go#L100`}},
		{key: "modURLEncode3", fn: testModWA, args: map[string]interface{}{"url": `https://golang.org/src/net/url/url.go#L100`}},
		{key: "modIfThen", fn: testModWA, args: map[string]interface{}{"allow": true}},
		{key: "modIfThenElse", fn: testModWA, args: map[string]interface{}{"logged": true, "userName": "foobar"}},
		{key: "modRound", fn: testModWA, args: map[string]interface{}{
			"f0": 7.243242,
			"f1": 3.1415,
			"f2": 11.39,
			"f3": 56.68734,
			"f4": 67.999,
			"f5": 20.214999,
		}},
	}

	for _, s := range modStages {
		t.Run(s.key, func(t *testing.T) {
			if s.fn == nil {
				s.fn = testMod
			}
			s.fn(t, &s)
		})
	}
}

func testMod(t *testing.T, st *modStage) {
	testTpl(t, st.key)
}

func testModWA(t *testing.T, st *modStage) {
	st1 := getStage(st.key)
	if st1 == nil {
		t.Error("stage not found")
		return
	}

	ctx := NewCtx()
	for k, v := range st.args {
		ctx.SetStatic(k, v)
	}
	result, err := Render(st.key, ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, st1.expect) {
		t.Errorf("%s mismatch", st.key)
	}
	if !bytes.Equal(result, st1.expect) {
		t.Errorf("%s mismatch", st.key)
	}
}

func BenchmarkTplModDefault(b *testing.B) {
	benchBase(b, "tplModDefault", expectModDefault, "mod def tpl mismatch")
}

func BenchmarkTplModDefaultStatic(b *testing.B) {
	benchBase(b, "tplModDefaultStatic", expectModDefault, "mod def static tpl mismatch")
}

func BenchmarkTplModArgs(b *testing.B) {
	benchBase(b, "tplModDefault1", expectModDefault1, "mod def (hex color arg) static tpl mismatch")
}

func BenchmarkTplModJSONQuote(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("userName", `Foo"bar`)
		err := Write(&buf, "tplModJSONQuoteShort", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModJSON) {
			b.Error("json quote tpl mismatch")
		}
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplModJSONEscape(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("userName", `Foo"bar`)
		err := Write(&buf, "tplModJSONEscapeShort", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModJSON) {
			b.Error("json escape tpl mismatch")
		}
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplModJSONEscapeDbl(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("valueWithQuotes", `He said: "welcome friend"`)
		err := Write(&buf, "tplModJSONEscapeDbl", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModJSONEscapeDbl) {
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
		err := Write(&buf, "tplModHtmlEscapeShort", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModHtml) {
			b.Error("html escape tpl mismatch")
		}
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplModLinkEscape(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("link", `http://x.com/link-with-"-and space-symbol`)
		err := Write(&buf, "tplModLinkEscape", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModLinkEscape) {
			b.Error("link escape tpl mismatch")
		}
		ReleaseCtx(ctx)
	}
}

func benchModURLEncode(b *testing.B, tplID string, expect []byte, failMsg string) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("url", `https://golang.org/src/net/url/url.go#L100`)
		err := Write(&buf, tplID, ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expect) {
			b.Error(failMsg)
		}
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplModURLEncode(b *testing.B) {
	benchModURLEncode(b, "tplModURLEncode", expectModURLEncode, "url encode tpl mismatch")
}

func BenchmarkTplModURLEncode2(b *testing.B) {
	benchModURLEncode(b, "tplModURLEncode2", expectModURLEncode2, "url encode 2 tpl mismatch")
}

func BenchmarkTplModURLEncode3(b *testing.B) {
	benchModURLEncode(b, "tplModURLEncode3", expectModURLEncode3, "url encode 3 tpl mismatch")
}

func BenchmarkTplModIfThen(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.SetStatic("allow", true)
		err := Write(&buf, "tplModIfThen", ctx)
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
		err := Write(&buf, "tplModIfThenElse", ctx)
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
		err := Write(&buf, "tplModRound", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectModRound) {
			b.Error("round tpl mismatch")
		}
		ReleaseCtx(ctx)
	}
}
