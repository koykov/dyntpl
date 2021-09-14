package dyntpl

import (
	"bytes"
	"io/ioutil"
	"testing"
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

	stages := []stage{
		{key: "cutComments", fn: tstRaw},
		{key: "cutFmt", fn: tstRaw},
		{key: "printVar"},
		{key: "unexpectedEOF", fn: tstErr, err: ErrUnexpectedEOF},
		{key: "prefixSuffix"},
		{key: "exit"},
		{key: "mod"},
		{key: "modNestedArg"},
		{key: "ctxDot"},
		{key: "ctxDot1"},
		{key: "ctxModDot"},
		{key: "ctxAsOK"},
		{key: "ctx"},
		{key: "counter"},
		{key: "condition"},
		{key: "conditionStr"},
		{key: "conditionNested"},
		{key: "conditionOK"},
		{key: "conditionNotOK"},
		{key: "loop"},
		{key: "loopSeparator"},
		{key: "loopBreak"},
		{key: "loopBreakNested"},
		{key: "switch"},
		{key: "switchNoCondition"},
		{key: "switchNoConditionWithHelper"},
		{key: "include"},
		{key: "includeDot"},
		{key: "locale"},
	}

	for _, s := range stages {
		fn := s.fn
		if fn == nil {
			fn = tst
		}
		s.origin, _ = ioutil.ReadFile("testdata/parser/" + s.key + ".tpl")
		s.expect, _ = ioutil.ReadFile("testdata/parser/" + s.key + ".txt")
		t.Run(s.key, func(t *testing.T) { fn(t, &s) })
	}
}
