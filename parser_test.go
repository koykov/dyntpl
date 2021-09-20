package dyntpl

import (
	"bytes"
	"testing"
)

func TestParser(t *testing.T) {
	t.Run("cutComments", func(t *testing.T) { testParser(t) })
	t.Run("cutFmt", func(t *testing.T) { testParser(t) })
	t.Run("printVar", func(t *testing.T) { testParser(t) })
	t.Run("unexpectedEOF", func(t *testing.T) { testParser(t) })
	t.Run("prefixSuffix", func(t *testing.T) { testParser(t) })
	t.Run("exit", func(t *testing.T) { testParser(t) })
	t.Run("mod", func(t *testing.T) { testParser(t) })
	t.Run("modNoVar", func(t *testing.T) { testParser(t) })
	t.Run("modNestedArg", func(t *testing.T) { testParser(t) })
	t.Run("ctxDot", func(t *testing.T) { testParser(t) })
	t.Run("ctxDot1", func(t *testing.T) { testParser(t) })
	t.Run("ctxModDot", func(t *testing.T) { testParser(t) })
	t.Run("ctxAsOK", func(t *testing.T) { testParser(t) })
	t.Run("ctx", func(t *testing.T) { testParser(t) })
	t.Run("counter", func(t *testing.T) { testParser(t) })
	t.Run("condition", func(t *testing.T) { testParser(t) })
	t.Run("conditionStr", func(t *testing.T) { testParser(t) })
	t.Run("conditionNested", func(t *testing.T) { testParser(t) })
	t.Run("conditionOK", func(t *testing.T) { testParser(t) })
	t.Run("conditionNotOK", func(t *testing.T) { testParser(t) })
	t.Run("loop", func(t *testing.T) { testParser(t) })
	t.Run("loopSeparator", func(t *testing.T) { testParser(t) })
	t.Run("loopBreak", func(t *testing.T) { testParser(t) })
	t.Run("loopBreakNested", func(t *testing.T) { testParser(t) })
	t.Run("switch", func(t *testing.T) { testParser(t) })
	t.Run("switchNoCondition", func(t *testing.T) { testParser(t) })
	t.Run("switchNoConditionWithHelper", func(t *testing.T) { testParser(t) })
	t.Run("include", func(t *testing.T) { testParser(t) })
	t.Run("includeDot", func(t *testing.T) { testParser(t) })
	t.Run("locale", func(t *testing.T) { testParser(t) })
}

func testParser(t *testing.T) {
	key := getTBName(t)
	st := getStage("parser/" + key)
	if st == nil {
		t.Error("stage not found")
		return
	}
	if len(st.expect) > 0 {
		tree, _ := Parse(st.origin, false)
		r := tree.HumanReadable()
		if !bytes.Equal(r, st.expect) {
			t.Errorf("%s test failed\nexp: %s\ngot: %s", key, string(st.expect), string(r))
		}
	} else if len(st.raw) > 0 {
		p := &Parser{tpl: st.origin}
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
