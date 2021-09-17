package dyntpl

import (
	"bytes"
	"testing"
)

var (
	tplRaw = []byte(`<h1>Raw template<h1><p>Lorem ipsum dolor sit amet, ...</p>`)

	tplSimple = []byte(`<h1>Welcome, {%= user.Name %}!</h1>
<p>Status: {%= user.Status %}</p>
<p>Your balance: {%= user.Finance.Balance %}; buy allowance: {%= user.Finance.AllowBuy %}</p>`)
	expectSimple = []byte(`<h1>Welcome, John!</h1><p>Status: 78</p><p>Your balance: 9000.015; buy allowance: false</p>`)

	tplCondition = []byte(`<h2>Status</h2><p>
{% if user.Status >= 60 %}Privileged user, your balance: {%= user.Finance.Balance %}.
{% else %}You don't have enough privileges.{% endif %}</p>`)
	tplConditionNoStatic = []byte(`<h2>Status</h2><p>
{% ctx permissionLimit = 60 %}
{% if user.Status >= permissionLimit %}Privileged user, your balance: {%= user.Finance.Balance %}.
{% else %}You don't have enough privileges.{% endif %}</p>`)
	expectCondition    = []byte(`<h2>Status</h2><p>Privileged user, your balance: 9000.015.</p>`)
	tplConditionHlp    = []byte(`{% if lenGt0(user.Id) %}greater than zero{% endif %}`)
	expectConditionHlp = []byte(`greater than zero`)
	tplConditionStr    = []byte(`{% if user.Id == "115" %}foo{% else %}bar{% endif %}`)
	expectConditionStr = []byte(`foo`)

	tplConditionOK = []byte(`<ul>{% for i:=0; i<5; i++ %}
	{% if h, ok := __testUserNextHistory999(user.Finance).(TestHistory); ok %}
		<li>{%= h.Cost %}</li>
	{% endif %}
{%endfor%}</ul>`)
	expectConditionOK = []byte(`<ul><li>14.345241</li><li>-3.0000342543</li><li>2325242534.3532453</li></ul>`)

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
	tplSwitchNoCondition = []byte(`{
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
	tplLoopRangeLazybreakN = []byte(`{
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
	tplLoopCountLazybreakN = []byte(`<h2>History</h2>
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
	tplLoopCountLazybreak = []byte(`<ul>{% for i := 0; i < 10; i++ %}
	<li>
		{%= i %}: {%= user.Finance.History[i].Cost|default(0) %}
		{% if i == 2 %}{% lazybreak %}{% endif %}
	</li>
{% endfor %}</ul>`)
	expectLoopCountLazybreak = []byte(`<ul><li>0: 14.345241</li><li>1: -3.0000342543</li><li>2: 2325242534.3532453</li></ul>`)

	tplCtxOK    = []byte(`foo{% ctx __testFin999, ok = user.Finance %}bar{%= __testFin999.Balance %}end`)
	expectCtxOK = []byte(`foobar9000.015end`)

	tplCounter0 = []byte(`{% counter c = 0 %}
[{% counter c+2 %}{%= c %},
{% counter c+2 %}{%= c %},
{% counter c+2 %}{%= c %},
{% counter c+2 %}{%= c %},
{% counter c+2 %}{%= c %}]`)
	tplCounter1 = []byte(`{% counter c = 0 %}
[
	{% for i := 0; i < 5; i++ separator , %}
		{% counter c++ %}
		{%= c %}
	{% endfor %}
]`)
	expectCounter0 = []byte(`[2,4,6,8,10]`)
	expectCounter1 = []byte(`[1,2,3,4,5]`)

	tplExit = []byte(`{% if user.Status < 100 %}{% exit %}{% endif %}foobar`)

	tplIncludeHost     = []byte(`foo {% include sub %} bar`)
	tplIncludeSub      = []byte(`welcome {%= user.Name %}!`)
	expectTplInclude   = []byte(`foo welcome John! bar`)
	tplIncludeHostJS   = []byte(`{"a":"{% include sub1 subjs %}"}`)
	tplIncludeSubJS    = []byte(`welcome {%j= user.Id|default('anon') %}!`)
	expectTplIncludeJS = []byte(`{"a":"welcome 115!"}`)

	tplI18n             = []byte(`<h1>{%= t("messages.welcome", "", {"!user": user.Name}) %}</h1>`)
	expectI18n          = []byte(`<h1>Welcome, John!</h1>`)
	tplI18nPlural       = []byte(`<div>Multithreading support: {%= tp("pc.cpu", "N/D", cores) %}</div>`)
	expectI18nPlural    = []byte(`<div>Multithreading support: yes</div>`)
	tplI18nPluralExt    = []byte(`<div>Age: {%= tp("me.age", "unknown", years) %}</div>`)
	expectI18nPluralExt = []byte(`<div>Age: you're dead</div>`)
	tplI18nSetLocale    = []byte(`{% locale "ru" %}<h1>{%= t("messages.welcome", "", {"!user": user.Name}) %}</h1>`)
	expectI18nSetLocale = []byte(`<h1>Привет, John!</h1>`)
)

func pretest() {
	tpl := map[string][]byte{
		"tplRaw":                 tplRaw,
		"tplSimple":              tplSimple,
		"tplCondition":           tplCondition,
		"tplConditionNoStatic":   tplConditionNoStatic,
		"tplConditionHlp":        tplConditionHlp,
		"tplConditionStr":        tplConditionStr,
		"tplConditionOK":         tplConditionOK,
		"tplSwitch":              tplSwitch,
		"tplSwitchNoCondition":   tplSwitchNoCondition,
		"tplLoopRange":           tplLoopRange,
		"tplLoopRangeLazybreakN": tplLoopRangeLazybreakN,
		"tplLoopCountStatic":     tplLoopCountStatic,
		"tplLoopCountBreak":      tplLoopCountBreak,
		"tplLoopCountBreakN":     tplLoopCountBreakN,
		"tplLoopCountLazybreak":  tplLoopCountLazybreak,
		"tplLoopCountLazybreakN": tplLoopCountLazybreakN,
		"tplLoopCountContinue":   tplLoopCountContinue,
		"tplLoopCount":           tplLoopCount,
		"tplLoopCountCtx":        tplLoopCountCtx,
		"tplCtxOK":               tplCtxOK,
		"tplCounter0":            tplCounter0,
		"tplCounter1":            tplCounter1,
		"tplExit":                tplExit,

		"tplModDefault":         tplModDefault,
		"tplModDefaultStatic":   tplModDefaultStatic,
		"tplModDefault1":        tplModDefault1,
		"tplModJSONEscape":      tplModJSONEscape,
		"tplModJSONEscapeShort": tplModJSONEscapeShort,
		"tplModJSONEscapeDbl":   tplModJSONEscapeDbl,
		"tplModJSONQuoteShort":  tplModJSONQuoteShort,
		"tplModHtmlEscape":      tplModHtmlEscape,
		"tplModHtmlEscapeShort": tplModHtmlEscapeShort,
		"tplModLinkEscape":      tplModLinkEscape,
		"tplModURLEncode":       tplModURLEncode,
		"tplModURLEncode2":      tplModURLEncode2,
		"tplModURLEncode3":      tplModURLEncode3,
		"tplModIfThen":          tplModIfThen,
		"tplModIfThenElse":      tplModIfThenElse,
		"tplModRound":           tplModRound,

		"tplIncludeHost":   tplIncludeHost,
		"sub":              tplIncludeSub,
		"tplIncludeHostJS": tplIncludeHostJS,
		"subjs":            tplIncludeSubJS,

		"tplI18n":          tplI18n,
		"tplI18nPlural":    tplI18nPlural,
		"tplI18nPluralExt": tplI18nPluralExt,
		"tplI18nSetLocale": tplI18nSetLocale,
	}
	for name, body := range tpl {
		tree, _ := Parse(body, false)
		RegisterTplKey(name, tree)
	}
}

type tplStage struct {
	key string
	fn  func(t *testing.T, key string)
}

func TestTpl(t *testing.T) {
	loadStages()

	tplStages := []tplStage{
		{key: "condition"},
		{key: "conditionHlp"},
		{key: "conditionNoStatic"},
		{key: "conditionOK"},
		{key: "conditionStr"},
		{key: "counter0"},
		{key: "counter1"},
		{key: "ctxOK"},
		{key: "exit"},
		{key: "includeHost"},
		{key: "includeHostJS"},
		{key: "loopCount", fn: fnTplLC},
		{key: "loopCountBreak"},
		{key: "loopCountBreakN"},
		{key: "loopCountContinue"},
		{key: "loopCountCtx"},
		{key: "loopCountLazybreak"},
		{key: "loopCountLazybreakN"},
		{key: "loopCountStatic"},
		{key: "loopRange"},
		{key: "loopRangeLazybreakN"},
		{key: "raw"},
		{key: "simple"},
		{key: "switch"},
		{key: "switchNoCondition"},
	}

	for _, s := range tplStages {
		t.Run(s.key, func(t *testing.T) {
			if s.fn == nil {
				s.fn = fnTpl
			}
			s.fn(t, s.key)
		})
	}
}

func fnTpl(t *testing.T, key string) {
	st := getStage(key)
	if st == nil {
		t.Error("stage not found")
		return
	}

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render(key, ctx)
	if err != nil {
		t.Error(err)
	}
	if len(st.expect) == 0 && len(result) != 0 {
		t.Errorf("%s mismatch", key)
		return
	}
	if !bytes.Equal(result, st.expect) {
		t.Errorf("%s mismatch", key)
	}
}

func fnTplLC(t *testing.T, key string) {
	st := getStage(key)
	if st == nil {
		t.Error("stage not found")
		return
	}

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	ctx.SetStatic("begin", 0)
	ctx.SetStatic("end", 3)
	result, err := Render(key, ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, st.expect) {
		t.Error("loop count tpl mismatch")
	}
}

func benchBase(b *testing.B, tplName string, expectResult []byte, errMsg string) {
	pretest()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		ctx.Set("user", user, &ins)
		buf.Reset()
		err := Write(&buf, tplName, ctx)
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

func BenchmarkTplSimple(b *testing.B) {
	benchBase(b, "tplSimple", expectSimple, "simple tpl mismatch")
}

func BenchmarkTplCondition(b *testing.B) {
	benchBase(b, "tplCondition", expectCondition, "cond tpl mismatch")
}

func BenchmarkTplConditionNoStatic(b *testing.B) {
	benchBase(b, "tplConditionNoStatic", expectCondition, "cond no static tpl mismatch")
}

func BenchmarkTplConditionHlp(b *testing.B) {
	benchBase(b, "tplConditionHlp", expectConditionHlp, "cond (helper) tpl mismatch")
}

func BenchmarkTplConditionStr(b *testing.B) {
	benchBase(b, "tplConditionStr", expectConditionStr, "cond (str comparison) tpl mismatch")
}

func BenchmarkTplConditionOK(b *testing.B) {
	benchBase(b, "tplConditionOK", expectConditionOK, "cond-ok tpl mismatch")
}

func BenchmarkTplSwitch(b *testing.B) {
	benchBase(b, "tplSwitch", expectSwitch, "switch tpl mismatch")
}

func BenchmarkTplSwitchNoCondition(b *testing.B) {
	benchBase(b, "tplSwitchNoCondition", expectSwitch, "switch no cond tpl mismatch")
}

func BenchmarkTplLoopRange(b *testing.B) {
	benchBase(b, "tplLoopRange", expectLoopRange, "loop range tpl mismatch")
}

func BenchmarkTplLoopRangeLazybreakN(b *testing.B) {
	benchBase(b, "tplLoopRangeLazybreakN", expectLoopRange, "loop range lazybreakN tpl mismatch")
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

func BenchmarkTplLoopCountLazybreak(b *testing.B) {
	benchBase(b, "tplLoopCountLazybreak", expectLoopCountLazybreak, "loop count lazybreak tpl mismatch")
}

func BenchmarkTplLoopCountLazybreakN(b *testing.B) {
	benchBase(b, "tplLoopCountLazybreakN", expectLoopCountBreakN, "loop count lazybreak N tpl mismatch")
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
		err := Write(&buf, "tplLoopCount", ctx)
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

func BenchmarkTplCounter0(b *testing.B) {
	benchBase(b, "tplCounter0", expectCounter0, "cntr 0 tpl mismatch")
}

func BenchmarkTplCounter1(b *testing.B) {
	benchBase(b, "tplCounter1", expectCounter1, "cntr 1 tpl mismatch")
}

func BenchmarkTplExit(b *testing.B) {
	benchBase(b, "tplExit", nil, "exit tpl mismatch")
}

func BenchmarkTplIncludelude(b *testing.B) {
	benchBase(b, "tplIncludeHost", expectTplInclude, "include tpl mismatch")
}

func BenchmarkTplIncludeludeJS(b *testing.B) {
	benchBase(b, "tplIncludeHostJS", expectTplIncludeJS, "include tpl (js) mismatch")
}
