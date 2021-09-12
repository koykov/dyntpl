package dyntpl

import (
	"bytes"
	"testing"

	"github.com/koykov/bytealg"
)

func TestParser(t *testing.T) {
	type stage struct {
		key string
		fn  func(*testing.T, *stage)
		origin,
		expect []byte
		err error
	}

	tst := func(t *testing.T, stage *stage) {
		tree, _ := Parse(stage.origin, false)
		r := tree.HumanReadable()
		r = bytealg.Trim(r, []byte("\n"))
		if !bytes.Equal(r, stage.expect) {
			t.Errorf("%s test failed\nexp: %s\ngot: %s", stage.key, string(stage.expect), string(r))
		}
	}
	tstRaw := func(t *testing.T, stage *stage) {
		p := &Parser{tpl: stage.origin}
		p.cutComments()
		p.cutFmt()
		if !bytes.Equal(stage.expect, p.tpl) {
			t.Errorf("%s test raw failed\nexp: %s\ngot: %s", stage.key, string(stage.expect), string(p.tpl))
		}
	}
	tstErr := func(t *testing.T, stage *stage) {
		_, err := Parse(stage.origin, false)
		if err != stage.err {
			t.Errorf("%s test error failed\nexp err: %s\ngot: %s", stage.key, stage.err, err)
		}
	}

	var (
		stages = []stage{
			{
				key: "cutComments",
				fn:  tstRaw,
				origin: []byte(`{# this is a test template #}
Payload line #0
{# some comment #}
Payload line #1
{# EOT #}`),
				expect: []byte(`Payload line #0Payload line #1`),
			},
			{
				key: "cutFmt",
				fn:  tstRaw,
				origin: []byte(`
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
`),
				expect: []byte(`[{{%= val0 %},{%= val1 %}},{{%= val2 %},{%= val3 %}}]`),
			},
			{
				key: "printVar",
				fn:  tst,
				origin: []byte(`
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
`),
				expect: []byte(`raw: [{
tpl: val0
raw: ,
tpl: val1
raw: },{
tpl: val2
raw: ,
tpl: val3
raw: }]`),
			},
			{
				key:    "unexpectedEOF",
				fn:     tstErr,
				origin: []byte(`foo {%= var1 %} bar {%= var1 end`),
				err:    ErrUnexpectedEOF,
			},
			{
				key: "prefixSuffix",
				fn:  tst,
				origin: []byte(`
{
	"key0": 15,
	"key1": {%= val0 %},
	"key2": {%= val1 suffix , %}
	{%= val2 prefix "key3" %},
	{%= val3 pfx "key4" sfx , %}
	"key5": "foo bar"
}
`),
				expect: []byte(`raw: {"key0": 15,"key1": 
tpl: val0
raw: ,"key2": 
tpl: val1 sfx ,
tpl: val2 pfx "key3"
raw: ,
tpl: val3 pfx "key4" sfx ,
raw: "key5": "foo bar"}`),
			},
			{
				key: "exit",
				fn:  tst,
				origin: []byte(`{% if user.Status == 0 %}
	{% exit %}
{% endif %}
Allowed items: ...`),
				expect: []byte(`cond: left user.Status op == right 0
	true: 
		exit
raw: Allowed items: ...`),
			},
			{
				key:    "mod",
				fn:     tst,
				origin: []byte(`Welcome, {%= user.Name|default("anonymous") %}!`),
				expect: []byte(`raw: Welcome, 
tpl: user.Name mod default("anonymous")
raw: !`),
			},
			{
				key:    "modNoVar",
				fn:     tst,
				origin: []byte(`Welcome, {%= testNameOf(user, "anonymous") %}!`),
				expect: []byte(`raw: Welcome, 
tpl:  mod testNameOf(user, "anonymous")
raw: !`),
			},
			{
				key:    "modNestedArg",
				fn:     tst,
				origin: []byte(`Welcome, {%= testNameOf(user, {"foo": "bar", "id": user.Id}, "qwe") %}`),
				expect: []byte(`raw: Welcome, 
tpl:  mod testNameOf(user, "foo":"bar", "id":user.Id, "qwe")`),
			},
			{
				key: "ctxAs",
				fn:  tst,
				origin: []byte(`{% ctx var0 = obj.Key as static %}
{% ctx var1=user.Id as static %}
{%= var1 %}`),
				expect: []byte(`ctx: var var0 src obj.Key ins static
ctx: var var1 src user.Id ins static
tpl: var1`),
			},
			{
				key: "ctxDot",
				fn:  tst,
				origin: []byte(`{% ctx var0 = obj.Key.(static) %}
{% ctx var1=user.Id.(static) %}
{%= var1 %}`),
				expect: []byte(`ctx: var var0 src obj.Key ins static
ctx: var var1 src user.Id ins static
tpl: var1`),
			},
			{
				key: "ctxDot1",
				fn:  tst,
				origin: []byte(`{% ctx fin = obj.Finance.(TestFinance) %}
{%= fin %}`),
				expect: []byte(`ctx: var fin src obj.Finance ins TestFinance
tpl: fin`),
			},
			{
				key:    "ctxModDot",
				fn:     tst,
				origin: []byte(`{% ctx history = user.History|default("date").(TestHistory) %}`),
				expect: []byte(`ctx: var history src user.History ins TestHistory mod default("date")`),
			},
			{
				key: "ctxAsOK",
				fn:  tst,
				origin: []byte(`{% ctx mask, ok = user.Flags[mask] %}
{%= mask %}`),
				expect: []byte(`ctx: var mask, ok src user.Flags[mask] ins static
tpl: mask`),
			},
			{
				key: "ctx",
				fn:  tst,
				origin: []byte(`{% ctx list, ok = user.History.(TestHistory) %}
{%= list %}`),
				expect: []byte(`ctx: var list, ok src user.History ins TestHistory
tpl: list`),
			},
			{
				key: "counter",
				fn:  tst,
				origin: []byte(`{% counter i = 0 %}
{% cntr j=5 %}
foo
{% counter i++ %}
bar
{% counter j+3 %}`),
				expect: []byte(`cntr: cntr i = 0
cntr: cntr j = 5
raw: foo
cntr: cntr i op ++ arg 1
raw: bar
cntr: cntr j op ++ arg 3`),
			},
			{
				key: "condition",
				fn:  tst,
				origin: []byte(`
<h1>Profile</h1>
{% if user.Id == 0 %}
You should <a href="{%= loginUrl %}">log in</a>.
{% else %}
Welcome, {%= user.Name %}.
{% endif %}
Copyright {%= brand %}
`),
				expect: []byte(`raw: <h1>Profile</h1>
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
tpl: brand`),
			},
			{
				key:    "conditionStr",
				fn:     tst,
				origin: []byte(`{% if val == "foo" %}foo{% else %}bar{% endif %}`),
				expect: []byte(`cond: left val op == right foo
	true: 
		raw: foo
	false: 
		raw: bar`),
			},
			{
				key: "conditionNested",
				fn:  tst,
				origin: []byte(`
{% if request.Secure == 1 %}
	{% if user.AllowBuy %}
		{%= user.Name %}, you may buy our items safely.
	{% else %}
		{%= user.Name %}, you should confirm your account first.
	{% endif %}
{% endif %}`),
				expect: []byte(`cond: left request.Secure op == right 1
	true: 
		cond: 
			true: 
				tpl: user.Name
				raw: , you may buy our items safely.
			false: 
				tpl: user.Name
				raw: , you should confirm your account first.`),
			},
			{
				key:    "conditionOK",
				fn:     tst,
				origin: []byte(`{%= x %}{% if v, ok := filterVar(vars); ok %}{%= v %}{% else %}N/D{%endif%}foo`),
				expect: []byte(`tpl: x
condOK: v, ok hlp filterVar(vars); left ok
	true: 
		tpl: v
	false: 
		raw: N/D
raw: foo`),
			},
			{
				key:    "conditionNotOK",
				fn:     tst,
				origin: []byte(`{%= x %}{% if v, ok := filterVar(vars); !ok %}{%= v %}{% else %}N/D{%endif%}foo`),
				expect: []byte(`tpl: x
condOK: v, ok hlp filterVar(vars); left ok op != right true
	true: 
		tpl: v
	false: 
		raw: N/D
raw: foo`),
			},
			{
				key: "loop",
				fn:  tst,
				origin: []byte(`
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
</select>`),
				expect: []byte(`raw: <h2>Export history</h2><label>Type</label><select name="type">
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
raw: </select>`),
			},
			{
				key: "loopSeparator",
				fn:  tst,
				origin: []byte(`
{
	"rules": [
		{% for _, rule := range config.Rules sep , %}
		{
			"key": {%= rule.Id %},
			"val": {%= rule.Slug %}
		}
		{% endfor %}
	]
}`),
				expect: []byte(`raw: {"rules": [
rloop: val rule src config.Rules sep ,
	raw: {"key": 
	tpl: rule.Id
	raw: ,"val": 
	tpl: rule.Slug
	raw: }
raw: ]}`),
			},
			{
				key: "loopBreak",
				fn:  tst,
				origin: []byte(`
{
	{% for i:=0; i<10; i++ %}
		foo
		{% if i == 8 %}{% break %}{% endif %}
		{% if i == 7 %}{% lazybreak %}{% endif %}
		{%= i %}
	{% endfor %}
}`),
				expect: []byte(`raw: {
cloop: cnt i cond < lim 10 op ++
	raw: foo
	cond: left i op == right 8
		true: 
			break
	cond: left i op == right 7
		true: 
			lazybreak
	tpl: i
raw: }`),
			},
			{
				key: "loopBreakNested",
				fn:  tst,
				origin: []byte(`
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
}`),
				expect: []byte(`raw: {
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
raw: }`),
			},
			{
				key: "switch",
				fn:  tst,
				origin: []byte(`<Tracking event="
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
</Tracking>`),
				expect: []byte(`raw: <Tracking event="
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
raw: ]]></Tracking>`),
			},
			{
				key: "switchNoCondition",
				fn:  tst,
				origin: []byte(`
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
]`),
				expect: []byte(`raw: [{
switch: 
	case: left item.Index op == right 0
		raw: "name": 
		tpl: item.Name
		raw: ,
	case: left item.Index op == right 1
		tpl: item.Slug pfx "slug": sfx ,
	case: left item.Uid op == right -1
		raw: "no_data": true
raw: }]`),
			},
			{
				key: "switchNoConditionWithHelper",
				fn:  tst,
				origin: []byte(`
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
]`),
				expect: []byte(`raw: [{
switch: 
	case: firstItem(item)
		raw: "name": 
		tpl: item.Name
		raw: ,
	case: secondItem(item, "false")
		tpl: item.Slug pfx "slug": sfx ,
	case: anonItem(item, "1")
		raw: "no_data": true
raw: }]`),
			},
			{
				key:    "include",
				fn:     tst,
				origin: []byte(`foo {% include sidebar/right %} bar`),
				expect: []byte(`raw: foo 
inc: sidebar/right 
raw:  bar`),
			},
			{
				key:    "includeDot",
				fn:     tst,
				origin: []byte(`foo {% . sidebar/right %} bar`),
				expect: []byte(`raw: foo 
inc: sidebar/right 
raw:  bar`),
			},
			{
				key:    "locale",
				fn:     tst,
				origin: []byte(`{% locale "ru-RU" %}{%= t("messages.welcome", "") %}`),
				expect: []byte(`locale: ru-RU
tpl:  mod t("messages.welcome", "")`),
			},
		}
	)

	for _, s := range stages {
		t.Run(s.key, func(t *testing.T) { s.fn(t, &s) })
	}
}
