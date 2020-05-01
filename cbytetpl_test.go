package cbytetpl

import (
	"bytes"
	"testing"

	"github.com/koykov/inspector/testobj"
	"github.com/koykov/inspector/testobj_ins"
)

var (
	buf bytes.Buffer

	user = &testobj.TestObject{
		Name:   []byte("John"),
		Status: 78,
		Finance: &testobj.TestFinance{
			Balance:  9000.015,
			AllowBuy: false,
		},
	}
	ins testobj_ins.TestObjectInspector

	tplRaw = []byte(`<h1>Raw template<h1><p>Lorem ipsum dolor sit amet, ...</p>`)

	tplSimple = []byte(`<h1>Welcome, {%= user.Name %}!</h1>
<p>Status: {%= user.Status %}</p>
<p>Your balance: {%= user.Finance.Balance %}; buy allowance: {%= user.Finance.AllowBuy %}</p>`)
	expectSimple = []byte(`<h1>Welcome, John!</h1><p>Status: 78</p><p>Your balance: 9000.015; buy allowance: false</p>`)
)

func pretest() {
	tree, _ := Parse(tplRaw, false)
	RegisterTpl("tplRaw", tree)

	tree, _ = Parse(tplSimple, false)
	RegisterTpl("tplSimple", tree)
}

func TestTplRaw(t *testing.T) {
	pretest()

	result, err := Render("tplRaw", nil)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, tplRaw) {
		t.Error("raw tpl mismatch")
	}
}

func TestTplSimple(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render("tplSimple", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectSimple) {
		t.Error("simple tpl mismatch")
	}
}

func BenchmarkTplSimple(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		ctx.Set("user", user, &ins)
		buf.Reset()
		err := RenderTo(&buf, "tplSimple", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectSimple) {
			b.Error("simple tpl mismatch")
		}
		CP.Put(ctx)
	}
}
