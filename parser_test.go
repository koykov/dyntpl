package dyntpl

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/koykov/bytealg"
)

func TestParser(t *testing.T) {
	type stage struct {
		key string
		origin,
		expect []byte
		fn  func(t *testing.T, key string, origin []byte, _ error)
		err error
	}

	fn := func(t *testing.T, key string, origin []byte, _ error) {
		expect, _ := ioutil.ReadFile("testdata/parser/" + key + ".xml")
		tree, _ := Parse(origin, false)
		r := tree.HumanReadable()
		if !bytes.Equal(r, expect) {
			t.Errorf("%s test failed\nexp: %s\ngot: %s", key, string(expect), string(r))
		}
	}
	fnRaw := func(t *testing.T, key string, origin []byte, _ error) {
		expect, _ := ioutil.ReadFile("testdata/parser/" + key + ".txt")
		expect = bytealg.Trim(expect, []byte("\n"))
		p := &Parser{tpl: origin}
		p.cutComments()
		p.cutFmt()
		if !bytes.Equal(expect, p.tpl) {
			t.Errorf("%s test raw failed\nexp: %s\ngot: %s", key, string(expect), string(p.tpl))
		}
	}
	fnErr := func(t *testing.T, key string, origin []byte, err error) {
		_, err1 := Parse(origin, false)
		if err != err1 {
			t.Errorf("%s test error failed\nexp err: %s\ngot: %s", key, err, err1)
		}
	}

	stages := []stage{
		{key: "cutComments", fn: fnRaw},
		{key: "cutFmt", fn: fnRaw},
		{key: "printVar"},
		{key: "unexpectedEOF", fn: fnErr, err: ErrUnexpectedEOF},
		{key: "prefixSuffix"},
		{key: "exit"},
		{key: "mod"},
		{key: "modNoVar"},
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
		t.Run(s.key, func(t *testing.T) {
			s.origin, _ = ioutil.ReadFile("testdata/parser/" + s.key + ".tpl")
			if s.fn == nil {
				s.fn = fn
			}
			s.fn(t, s.key, s.origin, s.err)
		})
	}
}
