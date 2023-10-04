package dyntpl

import (
	"bytes"
	"testing"

	"github.com/koykov/inspector"
)

func TestTpl(t *testing.T) {
	t.Run("condition", func(t *testing.T) { testTpl(t) })
	t.Run("conditionHelper", func(t *testing.T) { testTpl(t) })
	t.Run("conditionLC", func(t *testing.T) { testTpl(t) })
	t.Run("conditionNoStatic", func(t *testing.T) { testTpl(t) })
	t.Run("conditionOK", func(t *testing.T) { testTpl(t) })
	t.Run("conditionStr", func(t *testing.T) { testTpl(t) })
	t.Run("counter0", func(t *testing.T) { testTpl(t) })
	t.Run("counter1", func(t *testing.T) { testTpl(t) })
	t.Run("ctxOK", func(t *testing.T) { testTpl(t) })
	t.Run("exit", func(t *testing.T) { testTpl(t) })
	t.Run("includeHost", func(t *testing.T) { testTpl(t) })
	t.Run("includeHostJS", func(t *testing.T) { testTpl(t) })
	t.Run("loopCount", func(t *testing.T) { testTplLC(t) })
	t.Run("loopCountBreak", func(t *testing.T) { testTpl(t) })
	t.Run("loopCountBreakN", func(t *testing.T) { testTpl(t) })
	t.Run("loopCountContinue", func(t *testing.T) { testTpl(t) })
	t.Run("loopCountCtx", func(t *testing.T) { testTpl(t) })
	t.Run("loopCountLazybreak", func(t *testing.T) { testTpl(t) })
	t.Run("loopCountLazybreakN", func(t *testing.T) { testTpl(t) })
	t.Run("loopCountStatic", func(t *testing.T) { testTpl(t) })
	t.Run("loopRange", func(t *testing.T) { testTpl(t) })
	t.Run("loopRangeLazybreakN", func(t *testing.T) { testTpl(t) })
	t.Run("raw", func(t *testing.T) { testTpl(t) })
	t.Run("simple", func(t *testing.T) { testTpl(t) })
	t.Run("switch", func(t *testing.T) { testTpl(t) })
	t.Run("switchNoCondition", func(t *testing.T) { testTpl(t) })
	t.Run("field404", func(t *testing.T) { testTpl(t) })
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

func testTplLC(t *testing.T) {
	key := getTBName(t)
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
		t.Errorf("%s mismatch", key)
	}
}

func BenchmarkTpl(b *testing.B) {
	b.Run("condition", func(b *testing.B) { benchTpl(b) })
	b.Run("conditionHelper", func(b *testing.B) { benchTpl(b) })
	b.Run("conditionLC", func(b *testing.B) { benchTpl(b) })
	b.Run("conditionNoStatic", func(b *testing.B) { benchTpl(b) })
	b.Run("conditionOK", func(b *testing.B) { benchTpl(b) })
	b.Run("conditionStr", func(b *testing.B) { benchTpl(b) })
	b.Run("counter0", func(b *testing.B) { benchTpl(b) })
	b.Run("counter1", func(b *testing.B) { benchTpl(b) })
	b.Run("ctxOK", func(b *testing.B) { benchTpl(b) })
	b.Run("exit", func(b *testing.B) { benchTpl(b) })
	b.Run("includeHost", func(b *testing.B) { benchTpl(b) })
	b.Run("includeHostJS", func(b *testing.B) { benchTpl(b) })
	b.Run("loopCount", func(b *testing.B) { benchTplLC(b) })
	b.Run("loopCountBreak", func(b *testing.B) { benchTpl(b) })
	b.Run("loopCountBreakN", func(b *testing.B) { benchTpl(b) })
	b.Run("loopCountContinue", func(b *testing.B) { benchTpl(b) })
	b.Run("loopCountCtx", func(b *testing.B) { benchTpl(b) })
	b.Run("loopCountLazybreak", func(b *testing.B) { benchTpl(b) })
	b.Run("loopCountLazybreakN", func(b *testing.B) { benchTpl(b) })
	b.Run("loopCountStatic", func(b *testing.B) { benchTpl(b) })
	b.Run("loopRange", func(b *testing.B) { benchTpl(b) })
	b.Run("loopRangeLazybreakN", func(b *testing.B) { benchTpl(b) })
	b.Run("raw", func(b *testing.B) { benchTpl(b) })
	b.Run("simple", func(b *testing.B) { benchTpl(b) })
	b.Run("switch", func(b *testing.B) { benchTpl(b) })
	b.Run("switchNoCondition", func(b *testing.B) { benchTpl(b) })
	b.Run("field404", func(b *testing.B) { benchTpl(b) })
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
		ctx.Set("user", user, &ins)
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
		ctx.Set("user", user, &ins)
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
