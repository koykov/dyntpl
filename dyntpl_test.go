package dyntpl

import (
	"bytes"
	"testing"

	"github.com/koykov/inspector"
)

func TestTpl(t *testing.T) {
	t.Run("condition", testTpl)
	t.Run("conditionHelper", testTpl)
	t.Run("conditionLC", testTpl)
	t.Run("conditionNoStatic", testTpl)
	t.Run("conditionOK", testTpl)
	t.Run("conditionStr", testTpl)
	t.Run("counter0", testTpl)
	t.Run("counter1", testTpl)
	t.Run("ctxOK", testTpl)
	t.Run("exit", testTpl)
	t.Run("includeHost", testTpl)
	t.Run("includeHostJS", testTpl)
	t.Run("loopCount", testTplLC)
	t.Run("loopCountBreak", testTpl)
	t.Run("loopCountBreakN", testTpl)
	t.Run("loopCountContinue", testTpl)
	t.Run("loopCountCtx", testTpl)
	t.Run("loopCountLazybreak", testTpl)
	t.Run("loopCountLazybreakN", testTpl)
	t.Run("loopCountStatic", testTpl)
	t.Run("loopRange", testTpl)
	t.Run("loopRangeLazybreakN", testTpl)
	t.Run("loopRangeElse", testTpl)
	t.Run("raw", testTpl)
	t.Run("simple", testTpl)
	t.Run("switch", testTpl)
	t.Run("switchNoCondition", testTpl)
	t.Run("field404", testTpl)
	t.Run("strAnyMap", func(t *testing.T) {
		key := getTBName(t)
		st := getStage(key)
		if st == nil {
			t.Error("stage not found")
			return
		}
		ctx := NewCtx()
		ctx.Set("map_", map[string]any{
			"x": map[string]any{
				"y": map[string]any{
					"z": []string{"my", "substrings"},
				},
			},
		}, inspector.StringAnyMapInspector{})
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
	})
}

func testTpl(t *testing.T) {
	key := getTBName(t)
	st := getStage(key)
	if st == nil {
		t.Error("stage not found")
		return
	}

	ctx := NewCtx()
	ctx.Set("user", user, ins)
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

func testTplLC(t *testing.T) {
	key := getTBName(t)
	st := getStage(key)
	if st == nil {
		t.Error("stage not found")
		return
	}

	ctx := NewCtx()
	ctx.Set("user", user, ins)
	ctx.SetStatic("begin", 0)
	ctx.SetStatic("end", 3)
	result, err := Render(key, ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, st.expect) {
		t.Errorf("%s mismatch", key)
	}
}

func BenchmarkTpl(b *testing.B) {
	b.Run("condition", benchTpl)
	b.Run("conditionHelper", benchTpl)
	b.Run("conditionLC", benchTpl)
	b.Run("conditionNoStatic", benchTpl)
	b.Run("conditionOK", benchTpl)
	b.Run("conditionStr", benchTpl)
	b.Run("counter0", benchTpl)
	b.Run("counter1", benchTpl)
	b.Run("ctxOK", benchTpl)
	b.Run("exit", benchTpl)
	b.Run("includeHost", benchTpl)
	b.Run("includeHostJS", benchTpl)
	b.Run("loopCount", benchTplLC)
	b.Run("loopCountBreak", benchTpl)
	b.Run("loopCountBreakN", benchTpl)
	b.Run("loopCountContinue", benchTpl)
	b.Run("loopCountCtx", benchTpl)
	b.Run("loopCountLazybreak", benchTpl)
	b.Run("loopCountLazybreakN", benchTpl)
	b.Run("loopCountStatic", benchTpl)
	b.Run("loopRange", benchTpl)
	b.Run("loopRangeLazybreakN", benchTpl)
	b.Run("loopRangeElse", benchTpl)
	b.Run("raw", benchTpl)
	b.Run("simple", benchTpl)
	b.Run("switch", benchTpl)
	b.Run("switchNoCondition", benchTpl)
	b.Run("field404", benchTpl)
	b.Run("strAnyMap", func(b *testing.B) {
		key := getTBName(b)
		st := getStage(key)
		if st == nil {
			b.Error("stage not found")
			return
		}
		map_ := map[string]any{
			"x": map[string]any{
				"y": map[string]any{
					"z": []string{"my", "substrings"},
				},
			},
		}
		var buf_ bytes.Buffer
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx := AcquireCtx()
			ctx.Set("map_", map_, inspector.StringAnyMapInspector{})
			err := Write(&buf_, key, ctx)
			if err != nil {
				b.Error(err)
			}
			result := buf_.Bytes()
			if len(st.expect) == 0 && len(result) != 0 {
				b.Errorf("%s mismatch", key)
				return
			}
			if !bytes.Equal(result, st.expect) {
				b.Errorf("%s mismatch", key)
			}
			ReleaseCtx(ctx)
			buf_.Reset()
		}
	})
}

func benchTpl(b *testing.B) {
	key := getTBName(b)

	st := getStage(key)
	if st == nil {
		b.Error("stage not found")
		return
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		ctx.Set("user", user, ins)
		buf.Reset()
		err := Write(&buf, key, ctx)
		if err != nil {
			b.Error(err)
		}
		if len(st.expect) == 0 && buf.Len() != 0 {
			b.Errorf("%s mismatch", key)
		}
		if !bytes.Equal(buf.Bytes(), st.expect) {
			b.Errorf("%s mismatch", key)
		}
		ReleaseCtx(ctx)
	}
}

func benchTplLC(b *testing.B) {
	key := getTBName(b)

	st := getStage(key)
	if st == nil {
		b.Error("stage not found")
		return
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		buf.Reset()
		ctx.Set("user", user, ins)
		ctx.SetStatic("begin", 0)
		ctx.SetStatic("end", 3)
		err := Write(&buf, key, ctx)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(buf.Bytes(), st.expect) {
			b.Errorf("%s mismatch", key)
		}
		ReleaseCtx(ctx)
	}
}
