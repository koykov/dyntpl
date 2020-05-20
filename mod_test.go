package dyntpl

import (
	"bytes"
	"testing"
)

var (
	tplModDef       = []byte(`Cost is: {%= user.Cost|default(999.99) %} USD`)
	tplModDefStatic = []byte(`{% ctx defaultCost = 999.99 %}Cost is: {%= user.Cost|default(defaultCost) %} USD`)
	expectModDef    = []byte(`Cost is: 999.99 USD`)

	tplModJsonQ    = []byte(`{"id":"foo","name":{%= userName|jsonQuote pfx " sfx " %}}`)
	expectModJsonQ = []byte(`{"id":"foo","name":"Foo\"bar"}`)
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

func TestTplModJsonQ(t *testing.T) {
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

func BenchmarkTplModJsonQ(b *testing.B) {
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
