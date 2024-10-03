package dyntpl

import (
	"bytes"
	"testing"
)

func TestParser(t *testing.T) {
	t.Run("cutComments", testParser)
	t.Run("cutFmt", testParser)
	t.Run("printVar", testParser)
	t.Run("unexpectedEOF", testParser)
	t.Run("prefixSuffix", testParser)
	t.Run("exit", testParser)
	t.Run("mod", testParser)
	t.Run("modNoVar", testParser)
	t.Run("modNestedArg", testParser)
	t.Run("modPrint", testParser)
	t.Run("modPrintMulti", testParser)
	t.Run("modPrintMix", testParser)
	t.Run("modCallback", testParser)
	t.Run("namespace", testParser)
	t.Run("ctxDot", testParser)
	t.Run("ctxDot1", testParser)
	t.Run("ctxQB", testParser)
	t.Run("ctxModDot", testParser)
	t.Run("ctxAsOK", testParser)
	t.Run("ctx", testParser)
	t.Run("counter", testParser)
	t.Run("condition", testParser)
	t.Run("conditionStr", testParser)
	t.Run("conditionNested", testParser)
	t.Run("conditionOK", testParser)
	t.Run("conditionNotOK", testParser)
	t.Run("conditionLen", testParser)
	t.Run("conditionCap", testParser)
	t.Run("loop", testParser)
	t.Run("loopSeparator", testParser)
	t.Run("loopBreak", testParser)
	t.Run("loopBreakNested", testParser)
	t.Run("loopElse", testParser)
	t.Run("loopElseMixedCond", testParser)
	t.Run("loopBreakIf", testParser)
	t.Run("switch", testParser)
	t.Run("switchNoCondition", testParser)
	t.Run("switchNoConditionWithHelper", testParser)
	t.Run("include", testParser)
	t.Run("includeDot", testParser)
	t.Run("ctlSym", testParser)
	t.Run("raw", testParser)
}

func testParser(t *testing.T) {
	key := getTBName(t)
	st := getStage("parser/" + key)
	if st == nil {
		t.Error("stage not found")
		return
	}
	if len(st.expect) > 0 {
		tree, err := Parse(st.origin, false)
		if err != nil {
			t.Error(err)
		}
		r := tree.HumanReadable()
		if !bytes.Equal(r, st.expect) {
			t.Errorf("%s test failed\nexp: %s\ngot: %s", key, string(st.expect), string(r))
		}
	} else if len(st.raw) > 0 {
		p := &parser{tpl: st.origin}
		p.cutComments()
		p.cutFmt()
		if !bytes.Equal(st.raw, p.tpl) {
			t.Errorf("%s test raw failed\nexp: %s\ngot: %s", key, string(st.expect), string(p.tpl))
		}
	} else if len(st.err) > 0 {
		if _, err := Parse(st.origin, false); err != nil {
			if err.Error() != st.err {
				t.Errorf("%s test error failed\nexp err: %s\ngot: %s", key, err, st.err)
			}
		}
	}
}
