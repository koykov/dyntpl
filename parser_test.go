package dyntpl

import (
	"bytes"
	"testing"
)

var (
	cutCommentOrigin = []byte(`{# this is a test template #}
		Payload line #0
		{# some comment #}
		Payload line #1
		{# EOT #}`)
	cutCommentExpect = []byte(`Payload line #0Payload line #1`)
	cutFmtOrigin     = []byte(`
[
	{
		{%= val0 %},
		{%= val1 %}
	},
    {
        {%= val2 %},
        {%= val3 %}
    }
]
`)
	cutFmtExpect   = []byte(`[{{%= val0 %},{%= val1 %}},{{%= val2 %},{%= val3 %}}]`)
	cutFmtHrExpect = []byte(`raw: [{
tpl: val0
raw: ,
tpl: val1
raw: },{
tpl: val2
raw: ,
tpl: val3
raw: }]
`)

	uEOTOrigin = []byte(`foo {%= var1 %} bar {%= var1 end`)

	tplPS = []byte(`
{
	"key0": 15,
	"key1": {%= val0 %},
	"key2": {%= val1 suffix , %}
	{%= val2 prefix "key3" %},
	{%= val3 pfx "key4" sfx , %}
	"key5": "foo bar"
}
`)
	tplPSExpect = []byte(`raw: {"key0": 15,"key1": 
tpl: val0
raw: ,"key2": 
tpl: val1 sfx ,
tpl: val2 pfx "key3"
raw: ,
tpl: val3 pfx "key4" sfx ,
raw: "key5": "foo bar"}
`)
	tplModOrigin = []byte(`Welcome, {%= user.Name|default("anonymous") %}!`)
	tplModExpect = []byte(`raw: Welcome, 
tpl: user.Name mod default("anonymous")
raw: !
`)
	tplModNoVarOrigin = []byte(`Welcome, {%= testNameOf(user, "anonymous") %}!`)
	tplModNoVarExpect = []byte(`raw: Welcome, 
tpl:  mod testNameOf(user, "anonymous")
raw: !
`)
	tplModNestedArgOrigin = []byte(`Welcome, {%= testNameOf(user, {"foo": "bar", "id": user.Id}, "qwe") %}`)
	tplModNestedArgExpect = []byte(`raw: Welcome, 
tpl:  mod testNameOf(user, "foo":"bar", "id":user.Id, "qwe")
`)

	tplExitOrigin = []byte(`{% if user.Status == 0 %}
	{% exit %}
{% endif %}
Allowed items: ...`)
	tplExitExpect = []byte(`cond: left user.Status op == right 0
	true: 
		exit
raw: Allowed items: ...
`)

	ctxOriginAs = []byte(`{% ctx var0 = obj.Key as static %}
{% ctx var1=user.Id as static %}
{%= var1 %}`)
	ctxOriginDot = []byte(`{% ctx var0 = obj.Key.(static) %}
{% ctx var1=user.Id.(static) %}
{%= var1 %}`)
	ctxExpect = []byte(`ctx: var var0 src obj.Key ins static
ctx: var var1 src user.Id ins static
tpl: var1
`)
	ctxOriginDot1 = []byte(`{% ctx fin = obj.Finance.(TestFinance) %}
{%= fin %}`)
	ctxExpect1 = []byte(`ctx: var fin src obj.Finance ins TestFinance
tpl: fin
`)
	ctxOriginModDot = []byte(`{% ctx history = user.History|default("date").(TestHistory) %}`)
	ctxExpectModDot = []byte(`ctx: var history src user.History ins TestHistory mod default("date")
`)
	ctxOriginAsOK = []byte(`{% ctx mask, ok = user.Flags[mask] %}
{%= mask %}`)
	ctxExpectAsOK = []byte(`ctx: var mask, ok src user.Flags[mask] ins static
tpl: mask
`)
	ctxOriginAsOK1 = []byte(`{% ctx list, ok = user.History.(TestHistory) %}
{%= list %}`)
	ctxExpectAsOK1 = []byte(`ctx: var list, ok src user.History ins TestHistory
tpl: list
`)

	cntrOrigin = []byte(`{% counter i = 0 %}
{% cntr j=5 %}
foo
{% counter i++ %}
bar
{% counter j+3 %}`)
	cntrExpect = []byte(`cntr: cntr i = 0
cntr: cntr j = 5
raw: foo
cntr: cntr i op ++ arg 1
raw: bar
cntr: cntr j op ++ arg 3
`)

	condOrigin = []byte(`
<h1>Profile</h1>
{% if user.Id == 0 %}
You should <a href="{%= loginUrl %}">log in</a>.
{% else %}
Welcome, {%= user.Name %}.
{% endif %}
Copyright {%= brand %}
`)
	condExpect = []byte(`raw: <h1>Profile</h1>
cond: left user.Id op == right 0
	true: 
		raw: You should <a href="
		tpl: loginUrl
		raw: ">log in</a>.
	false: 
		raw: Welcome, 
		tpl: user.Name
		raw: .
raw: Copyright 
tpl: brand
`)
	condStrOrigin = []byte(`{% if val == "foo" %}foo{% else %}bar{% endif %}`)
	condStrExpect = []byte(`cond: left val op == right foo
	true: 
		raw: foo
	false: 
		raw: bar
`)
	condNestedOrigin = []byte(`
{% if request.Secure == 1 %}
	{% if user.AllowBuy %}
		{%= user.Name %}, you may buy our items safely.
	{% else %}
		{%= user.Name %}, you should confirm your account first.
	{% endif %}
{% endif %}
`)
	condNestedExpect = []byte(`cond: left request.Secure op == right 1
	true: 
		cond: 
			true: 
				tpl: user.Name
				raw: , you may buy our items safely.
			false: 
				tpl: user.Name
				raw: , you should confirm your account first.
`)
	condOriginOK = []byte(`{%= x %}{% if v, ok := filterVar(vars); ok %}{%= v %}{% else %}N/D{%endif%}foo`)
	condExpectOK = []byte(`tpl: x
condOK: v, ok hlp filterVar(vars); left ok
	true: 
		tpl: v
	false: 
		raw: N/D
raw: foo
`)
	condOriginNotOK = []byte(`{%= x %}{% if v, ok := filterVar(vars); !ok %}{%= v %}{% else %}N/D{%endif%}foo`)
	condExpectNotOK = []byte(`tpl: x
condOK: v, ok hlp filterVar(vars); left ok op != right true
	true: 
		tpl: v
	false: 
		raw: N/D
raw: foo
`)

	loopOrigin = []byte(`
<h2>Export history</h2>
<label>Type</label>
<select name="type">
	{% for k, v := range user.historyTags %}
	<option value="{%= k %}">{%= v %}</option>
	{% endfor %}
</select>
<label>Format</label>
<select name="fmt">
	{% for i:=0; i<4; i++ %}
	<option value="{%= i %}">{%= allowFmt[i] %}</option>
	{% endfor %}
</select>
`)
	loopExpect = []byte(`raw: <h2>Export history</h2><label>Type</label><select name="type">
rloop: key k val v src user.historyTags
	raw: <option value="
	tpl: k
	raw: ">
	tpl: v
	raw: </option>
raw: </select><label>Format</label><select name="fmt">
cloop: cnt i cond < lim 4 op ++
	raw: <option value="
	tpl: i
	raw: ">
	tpl: allowFmt[i]
	raw: </option>
raw: </select>
`)
	loopSepOrigin = []byte(`
{
	"rules": [
		{% for _, rule := range config.Rules sep , %}
		{
			"key": {%= rule.Id %},
			"val": {%= rule.Slug %}
		}
		{% endfor %}
	]
}
`)
	loopSepExpect = []byte(`raw: {"rules": [
rloop: val rule src config.Rules sep ,
	raw: {"key": 
	tpl: rule.Id
	raw: ,"val": 
	tpl: rule.Slug
	raw: }
raw: ]}
`)
	loopBrkOrigin = []byte(`
{
	{% for i:=0; i<10; i++ %}
		foo
		{% if i == 8 %}{% break %}{% endif %}
		{% if i == 7 %}{% lazybreak %}{% endif %}
		{%= i %}
	{% endfor %}
}
`)
	loopBrkExpect = []byte(`raw: {
cloop: cnt i cond < lim 10 op ++
	raw: foo
	cond: left i op == right 8
		true: 
			break
	cond: left i op == right 7
		true: 
			lazybreak
	tpl: i
raw: }
`)
	loopBrkNestedOrigin = []byte(`
{
	{% for i:=0; i<10; i++ %}
		bar
		{% for j:=0; i<10; i++ %}
			foo
			{% if j == 8 %}{% break 2 %}{% endif %}
			{% if j == 7 %}{% lazybreak 2 %}{% endif %}
			{%= j %}
		{% endfor %}
		{%= i %}
	{% endfor %}
}
`)
	loopBrkNestedExpect = []byte(`raw: {
cloop: cnt i cond < lim 10 op ++
	raw: bar
	cloop: cnt j cond < lim 10 op ++
		raw: foo
		cond: left j op == right 8
			true: 
				break 2
		cond: left j op == right 7
			true: 
				lazybreak 2
		tpl: j
	tpl: i
raw: }
`)

	switchOrigin = []byte(`<Tracking event="
	{% switch track.Event %}
	{% case param.Start %}
		start
	{% case param.FirstQuartile %}
		firstQuartile
	{% case param.Midpoint %}
		midpoint
	{% case param.ThirdQuartile %}
		thirdQuartile
	{% case param.Complete %}
		complete
	{% default %}
		unknown
	{% endswitch %}
">
	<![CDATA[{%= track.Value %}]]>
</Tracking>`)
	switchExpect = []byte(`raw: <Tracking event="
switch: arg track.Event
	case: val param.Start
		raw: start
	case: val param.FirstQuartile
		raw: firstQuartile
	case: val param.Midpoint
		raw: midpoint
	case: val param.ThirdQuartile
		raw: thirdQuartile
	case: val param.Complete
		raw: complete
	def: 
		raw: unknown
raw: "><![CDATA[
tpl: track.Value
raw: ]]></Tracking>
`)
	switchNoCondOrigin = []byte(`
[
	{
		{% switch %}
		{% case item.Index == 0 %}
			"name": {%= item.Name %},
		{% case item.Index == 1 %}
			{%= item.Slug pfx "slug": sfx , %}
		{% case item.Uid == -1 %}
			"no_data": true
		{% endswitch %}
	}
]`)
	switchNoCondExpect = []byte(`raw: [{
switch: 
	case: left item.Index op == right 0
		raw: "name": 
		tpl: item.Name
		raw: ,
	case: left item.Index op == right 1
		tpl: item.Slug pfx "slug": sfx ,
	case: left item.Uid op == right -1
		raw: "no_data": true
raw: }]
`)
	switchNoCondHelperOrigin = []byte(`
[
	{
		{% switch %}
		{% case firstItem(item) %}
			"name": {%= item.Name %},
		{% case secondItem(item, false) %}
			{%= item.Slug pfx "slug": sfx , %}
		{% case anonItem(item, 1) %}
			"no_data": true
		{% endswitch %}
	}
]`)
	switchNoCondHelperExpect = []byte(`raw: [{
switch: 
	case: firstItem(item)
		raw: "name": 
		tpl: item.Name
		raw: ,
	case: secondItem(item, "false")
		tpl: item.Slug pfx "slug": sfx ,
	case: anonItem(item, "1")
		raw: "no_data": true
raw: }]
`)
	incOrigin    = []byte(`foo {% include sidebar/right %} bar`)
	incOriginDot = []byte(`foo {% . sidebar/right %} bar`)
	incExpect    = []byte(`raw: foo 
inc: sidebar/right 
raw:  bar
`)
	locOrigin = []byte(`{% locale "ru-RU" %}{%= t("messages.welcome", "") %}`)
	locExpect = []byte(`locale: ru-RU
tpl:  mod t("messages.welcome", "")
`)
)

func TestParser(t *testing.T) {
	tst := func(t *testing.T, key string, tpl, expect []byte) {
		tree, _ := Parse(tpl, false)
		r := tree.HumanReadable()
		if !bytes.Equal(r, expect) {
			t.Errorf("%s test failed\nexp: %s\ngot: %s", key, string(expect), string(r))
		}
	}
	tstErr := func(t *testing.T, key string, tpl []byte, errExp error) {
		_, err := Parse(tpl, false)
		if err != errExp {
			t.Errorf("%s failed\nexp err: %s\ngot: %s", key, errExp, err)
		}
	}

	t.Run("cutComments", func(t *testing.T) {
		p := &Parser{tpl: cutCommentOrigin}
		p.cutComments()
		p.cutFmt()
		if !bytes.Equal(cutCommentExpect, p.tpl) {
			t.Errorf("comment cut test failed\nexp: %s\ngot: %s", string(cutCommentExpect), string(p.tpl))
		}
	})
	t.Run("cutFmt", func(t *testing.T) {
		p := &Parser{tpl: cutFmtOrigin}
		p.cutFmt()
		if !bytes.Equal(cutFmtExpect, p.tpl) {
			t.Errorf("fmt cut test failed\nexp: %s\ngot: %s", string(cutFmtExpect), string(p.tpl))
		}
	})
	t.Run("printVar", func(t *testing.T) { tst(t, "printVar", cutFmtOrigin, cutFmtHrExpect) })
	t.Run("unexpectedEOF", func(t *testing.T) { tstErr(t, "unexpectedEOF", uEOTOrigin, ErrUnexpectedEOF) })
	t.Run("prefixSuffix", func(t *testing.T) { tst(t, "prefixSuffix", tplPS, tplPSExpect) })
	t.Run("exit", func(t *testing.T) { tst(t, "exit", tplExitOrigin, tplExitExpect) })
	t.Run("mod", func(t *testing.T) { tst(t, "mod", tplModOrigin, tplModExpect) })
	t.Run("modNoVar", func(t *testing.T) { tst(t, "modNoVar", tplModNoVarOrigin, tplModNoVarExpect) })
	t.Run("modNestedArg", func(t *testing.T) { tst(t, "modNestedArg", tplModNestedArgOrigin, tplModNestedArgExpect) })

	t.Run("ctxOriginAs", func(t *testing.T) { tst(t, "ctxOriginAs", ctxOriginAs, ctxExpect) })
	t.Run("ctxDot", func(t *testing.T) { tst(t, "ctxDot", ctxOriginDot, ctxExpect) })
	t.Run("ctxDot1", func(t *testing.T) { tst(t, "ctxDot1", ctxOriginDot1, ctxExpect1) })
	t.Run("ctxModDot", func(t *testing.T) { tst(t, "ctxModDot", ctxOriginModDot, ctxExpectModDot) })
	t.Run("ctxOriginAsOK", func(t *testing.T) { tst(t, "ctxOriginAsOK", ctxOriginAsOK, ctxExpectAsOK) })
	t.Run("ctxOriginAsOK1", func(t *testing.T) { tst(t, "ctxOriginAsOK1", ctxOriginAsOK1, ctxExpectAsOK1) })

	t.Run("counter", func(t *testing.T) { tst(t, "counter", cntrOrigin, cntrExpect) })
	t.Run("condition", func(t *testing.T) { tst(t, "condition", condOrigin, condExpect) })
	t.Run("conditionStr", func(t *testing.T) { tst(t, "conditionStr", condStrOrigin, condStrExpect) })
	t.Run("conditionNested", func(t *testing.T) { tst(t, "conditionNested", condNestedOrigin, condNestedExpect) })
	t.Run("conditionOK", func(t *testing.T) { tst(t, "conditionOK", condOriginOK, condExpectOK) })
	t.Run("conditionNotOK", func(t *testing.T) { tst(t, "conditionNotOK", condOriginNotOK, condExpectNotOK) })

	t.Run("loop", func(t *testing.T) { tst(t, "loop", loopOrigin, loopExpect) })
	t.Run("loopSep", func(t *testing.T) { tst(t, "loopSep", loopSepOrigin, loopSepExpect) })
	t.Run("loopBrk", func(t *testing.T) { tst(t, "loopBrk", loopBrkOrigin, loopBrkExpect) })
	t.Run("loopBrkNested", func(t *testing.T) { tst(t, "loopBrkNested", loopBrkNestedOrigin, loopBrkNestedExpect) })

	t.Run("switch", func(t *testing.T) { tst(t, "switch", switchOrigin, switchExpect) })
	t.Run("switchNoCondition", func(t *testing.T) { tst(t, "switchNoCondition", switchNoCondOrigin, switchNoCondExpect) })
	t.Run("switchNoConditionHelper", func(t *testing.T) {
		tst(t, "switchNoConditionHelper", switchNoCondHelperOrigin, switchNoCondHelperExpect)
	})

	t.Run("include", func(t *testing.T) { tst(t, "include", incOrigin, incExpect) })
	t.Run("includeDot", func(t *testing.T) { tst(t, "includeDot", incOriginDot, incExpect) })

	t.Run("locale", func(t *testing.T) { tst(t, "locale", locOrigin, locExpect) })
}
