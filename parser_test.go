package cbytetpl

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
)

func TestParse_CutComments(t *testing.T) {
	tpl := []byte(`{# this is a test template #}
		Payload line #0
		{# some comment #}
		Payload line #1
		{# EOT #}`)
	exp := []byte(`		Payload line #0
		Payload line #1
`)
	p := &Parser{tpl: tpl}
	p.cutComments()
	if !bytes.Equal(exp, p.tpl) {
		t.Errorf("comment cut test failed\nexp: %s\ngot: %s", string(exp), string(p.tpl))
	}
}

func TestParse_CutFmt(t *testing.T) {
	p := &Parser{tpl: cutFmtOrigin}
	p.cutFmt()
	if !bytes.Equal(cutFmtExpect, p.tpl) {
		t.Errorf("fmt cut test failed\nexp: %s\ngot: %s", string(cutFmtExpect), string(p.tpl))
	}
}

func TestParse_Primitive(t *testing.T) {
	tree, _ := Parse(cutFmtOrigin, false)
	r := tree.humanReadable()
	if !bytes.Equal(r, primHrExpect) {
		t.Errorf("prim test failed\nexp: %s\ngot: %s", string(primHrExpect), string(r))
	}
}

func TestParse_uEOT(t *testing.T) {
	_, err := Parse(uEOTOrigin, false)
	if err != ErrUnexpectedEOF {
		t.Errorf("unexpected EOT fail\nexp: %s\ngot: %s", ErrUnexpectedEOF, err)
	}
}

func TestParse_PrefixSuffix(t *testing.T) {
	tree, _ := Parse(tplPS, false)
	r := tree.humanReadable()
	if !bytes.Equal(r, tplPSExpect) {
		t.Errorf("prefix/suffix test failed\nexp: %s\ngot: %s", string(tplPSExpect), string(r))
	}
}

func TestParse_Condition(t *testing.T) {
	tree, _ := Parse(condOrigin, false)
	r := tree.humanReadable()
	if !bytes.Equal(r, condExpect) {
		t.Errorf("condition test failed\nexp: %s\ngot: %s", string(condExpect), string(r))
	}

	treeNested, _ := Parse(condNestedOrigin, false)
	rNested := treeNested.humanReadable()
	if !bytes.Equal(rNested, condNestedExpect) {
		t.Errorf("nested condition test failed\nexp: %s\ngot: %s", string(condNestedExpect), string(rNested))
	}
}
