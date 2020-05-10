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
		Id:     "115",
		Name:   []byte("John"),
		Status: 78,
		Flags: testobj.TestFlag{
			"export": 17,
			"ro":     4,
			"rw":     7,
			"Valid":  1,
		},
		Finance: &testobj.TestFinance{
			Balance:  9000.015,
			AllowBuy: false,
			History: []testobj.TestHistory{
				{
					152354345634,
					14.345241,
					[]byte("pay for domain"),
				},
				{
					153465345246,
					-3.0000342543,
					[]byte("got refund"),
				},
				{
					156436535640,
					2325242534.35324523,
					[]byte("maintenance"),
				},
			},
		},
	}
	ins testobj_ins.TestObjectInspector

	tplRaw = []byte(`<h1>Raw template<h1><p>Lorem ipsum dolor sit amet, ...</p>`)

	tplSimple = []byte(`<h1>Welcome, {%= user.Name %}!</h1>
<p>Status: {%= user.Status %}</p>
<p>Your balance: {%= user.Finance.Balance %}; buy allowance: {%= user.Finance.AllowBuy %}</p>`)
	expectSimple = []byte(`<h1>Welcome, John!</h1><p>Status: 78</p><p>Your balance: 9000.015; buy allowance: false</p>`)

	tplCond = []byte(`<h2>Status</h2><p>
{% if user.Status >= 60 %}Privileged user, your balance: {%= user.Finance.Balance %}.
{% else %}You don't have enough privileges.{% endif %}</p>`)
	tplCondNoStatic = []byte(`<h2>Status</h2><p>
{%ctx permissionLimit = 60 %}
{% if user.Status >= permissionLimit %}Privileged user, your balance: {%= user.Finance.Balance %}.
{% else %}You don't have enough privileges.{% endif %}</p>`)
	expectCond = []byte(`<h2>Status</h2><p>Privileged user, your balance: 9000.015.</p>`)

	tplSwitch = []byte(`{%ctx exactStatus = 78 %}{
	"permission": "{% switch user.Status %}
	{% case 10 %}
		anonymous
	{% case 45 %}
		logged in
	{% case exactStatus %}
		privileged
	{% default %}
		unknown
{% endswitch %}"
}`)
	tplSwitchNoCond = []byte(`{
	"permission": "{% switch %}
	{% case user.Status <= 10 %}
		anonymous
	{% case user.Status <= 45 %}
		logged in
	{% case user.Status >= 60 %}
		privileged
	{% default %}
		unknown
{% endswitch %}"
}`)
	expectSwitch = []byte(`{"permission": "privileged"}`)

	tplLoopRange = []byte(`{
	"id":"{%= user.Id %}",
	"name":"{%= user.Name %}",
	"fin_history":[
		{% for k, item := range user.Finance.History sep , %}
		{%= k %}:{
			"utime":{%= item.DateUnix %},
			"cost":{%= item.Cost %},
			"desc":"{%= item.Comment %}"
		}
		{% endfor %}
	]
}`)
	expectLoopRange = []byte(`{"id":"115","name":"John","fin_history":[0:{"utime":152354345634,"cost":14.345241,"desc":"pay for domain"},1:{"utime":153465345246,"cost":-3.0000342543,"desc":"got refund"},2:{"utime":156436535640,"cost":2325242534.3532453,"desc":"maintenance"}]}`)

	tplLoopCountStatic = []byte(`<h2>History</h2>
<ul>
	{% for i := 0; i < 3; i++ %}
	<li>Amount: {%= user.Finance.History[i].Cost %}<br/>
		Description: {%= user.Finance.History[i].Comment %}<br/>
		Date: {%= user.Finance.History[i].DateUnix %}
	</li>
	{% endfor %}
</ul>`)
	tplLoopCount = []byte(`<h2>History</h2>
<ul>
	{% for i := begin; i < end; i++ %}
	<li>Amount: {%= user.Finance.History[i].Cost %}<br/>
		Description: {%= user.Finance.History[i].Comment %}<br/>
		Date: {%= user.Finance.History[i].DateUnix %}
	</li>
	{% endfor %}
</ul>`)
	tplLoopCountCtx = []byte(`<h2>History</h2>
{%ctx begin = 0 %}
{%ctx end = 3 %}
<ul>
	{% for i := begin; i < end; i++ %}
	<li>Amount: {%= user.Finance.History[i].Cost %}<br/>
		Description: {%= user.Finance.History[i].Comment %}<br/>
		Date: {%= user.Finance.History[i].DateUnix %}
	</li>
	{% endfor %}
</ul>`)
	expectLoopCount = []byte(`<h2>History</h2><ul><li>Amount: 14.345241<br/>Description: pay for domain<br/>Date: 152354345634</li><li>Amount: -3.0000342543<br/>Description: got refund<br/>Date: 153465345246</li><li>Amount: 2325242534.3532453<br/>Description: maintenance<br/>Date: 156436535640</li></ul>`)
)

func pretest() {
	tree, _ := Parse(tplRaw, false)
	RegisterTpl("tplRaw", tree)

	tree, _ = Parse(tplSimple, false)
	RegisterTpl("tplSimple", tree)

	tree, _ = Parse(tplCond, false)
	RegisterTpl("tplCond", tree)
	tree, _ = Parse(tplCondNoStatic, false)
	RegisterTpl("tplCondNoStatic", tree)

	tree, _ = Parse(tplSwitch, false)
	RegisterTpl("tplSwitch", tree)
	tree, _ = Parse(tplSwitchNoCond, false)
	RegisterTpl("tplSwitchNoCond", tree)

	tree, _ = Parse(tplLoopRange, false)
	RegisterTpl("tplLoopRange", tree)

	tree, _ = Parse(tplLoopCountStatic, false)
	RegisterTpl("tplLoopCountStatic", tree)
	tree, _ = Parse(tplLoopCount, false)
	RegisterTpl("tplLoopCount", tree)
	tree, _ = Parse(tplLoopCountCtx, false)
	RegisterTpl("tplLoopCountCtx", tree)
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

func TestTplCond(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render("tplCond", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectCond) {
		t.Error("cond tpl mismatch")
	}
}

func TestTplCondNoStatic(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render("tplCondNoStatic", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectCond) {
		t.Error("cond no static tpl mismatch")
	}
}

func TestTplSwitch(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render("tplSwitch", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectSwitch) {
		t.Error("switch tpl mismatch")
	}
}

func TestTplSwitchNoCond(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render("tplSwitchNoCond", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectSwitch) {
		t.Error("switch tpl mismatch")
	}
}

func TestTplLoopRange(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render("tplLoopRange", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectLoopRange) {
		t.Error("loop range tpl mismatch")
	}
}

func TestTplLoopCountStatic(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render("tplLoopCountStatic", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectLoopCount) {
		t.Error("loop count static tpl mismatch")
	}
}

func TestTplLoopCount(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	ctx.SetStatic("begin", 0)
	ctx.SetStatic("end", 3)
	result, err := Render("tplLoopCount", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectLoopCount) {
		t.Error("loop count tpl mismatch")
	}
}

func TestTplLoopCountCtx(t *testing.T) {
	pretest()

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render("tplLoopCountCtx", ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, expectLoopCount) {
		t.Error("loop count ctx tpl mismatch")
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

func BenchmarkTplCond(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		ctx.Set("user", user, &ins)
		buf.Reset()
		err := RenderTo(&buf, "tplCond", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectCond) {
			b.Error("cond tpl mismatch")
		}
		CP.Put(ctx)
	}
}

func BenchmarkTplCondNoStatic(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		ctx.Set("user", user, &ins)
		buf.Reset()
		err := RenderTo(&buf, "tplCondNoStatic", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectCond) {
			b.Error("cond no static tpl mismatch")
		}
		CP.Put(ctx)
	}
}

func BenchmarkTplSwitch(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		ctx.Set("user", user, &ins)
		buf.Reset()
		err := RenderTo(&buf, "tplSwitch", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectSwitch) {
			b.Error("switch tpl mismatch")
		}
		CP.Put(ctx)
	}
}

func BenchmarkTplSwitchNoCond(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		ctx.Set("user", user, &ins)
		buf.Reset()
		err := RenderTo(&buf, "tplSwitchNoCond", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectSwitch) {
			b.Error("switch no cond tpl mismatch")
		}
		CP.Put(ctx)
	}
}

func BenchmarkTplLoopRange(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		buf.Reset()
		ctx.Set("user", user, &ins)
		err := RenderTo(&buf, "tplLoopRange", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectLoopRange) {
			b.Error("loop range tpl mismatch")
		}
		CP.Put(ctx)
	}
}

func BenchmarkTplLoopCountStatic(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		buf.Reset()
		ctx.Set("user", user, &ins)
		err := RenderTo(&buf, "tplLoopCountStatic", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectLoopCount) {
			b.Error("loop count static tpl mismatch")
		}
		CP.Put(ctx)
	}
}

func BenchmarkTplLoopCount(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		buf.Reset()
		ctx.Set("user", user, &ins)
		ctx.SetStatic("begin", 0)
		ctx.SetStatic("end", 3)
		err := RenderTo(&buf, "tplLoopCount", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectLoopCount) {
			b.Error("loop count tpl mismatch")
		}
		CP.Put(ctx)
	}
}

func BenchmarkTplLoopCountCtx(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := CP.Get()
		buf.Reset()
		ctx.Set("user", user, &ins)
		err := RenderTo(&buf, "tplLoopCountCtx", ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), expectLoopCount) {
			b.Error("loop count ctx tpl mismatch")
		}
		CP.Put(ctx)
	}
}
