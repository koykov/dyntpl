package dyntpl

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
{% ctx permissionLimit = 60 %}
{% if user.Status >= permissionLimit %}Privileged user, your balance: {%= user.Finance.Balance %}.
{% else %}You don't have enough privileges.{% endif %}</p>`)
	expectCond    = []byte(`<h2>Status</h2><p>Privileged user, your balance: 9000.015.</p>`)
	tplCondHlp    = []byte(`{% if lenGt0(user.Id) %}greater than zero{% endif %}`)
	expectCondHlp = []byte(`greater than zero`)

	tplSwitch = []byte(`{% ctx exactStatus = 78 %}{
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
{% ctx begin = 0 %}
{% ctx end = 3 %}
<ul>
	{% for i := begin; i < end; i++ %}
	<li>Amount: {%= user.Finance.History[i].Cost %}<br/>
		Description: {%= user.Finance.History[i].Comment %}<br/>
		Date: {%= user.Finance.History[i].DateUnix %}
	</li>
	{% endfor %}
</ul>`)
	tplLoopCountBreak = []byte(`<h2>History</h2>
<ul>
	{% for i := 0; i < 10; i++ %}
	<li>Amount: {%= user.Finance.History[i].Cost %}<br/>
		Description: {%= user.Finance.History[i].Comment %}<br/>
		Date: {%= user.Finance.History[i].DateUnix %}
	</li>
	{% if i == 2 %}{% break %}{% endif %}
	{% endfor %}
</ul>`)
	tplLoopCountContinue = []byte(`<h2>History</h2>
<ul>
	{% for i := 0; i < 10; i++ %}
	{% if i > 2 %}{% continue %}{% endif %}
	<li>Amount: {%= user.Finance.History[i].Cost %}<br/>
		Description: {%= user.Finance.History[i].Comment %}<br/>
		Date: {%= user.Finance.History[i].DateUnix %}
	</li>
	{% endfor %}
</ul>`)
	expectLoopCount = []byte(`<h2>History</h2><ul><li>Amount: 14.345241<br/>Description: pay for domain<br/>Date: 152354345634</li><li>Amount: -3.0000342543<br/>Description: got refund<br/>Date: 153465345246</li><li>Amount: 2325242534.3532453<br/>Description: maintenance<br/>Date: 156436535640</li></ul>`)

	tplCntr0 = []byte(`{% counter c = 0 %}
[{% counter c+2 %}{%= c %},
{% counter c+2 %}{%= c %},
{% counter c+2 %}{%= c %},
{% counter c+2 %}{%= c %},
{% counter c+2 %}{%= c %}]`)
	tplCntr1 = []byte(`{% counter c = 0 %}
[
	{% for i := 0; i < 5; i++ separator , %}
		{% counter c++ %}
		{%= c %}
	{% endfor %}
]`)
	expectCntr0 = []byte(`[2,4,6,8,10]`)
	expectCntr1 = []byte(`[1,2,3,4,5]`)

	tplExit = []byte(`{% if user.Status < 100 %}{% exit %}{% endif %}foobar`)

	tplIncHost     = []byte(`foo {% include sub %} bar`)
	tplIncSub      = []byte(`welcome {%= user.Name %}!`)
	expectTplInc   = []byte(`foo welcome John! bar`)
	tplIncHostJS   = []byte(`{"a":"{% include sub1 subjs %}"}`)
	tplIncSubJS    = []byte(`welcome {%j= user.Id|default('anon') %}!`)
	expectTplIncJS = []byte(`{"a":"welcome 115!"}`)
)

func pretest() {
	tpl := map[string][]byte{
		"tplRaw":               tplRaw,
		"tplSimple":            tplSimple,
		"tplCond":              tplCond,
		"tplCondNoStatic":      tplCondNoStatic,
		"tplCondHlp":           tplCondHlp,
		"tplSwitch":            tplSwitch,
		"tplSwitchNoCond":      tplSwitchNoCond,
		"tplLoopRange":         tplLoopRange,
		"tplLoopCountStatic":   tplLoopCountStatic,
		"tplLoopCountBreak":    tplLoopCountBreak,
		"tplLoopCountContinue": tplLoopCountContinue,
		"tplLoopCount":         tplLoopCount,
		"tplLoopCountCtx":      tplLoopCountCtx,
		"tplCntr0":             tplCntr0,
		"tplCntr1":             tplCntr1,
		"tplExit":              tplExit,

		"tplModDef":             tplModDef,
		"tplModDefStatic":       tplModDefStatic,
		"tplModJsonEscape":      tplModJsonEscape,
		"tplModJsonEscapeShort": tplModJsonEscapeShort,
		"tplModJsonEscapeDbl":   tplModJsonEscapeDbl,
		"tplModJsonQuoteShort":  tplModJsonQuoteShort,
		"tplModHtmlEscape":      tplModHtmlEscape,
		"tplModHtmlEscapeShort": tplModHtmlEscapeShort,
		"tplModUrlEncode":       tplModUrlEncode,
		"tplModUrlEncode2":      tplModUrlEncode2,
		"tplModUrlEncode3":      tplModUrlEncode3,
		"tplModIfThen":          tplModIfThen,
		"tplModIfThenElse":      tplModIfThenElse,
		"tplModRound":           tplModRound,

		"tplIncHost":   tplIncHost,
		"sub":          tplIncSub,
		"tplIncHostJS": tplIncHostJS,
		"subjs":        tplIncSubJS,
	}
	for name, body := range tpl {
		tree, _ := Parse(body, false)
		RegisterTpl(name, tree)
	}
}

func testBase(t *testing.T, tplName string, expectResult []byte, errMsg string) {
	pretest()

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render(tplName, ctx)
	if err != nil {
		t.Error(err)
	}
	if len(expectResult) == 0 && len(result) != 0 {
		t.Error(errMsg)
	}
	if !bytes.Equal(result, expectResult) {
		t.Error(errMsg)
	}
}

func benchBase(b *testing.B, tplName string, expectResult []byte, errMsg string) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		ctx.Set("user", user, &ins)
		buf.Reset()
		err := RenderTo(&buf, tplName, ctx)
		if err != nil {
			b.Error(err)
		}
		if len(expectResult) == 0 && buf.Len() != 0 {
			b.Error(errMsg)
		}
		if !bytes.Equal(buf.Bytes(), expectResult) {
			b.Error(errMsg)
		}
		ReleaseCtx(ctx)
	}
}

func TestTplRaw(t *testing.T) {
	testBase(t, "tplRaw", tplRaw, "raw tpl mismatch")
}

func TestTplSimple(t *testing.T) {
	testBase(t, "tplSimple", expectSimple, "simple tpl mismatch")
}

func TestTplCond(t *testing.T) {
	testBase(t, "tplCond", expectCond, "cond tpl mismatch")
}

func TestTplCondNoStatic(t *testing.T) {
	testBase(t, "tplCondNoStatic", expectCond, "cond (no static) tpl mismatch")
}

func TestTplCondHlp(t *testing.T) {
	testBase(t, "tplCondHlp", expectCondHlp, "cond (helper) tpl mismatch")
}

func TestTplSwitch(t *testing.T) {
	testBase(t, "tplSwitch", expectSwitch, "switch tpl mismatch")
}

func TestTplSwitchNoCond(t *testing.T) {
	testBase(t, "tplSwitchNoCond", expectSwitch, "switch (no cond) tpl mismatch")
}

func TestTplLoopRange(t *testing.T) {
	testBase(t, "tplLoopRange", expectLoopRange, "loop range tpl mismatch")
}

func TestTplLoopCountStatic(t *testing.T) {
	testBase(t, "tplLoopCountStatic", expectLoopCount, "loop count static tpl mismatch")
}

func TestTplLoopCountBreak(t *testing.T) {
	testBase(t, "tplLoopCountBreak", expectLoopCount, "loop count break tpl mismatch")
}

func TestTplLoopCountContinue(t *testing.T) {
	testBase(t, "tplLoopCountContinue", expectLoopCount, "loop count continue tpl mismatch")
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
	testBase(t, "tplLoopCountCtx", expectLoopCount, "loop count ctx tpl mismatch")
}

func TestCntr0Exit(t *testing.T) {
	testBase(t, "tplCntr0", expectCntr0, "cntr 0 tpl mismatch")
}

func TestCntr1Exit(t *testing.T) {
	testBase(t, "tplCntr1", expectCntr1, "cntr 1 tpl mismatch")
}

func TestTplExit(t *testing.T) {
	testBase(t, "tplExit", nil, "exit tpl mismatch")
}

func TestTplInclude(t *testing.T) {
	testBase(t, "tplIncHost", expectTplInc, "include tpl mismatch")
}

func TestTplIncludeJS(t *testing.T) {
	testBase(t, "tplIncHostJS", expectTplIncJS, "include tpl (js) mismatch")
}

func BenchmarkTplSimple(b *testing.B) {
	benchBase(b, "tplSimple", expectSimple, "simple tpl mismatch")
}

func BenchmarkTplCond(b *testing.B) {
	benchBase(b, "tplCond", expectCond, "cond tpl mismatch")
}

func BenchmarkTplCondNoStatic(b *testing.B) {
	benchBase(b, "tplCondNoStatic", expectCond, "cond no static tpl mismatch")
}

func BenchmarkTplCondHlp(b *testing.B) {
	benchBase(b, "tplCondHlp", expectCondHlp, "cond (helper) tpl mismatch")
}

func BenchmarkTplSwitch(b *testing.B) {
	benchBase(b, "tplSwitch", expectSwitch, "switch tpl mismatch")
}

func BenchmarkTplSwitchNoCond(b *testing.B) {
	benchBase(b, "tplSwitchNoCond", expectSwitch, "switch no cond tpl mismatch")
}

func BenchmarkTplLoopRange(b *testing.B) {
	benchBase(b, "tplLoopRange", expectLoopRange, "loop range tpl mismatch")
}

func BenchmarkTplLoopCountStatic(b *testing.B) {
	benchBase(b, "tplLoopCountStatic", expectLoopCount, "loop count tpl mismatch")
}

func BenchmarkTplLoopCountBreak(b *testing.B) {
	benchBase(b, "tplLoopCountBreak", expectLoopCount, "loop count break tpl mismatch")
}

func BenchmarkTplLoopCountContinue(b *testing.B) {
	benchBase(b, "tplLoopCountContinue", expectLoopCount, "loop count continue tpl mismatch")
}

func BenchmarkTplLoopCount(b *testing.B) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
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
		ReleaseCtx(ctx)
	}
}

func BenchmarkTplLoopCountCtx(b *testing.B) {
	benchBase(b, "tplLoopCountCtx", expectLoopCount, "loop count ctx tpl mismatch")
}

func BenchmarkTplCntr0(b *testing.B) {
	benchBase(b, "tplCntr0", expectCntr0, "cntr 0 tpl mismatch")
}

func BenchmarkTplCntr1(b *testing.B) {
	benchBase(b, "tplCntr1", expectCntr1, "cntr 1 tpl mismatch")
}

func BenchmarkTplExit(b *testing.B) {
	benchBase(b, "tplExit", nil, "exit tpl mismatch")
}

func BenchmarkTplInclude(b *testing.B) {
	benchBase(b, "tplIncHost", expectTplInc, "include tpl mismatch")
}

func BenchmarkTplIncludeJS(b *testing.B) {
	benchBase(b, "tplIncHostJS", expectTplIncJS, "include tpl (js) mismatch")
}
