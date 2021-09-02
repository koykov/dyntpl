package dyntpl

import (
	"bytes"
	"testing"

	"github.com/koykov/hash/fnv"
	"github.com/koykov/i18n"
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
	tplCondStr    = []byte(`{% if user.Id == "115" %}foo{% else %}bar{% endif %}`)
	expectCondStr = []byte(`foo`)

	tplCondOK = []byte(`<ul>{% for i:=0; i<5; i++ %}
	{% if h, ok := __testUserNextHistory999(user.Finance).(TestHistory); ok %}
		<li>{%= h.Cost %}</li>
	{% endif %}
{%endfor%}</ul>`)
	expectCondOK = []byte(`<ul><li>14.345241</li><li>-3.0000342543</li><li>2325242534.3532453</li></ul>`)

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
	tplLoopRangeLBreakN = []byte(`{
	"id":"{%= user.Id %}",
	"name":"{%= user.Name %}",
	"fin_history":[
		{% for _, x := range user.Finance.History %}
			{% for k, item := range user.Finance.History sep , %}
			{%= k %}:{
				"utime":{%= item.DateUnix %},
				"cost":{%= item.Cost %},
				"desc":"{%= item.Comment %}"
				{% if k == 2 %}{% lazybreak 2 %}{% endif %}
			}
			{% endfor %}
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

	tplLoopCountBreakN = []byte(`<h2>History</h2>
{% for i := 0; i < 10; i++ %}
	<ul>
		{% for j := 0; j < 10; j++ %}
		<li>Amount: {%= user.Finance.History[j].Cost %}<br/>
			Description: {%= user.Finance.History[j].Comment %}<br/>
			Date: {%= user.Finance.History[j].DateUnix %}
		</li>
		{% if j == 2 %}{% break 2 %}{% endif %}
		{% endfor %}
	</ul>
{% endfor %}`)
	tplLoopCountLBreakN = []byte(`<h2>History</h2>
{% for i := 0; i < 10; i++ %}
	<ul>
		{% for j := 0; j < 10; j++ %}
		<li>Amount: {%= user.Finance.History[j].Cost %}<br/>
			Description: {%= user.Finance.History[j].Comment %}<br/>
			Date: {%= user.Finance.History[j].DateUnix %}
			{% if j == 2 %}{% lazybreak 2 %}{% endif %}
		</li>
		{% endfor %}
	</ul>
{% endfor %}`)
	expectLoopCountBreakN = []byte(`<h2>History</h2><ul><li>Amount: 14.345241<br/>Description: pay for domain<br/>Date: 152354345634</li><li>Amount: -3.0000342543<br/>Description: got refund<br/>Date: 153465345246</li><li>Amount: 2325242534.3532453<br/>Description: maintenance<br/>Date: 156436535640</li></ul>`)
	tplLoopCountLBreak    = []byte(`<ul>{% for i := 0; i < 10; i++ %}
	<li>
		{%= i %}: {%= user.Finance.History[i].Cost|default(0) %}
		{% if i == 2 %}{% lazybreak %}{% endif %}
	</li>
{% endfor %}</ul>`)
	expectLoopCountLBreak = []byte(`<ul><li>0: 14.345241</li><li>1: -3.0000342543</li><li>2: 2325242534.3532453</li></ul>`)

	tplCtxOK    = []byte(`foo{% ctx __testFin999, ok = user.Finance %}bar{%= __testFin999.Balance %}end`)
	expectCtxOK = []byte(`foobar9000.015end`)

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

	tplI18n             = []byte(`<h1>{%= t("messages.welcome", "", {"!user": user.Name}) %}</h1>`)
	expectI18n          = []byte(`<h1>Welcome, John!</h1>`)
	tplI18nPlural       = []byte(`<div>Multithreading support: {%= tp("pc.cpu", "N/D", cores) %}</div>`)
	expectI18nPlural    = []byte(`<div>Multithreading support: yes</div>`)
	tplI18nPluralExt    = []byte(`<div>Age: {%= tp("me.age", "unknown", years) %}</div>`)
	expectI18nPluralExt = []byte(`<div>Age: you're dead</div>`)
)

func pretest() {
	tpl := map[string][]byte{
		"tplRaw":               tplRaw,
		"tplSimple":            tplSimple,
		"tplCond":              tplCond,
		"tplCondNoStatic":      tplCondNoStatic,
		"tplCondHlp":           tplCondHlp,
		"tplCondStr":           tplCondStr,
		"tplCondOK":            tplCondOK,
		"tplSwitch":            tplSwitch,
		"tplSwitchNoCond":      tplSwitchNoCond,
		"tplLoopRange":         tplLoopRange,
		"tplLoopRangeLBreakN":  tplLoopRangeLBreakN,
		"tplLoopCountStatic":   tplLoopCountStatic,
		"tplLoopCountBreak":    tplLoopCountBreak,
		"tplLoopCountBreakN":   tplLoopCountBreakN,
		"tplLoopCountLBreak":   tplLoopCountLBreak,
		"tplLoopCountLBreakN":  tplLoopCountLBreakN,
		"tplLoopCountContinue": tplLoopCountContinue,
		"tplLoopCount":         tplLoopCount,
		"tplLoopCountCtx":      tplLoopCountCtx,
		"tplCtxOK":             tplCtxOK,
		"tplCntr0":             tplCntr0,
		"tplCntr1":             tplCntr1,
		"tplExit":              tplExit,

		"tplModDef":             tplModDef,
		"tplModDefStatic":       tplModDefStatic,
		"tplModDef1":            tplModDef1,
		"tplModJsonEscape":      tplModJsonEscape,
		"tplModJsonEscapeShort": tplModJsonEscapeShort,
		"tplModJsonEscapeDbl":   tplModJsonEscapeDbl,
		"tplModJsonQuoteShort":  tplModJsonQuoteShort,
		"tplModHtmlEscape":      tplModHtmlEscape,
		"tplModHtmlEscapeShort": tplModHtmlEscapeShort,
		"tplModLinkEscape":      tplModLinkEscape,
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

		"tplI18n":          tplI18n,
		"tplI18nPlural":    tplI18nPlural,
		"tplI18nPluralExt": tplI18nPluralExt,
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
		return
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

func TestCondOK(t *testing.T) {
	testBase(t, "tplCondOK", expectCondOK, "cond-ok tpl mismatch")
}

func TestTplCondHlp(t *testing.T) {
	testBase(t, "tplCondHlp", expectCondHlp, "cond (helper) tpl mismatch")
}

func TestTplCondStr(t *testing.T) {
	testBase(t, "tplCondStr", expectCondStr, "cond (str comparison) tpl mismatch")
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

func TestTplLoopRangeLBreakN(t *testing.T) {
	testBase(t, "tplLoopRangeLBreakN", expectLoopRange, "loop range lazybreakN tpl mismatch")
}

func TestTplLoopCountStatic(t *testing.T) {
	testBase(t, "tplLoopCountStatic", expectLoopCount, "loop count static tpl mismatch")
}

func TestTplLoopCountBreak(t *testing.T) {
	testBase(t, "tplLoopCountBreak", expectLoopCount, "loop count break tpl mismatch")
}

func TestTplLoopCountBreakN(t *testing.T) {
	testBase(t, "tplLoopCountBreakN", expectLoopCountBreakN, "loop count break N tpl mismatch")
}

func TestTplLoopCountLBreak(t *testing.T) {
	testBase(t, "tplLoopCountLBreak", expectLoopCountLBreak, "loop count lazybreak tpl mismatch")
}

func TestTplLoopCountLBreakN(t *testing.T) {
	testBase(t, "tplLoopCountLBreakN", expectLoopCountBreakN, "loop count lazybreak N tpl mismatch")
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

func TestCtxOK(t *testing.T) {
	testBase(t, "tplCtxOK", expectCtxOK, "ctx-ok tpl mismatch")
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

func TestI18n(t *testing.T) {
	pretest()
	fn := func(t *testing.T, tplName string, expect []byte, errMsg string) {
		db, _ := i18n.New(fnv.Hasher{})
		db.Set("en.messages.welcome", "Welcome, !user!")
		db.Set("en.pc.cpu", "no|yes")
		db.Set("en.me.age", "{0} you just born|[1,10] you're a child|[10,18] you're teenager|[18,40] you're adult|[40,80] you're old|[80,*] you're dead")

		ctx := NewCtx()
		ctx.I18n("en", db)
		ctx.Set("user", user, &ins)
		ctx.SetStatic("cores", 4)
		ctx.SetStatic("years", 90)
		result, err := Render(tplName, ctx)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(result, expect) {
			t.Error(errMsg)
		}
	}

	t.Run("tplI18n", func(t *testing.T) { fn(t, "tplI18n", expectI18n, "tplI18n mismatch") })
	t.Run("tplI18nPlural", func(t *testing.T) { fn(t, "tplI18nPlural", expectI18nPlural, "tplI18nPlural mismatch") })
	t.Run("tplI18nPluralExt", func(t *testing.T) { fn(t, "tplI18nPluralExt", expectI18nPluralExt, "tplI18nPluralExt mismatch") })
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

func BenchmarkTplCondStr(b *testing.B) {
	benchBase(b, "tplCondStr", expectCondStr, "cond (str comparison) tpl mismatch")
}

func BenchmarkTplCondOK(b *testing.B) {
	benchBase(b, "tplCondOK", expectCondOK, "cond-ok tpl mismatch")
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

func BenchmarkTplLoopRangeLBreakN(b *testing.B) {
	benchBase(b, "tplLoopRangeLBreakN", expectLoopRange, "loop range lazybreakN tpl mismatch")
}

func BenchmarkTplLoopCountStatic(b *testing.B) {
	benchBase(b, "tplLoopCountStatic", expectLoopCount, "loop count tpl mismatch")
}

func BenchmarkTplLoopCountBreak(b *testing.B) {
	benchBase(b, "tplLoopCountBreak", expectLoopCount, "loop count break tpl mismatch")
}

func BenchmarkTplLoopCountBreakN(b *testing.B) {
	benchBase(b, "tplLoopCountBreakN", expectLoopCountBreakN, "loop count break N tpl mismatch")
}

func BenchmarkTplLoopCountLBreak(b *testing.B) {
	benchBase(b, "tplLoopCountLBreak", expectLoopCountLBreak, "loop count lazybreak tpl mismatch")
}

func BenchmarkTplLoopCountLBreakN(b *testing.B) {
	benchBase(b, "tplLoopCountLBreakN", expectLoopCountBreakN, "loop count lazybreak N tpl mismatch")
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

func BenchmarkTplCtxOK(b *testing.B) {
	benchBase(b, "tplCtxOK", expectCtxOK, "ctx-ok tpl mismatch")
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

func BenchmarkI18n(b *testing.B) {
	pretest()
	fn := func(b *testing.B, tplName string, expect []byte, errMsg string) {
		db, _ := i18n.New(fnv.Hasher{})
		db.Set("en.messages.welcome", "Welcome, !user!")
		db.Set("en.pc.cpu", "no|yes")
		db.Set("en.me.age", "{0} you just born|[1,10] you're a child|[10,18] you're teenager|[18,40] you're adult|[40,80] you're old|[80,*] you're dead")

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			ctx := AcquireCtx()
			ctx.I18n("en", db)
			ctx.Set("user", user, &ins)
			ctx.SetStatic("cores", 4)
			ctx.SetStatic("years", 90)
			buf.Reset()
			err := RenderTo(&buf, tplName, ctx)
			if err != nil {
				b.Error(err)
			}
			if !bytes.Equal(buf.Bytes(), expect) {
				b.Error(errMsg)
			}
			ReleaseCtx(ctx)
		}
	}

	b.Run("tplI18n", func(b *testing.B) { fn(b, "tplI18n", expectI18n, "tplI18n mismatch") })
	b.Run("tplI18nPlural", func(b *testing.B) { fn(b, "tplI18nPlural", expectI18nPlural, "tplI18nPlural mismatch") })
	b.Run("tplI18nPluralExt", func(b *testing.B) { fn(b, "tplI18nPluralExt", expectI18nPluralExt, "tplI18nPluralExt mismatch") })
}
