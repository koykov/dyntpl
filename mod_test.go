package dyntpl

import (
	"bytes"
	"testing"
)

type modStage struct {
	key  string
	args map[string]interface{}
	fn   func(tb testing.TB, st *modStage)
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

func testMod(tb testing.TB, st *modStage) {
	testTpl(tb, st.key)
}

func testModWA(tb testing.TB, st *modStage) {
	st1 := getStage(st.key)
	if st1 == nil {
		tb.Error("stage not found")
		return
	}

	ctx := NewCtx()
	for k, v := range st.args {
		ctx.SetStatic(k, v)
	}
	result, err := Render(st.key, ctx)
	if err != nil {
		tb.Error(err)
	}
	if !bytes.Equal(result, st1.expect) {
		tb.Errorf("%s mismatch", st.key)
	}
}

func BenchmarkMod(b *testing.B) {
	loadStages()

	b.Run("modDefault", func(b *testing.B) { benchMod(b) })
	b.Run("modDefaultStatic", func(b *testing.B) { benchMod(b) })
	b.Run("modDefault1", func(b *testing.B) { benchMod(b) })
	b.Run("modJSONEscape", func(b *testing.B) { benchModWA(b, map[string]interface{}{"userName": `Foo"bar`}) })
	b.Run("modJSONEscapeShort", func(b *testing.B) { benchModWA(b, map[string]interface{}{"userName": `Foo"bar`}) })
	b.Run("modJSONEscapeDbl", func(b *testing.B) {
		benchModWA(b, map[string]interface{}{"valueWithQuotes": `He said: "welcome friend"`})
	})
	b.Run("modJSONQuoteShort", func(b *testing.B) { benchModWA(b, map[string]interface{}{"userName": `Foo"bar`}) })
	b.Run("modHtmlEscape", func(b *testing.B) {
		benchModWA(b, map[string]interface{}{
			"title": `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`,
			"text":  `Visit >`,
		})
	})
	b.Run("modHtmlEscapeShort", func(b *testing.B) {
		benchModWA(b, map[string]interface{}{
			"title": `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`,
			"text":  `Visit >`,
		})
	})
	b.Run("modLinkEscape", func(b *testing.B) {
		benchModWA(b, map[string]interface{}{"link": `http://x.com/link-with-"-and space-symbol`})
	})
	b.Run("modURLEncode", func(b *testing.B) {
		benchModWA(b, map[string]interface{}{"url": `https://golang.org/src/net/url/url.go#L100`})
	})
	b.Run("modURLEncode2", func(b *testing.B) {
		benchModWA(b, map[string]interface{}{"url": `https://golang.org/src/net/url/url.go#L100`})
	})
	b.Run("modURLEncode3", func(b *testing.B) {
		benchModWA(b, map[string]interface{}{"url": `https://golang.org/src/net/url/url.go#L100`})
	})
	b.Run("modIfThen", func(b *testing.B) { benchModWA(b, map[string]interface{}{"allow": true}) })
	b.Run("modIfThenElse", func(b *testing.B) { benchModWA(b, map[string]interface{}{"logged": true, "userName": "foobar"}) })
	b.Run("modRound", func(b *testing.B) {
		benchModWA(b, map[string]interface{}{
			"f0": 7.243242,
			"f1": 3.1415,
			"f2": 11.39,
			"f3": 56.68734,
			"f4": 67.999,
			"f5": 20.214999,
		})
	})
}

func benchMod(tb testing.TB) {
	benchTpl(tb)
}

func benchModWA(tb testing.TB, args map[string]interface{}) {
	b := interface{}(tb).(*testing.B)
	key := getTBName(b)

	st := getStage(key)
	if st == nil {
		tb.Error("stage not found")
		return
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		for k, v := range args {
			ctx.SetStatic(k, v)
		}
		buf.Reset()
		err := Write(&buf, key, ctx)
		if err != nil {
			tb.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), st.expect) {
			tb.Errorf("%s mismatch", key)
		}
		ReleaseCtx(ctx)
	}
}
