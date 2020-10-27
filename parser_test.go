package dyntpl

import (
	"bytes"
	"testing"
)

var (
	cutFmtOrigin = []byte(`
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
	cutFmtExpect = []byte(`[{{%= val0 %},{%= val1 %}},{{%= val2 %},{%= val3 %}}]`)

	primHrExpect = []byte(`raw: [{
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

	tplExitOrigin = []byte(`{% if user.Status == 0 %}
	{% exit %}
{% endif %}
Allowed items: ...`)
	tplExitExpect = []byte(`cond: left user.Status op == right 0
	true: 
		exit
raw: Allowed items: ...
`)

	ctxOrigin = []byte(`{% ctx var0 = obj.Key as static %}
{% ctx var1=user.Id as static %}
{%= var1 %}`)
	ctxExpect = []byte(`ctx: var var0 src obj.Key ins static
ctx: var var1 src user.Id ins static
tpl: var1
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
)

func TestParseCutComments(t *testing.T) {
	tpl := []byte(`{# this is a test template #}
		Payload line #0
		{# some comment #}
		Payload line #1
		{# EOT #}`)
	exp := []byte(`Payload line #0Payload line #1`)
	p := &Parser{tpl: tpl}
	p.cutComments()
	p.cutFmt()
	if !bytes.Equal(exp, p.tpl) {
		t.Errorf("comment cut test failed\nexp: %s\ngot: %s", string(exp), string(p.tpl))
	}
}

func TestParseCutFmt(t *testing.T) {
	p := &Parser{tpl: cutFmtOrigin}
	p.cutFmt()
	if !bytes.Equal(cutFmtExpect, p.tpl) {
		t.Errorf("fmt cut test failed\nexp: %s\ngot: %s", string(cutFmtExpect), string(p.tpl))
	}
}

func TestParsePrimitive(t *testing.T) {
	tree, _ := Parse(cutFmtOrigin, false)
	r := tree.HumanReadable()
	if !bytes.Equal(r, primHrExpect) {
		t.Errorf("prim test failed\nexp: %s\ngot: %s", string(primHrExpect), string(r))
	}
}

func TestParseuEOT(t *testing.T) {
	_, err := Parse(uEOTOrigin, false)
	if err != ErrUnexpectedEOF {
		t.Errorf("unexpected EOT fail\nexp: %s\ngot: %s", ErrUnexpectedEOF, err)
	}
}

func TestParsePrefixSuffix(t *testing.T) {
	tree, _ := Parse(tplPS, false)
	r := tree.HumanReadable()
	if !bytes.Equal(r, tplPSExpect) {
		t.Errorf("prefix/suffix test failed\nexp: %s\ngot: %s", string(tplPSExpect), string(r))
	}
}

func TestParseExit(t *testing.T) {
	tree, _ := Parse(tplExitOrigin, false)
	r := tree.HumanReadable()
	if !bytes.Equal(r, tplExitExpect) {
		t.Errorf("exit test failed\nexp: %s\ngot: %s", string(tplExitExpect), string(r))
	}
}

func TestParseMod(t *testing.T) {
	tree, _ := Parse(tplModOrigin, false)
	r := tree.HumanReadable()
	if !bytes.Equal(r, tplModExpect) {
		t.Errorf("mod test failed\nexp: %s\ngot: %s", string(tplModExpect), string(r))
	}
}

func TestParseModNoVar(t *testing.T) {
	tree, _ := Parse(tplModNoVarOrigin, false)
	r := tree.HumanReadable()
	if !bytes.Equal(r, tplModNoVarExpect) {
		t.Errorf("mod (no var) test failed\nexp: %s\ngot: %s", string(tplModNoVarExpect), string(r))
	}
}

func TestParseCtx(t *testing.T) {
	tree, _ := Parse(ctxOrigin, false)
	r := tree.HumanReadable()
	if !bytes.Equal(r, ctxExpect) {
		t.Errorf("ctx test failed\nexp: %s\ngot: %s", string(ctxExpect), string(r))
	}
}

func TestParseCntr(t *testing.T) {
	tree, _ := Parse(cntrOrigin, false)
	r := tree.HumanReadable()
	if !bytes.Equal(r, cntrExpect) {
		t.Errorf("cntr test failed\nexp: %s\ngot: %s", string(ctxExpect), string(r))
	}
}

func TestParseCondition(t *testing.T) {
	tree, _ := Parse(condOrigin, false)
	r := tree.HumanReadable()
	if !bytes.Equal(r, condExpect) {
		t.Errorf("condition test failed\nexp: %s\ngot: %s", string(condExpect), string(r))
	}

	treeNested, _ := Parse(condNestedOrigin, false)
	rNested := treeNested.HumanReadable()
	if !bytes.Equal(rNested, condNestedExpect) {
		t.Errorf("nested condition test failed\nexp: %s\ngot: %s", string(condNestedExpect), string(rNested))
	}
}

func TestParseLoop(t *testing.T) {
	tree, _ := Parse(loopOrigin, false)
	r := tree.HumanReadable()
	if !bytes.Equal(r, loopExpect) {
		t.Errorf("loop test failed\nexp: %s\ngot: %s", string(loopExpect), string(r))
	}

	treeSep, _ := Parse(loopSepOrigin, false)
	rSep := treeSep.HumanReadable()
	if !bytes.Equal(rSep, loopSepExpect) {
		t.Errorf("loop with sep test failed\nexp: %s\ngot: %s", string(loopSepExpect), string(rSep))
	}
}

func TestParseSwitch(t *testing.T) {
	tree, _ := Parse(switchOrigin, false)
	r := tree.HumanReadable()
	if !bytes.Equal(r, switchExpect) {
		t.Errorf("switch test failed\nexp: %s\ngot: %s", string(switchExpect), string(r))
	}

	treeNC, _ := Parse(switchNoCondOrigin, false)
	rNC := treeNC.HumanReadable()
	if !bytes.Equal(rNC, switchNoCondExpect) {
		t.Errorf("switch no cond test failed\nexp: %s\ngot: %s", string(switchNoCondExpect), string(rNC))
	}
}
