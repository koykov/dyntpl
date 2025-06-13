package dyntpl

import (
	"bytes"
	"testing"
)

type modArgs map[string]any

func TestMod(t *testing.T) {
	t.Run("modDefault", testMod)
	t.Run("modDefaultStatic", testMod)
	t.Run("modDefault1", testMod)
	t.Run("modJSONEscape", func(t *testing.T) { testModWA(t, modArgs{"userName": `Foo"bar`}) })
	t.Run("modJSONEscapeShort", func(t *testing.T) { testModWA(t, modArgs{"userName": `Foo"bar`}) })
	t.Run("modJSONEscapeDbl", func(t *testing.T) {
		testModWA(t, modArgs{"valueWithQuotes": `He said: "welcome friend"`})
	})
	t.Run("modJSONQuoteShort", func(t *testing.T) { testModWA(t, modArgs{"userName": `Foo"bar`}) })
	t.Run("modJSONQuoteNoesc", func(t *testing.T) { testModWA(t, modArgs{"userName": `Foo"bar`}) })
	t.Run("modHtmlEscape", func(t *testing.T) {
		testModWA(t, modArgs{
			"title": `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`,
			"text":  `Visit >`,
		})
	})
	t.Run("modHtmlEscapeShort", func(t *testing.T) {
		testModWA(t, modArgs{
			"title": `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`,
			"text":  `Visit >`,
		})
	})
	t.Run("modAttrEscape", func(t *testing.T) { testModWA(t, modArgs{"var1": "foo&<>\"'`!@$%()=+{}[]#;bar"}) })
	t.Run("modAttrEscapeMB", func(t *testing.T) { testModWA(t, modArgs{"var1": "Привет мир😀!"}) })
	t.Run("modCSSEscape", func(t *testing.T) { testModWA(t, modArgs{"var1": "<>'\"&Ā,._aAzZ09 !😀"}) })
	t.Run("modJSEscape", func(t *testing.T) { testModWA(t, modArgs{"var1": "<>'\"&/,._aAzZ09 Ā😀"}) })
	t.Run("modJSEscapeMB", func(t *testing.T) { testModWA(t, modArgs{"var1": "Привет мир😀!"}) })
	t.Run("modLinkEscape", func(t *testing.T) {
		testModWA(t, modArgs{"link": `http://x.com/link-with-"-and space-symbol`})
	})
	t.Run("modURLEncode", func(t *testing.T) {
		testModWA(t, modArgs{"url": `https://golang.org/src/net/url/url.go#L100`})
	})
	t.Run("modURLEncode2", func(t *testing.T) {
		testModWA(t, modArgs{"url": `https://golang.org/src/net/url/url.go#L100`})
	})
	t.Run("modURLEncode3", func(t *testing.T) {
		testModWA(t, modArgs{"url": `https://golang.org/src/net/url/url.go#L100`})
	})
	t.Run("modIfThen", func(t *testing.T) { testModWA(t, modArgs{"allow": true}) })
	t.Run("modIfThenElse", func(t *testing.T) { testModWA(t, modArgs{"logged": true, "userName": "foobar"}) })
	t.Run("modRound", func(t *testing.T) {
		testModWA(t, modArgs{
			"f0": 7.243242,
			"f1": 3.1415,
			"f2": 11.39,
			"f3": 56.68734,
			"f4": 67.999,
			"f5": 20.214999,
		})
	})
}

func testMod(t *testing.T) {
	testTpl(t)
}

func testModWA(t *testing.T, args modArgs) {
	key := getTBName(t)
	st := getStage(key)
	if st == nil {
		t.Error("stage not found")
		return
	}

	ctx := NewCtx()
	for k, v := range args {
		ctx.SetStatic(k, v)
	}
	result, err := Render(key, ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, st.expect) {
		t.Errorf("%s mismatch: need %s\ngot %s", key, st.expect, result)
	}
}

func BenchmarkMod(b *testing.B) {
	b.Run("modDefault", benchMod)
	b.Run("modDefaultStatic", benchMod)
	b.Run("modDefault1", benchMod)
	b.Run("modJSONEscape", func(b *testing.B) { benchModWA(b, modArgs{"userName": `Foo"bar`}) })
	b.Run("modJSONEscapeShort", func(b *testing.B) { benchModWA(b, modArgs{"userName": `Foo"bar`}) })
	b.Run("modJSONEscapeDbl", func(b *testing.B) {
		benchModWA(b, modArgs{"valueWithQuotes": `He said: "welcome friend"`})
	})
	b.Run("modJSONQuoteShort", func(b *testing.B) { benchModWA(b, modArgs{"userName": `Foo"bar`}) })
	b.Run("modJSONQuoteNoesc", func(b *testing.B) { benchModWA(b, modArgs{"userName": `Foo"bar`}) })
	b.Run("modHtmlEscape", func(b *testing.B) {
		benchModWA(b, modArgs{
			"title": `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`,
			"text":  `Visit >`,
		})
	})
	b.Run("modHtmlEscapeShort", func(b *testing.B) {
		benchModWA(b, modArgs{
			"title": `<h1>Go is an open source programming language that makes it easy to build <strong>simple<strong>, <strong>reliable</strong>, and <strong>efficient</strong> software.</h1>`,
			"text":  `Visit >`,
		})
	})
	b.Run("modAttrEscape", func(b *testing.B) { benchModWA(b, modArgs{"var1": "foo&<>\"'`!@$%()=+{}[]#;bar"}) })
	b.Run("modAttrEscapeMB", func(b *testing.B) { benchModWA(b, modArgs{"var1": "Привет мир😀!"}) })
	b.Run("modCSSEscape", func(b *testing.B) { benchModWA(b, modArgs{"var1": "<>'\"&Ā,._aAzZ09 !😀"}) })
	b.Run("modJSEscape", func(b *testing.B) { benchModWA(b, modArgs{"var1": "<>'\"&/,._aAzZ09 Ā😀"}) })
	b.Run("modJSEscapeMB", func(b *testing.B) { benchModWA(b, modArgs{"var1": "Привет мир😀!"}) })
	b.Run("modLinkEscape", func(b *testing.B) {
		benchModWA(b, modArgs{"link": `http://x.com/link-with-"-and space-symbol`})
	})
	b.Run("modURLEncode", func(b *testing.B) {
		benchModWA(b, modArgs{"url": `https://golang.org/src/net/url/url.go#L100`})
	})
	b.Run("modURLEncode2", func(b *testing.B) {
		benchModWA(b, modArgs{"url": `https://golang.org/src/net/url/url.go#L100`})
	})
	b.Run("modURLEncode3", func(b *testing.B) {
		benchModWA(b, modArgs{"url": `https://golang.org/src/net/url/url.go#L100`})
	})
	b.Run("modIfThen", func(b *testing.B) { benchModWA(b, modArgs{"allow": true}) })
	b.Run("modIfThenElse", func(b *testing.B) { benchModWA(b, modArgs{"logged": true, "userName": "foobar"}) })
	b.Run("modRound", func(b *testing.B) {
		benchModWA(b, modArgs{
			"f0": 7.243242,
			"f1": 3.1415,
			"f2": 11.39,
			"f3": 56.68734,
			"f4": 67.999,
			"f5": 20.214999,
		})
	})
}

func benchMod(b *testing.B) {
	benchTpl(b)
}

func benchModWA(b *testing.B, args modArgs) {
	key := getTBName(b)

	st := getStage(key)
	if st == nil {
		b.Error("stage not found")
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
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), st.expect) {
			b.Errorf("%s mismatch", key)
		}
		ReleaseCtx(ctx)
	}
}
