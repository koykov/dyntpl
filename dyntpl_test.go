package dyntpl

import (
	"bytes"
	"testing"
)

type tplStage struct {
	key string
	fn  func(tb testing.TB, key string)
}

func TestTpl(t *testing.T) {
	loadStages()

	tplStages := []tplStage{
		{key: "condition"},
		{key: "conditionHlp"},
		{key: "conditionNoStatic"},
		{key: "conditionOK"},
		{key: "conditionStr"},
		{key: "counter0"},
		{key: "counter1"},
		{key: "ctxOK"},
		{key: "exit"},
		{key: "includeHost"},
		{key: "includeHostJS"},
		{key: "loopCount", fn: testTplLC},
		{key: "loopCountBreak"},
		{key: "loopCountBreakN"},
		{key: "loopCountContinue"},
		{key: "loopCountCtx"},
		{key: "loopCountLazybreak"},
		{key: "loopCountLazybreakN"},
		{key: "loopCountStatic"},
		{key: "loopRange"},
		{key: "loopRangeLazybreakN"},
		{key: "raw"},
		{key: "simple"},
		{key: "switch"},
		{key: "switchNoCondition"},
	}

	for _, s := range tplStages {
		t.Run(s.key, func(t *testing.T) {
			if s.fn == nil {
				s.fn = testTpl
			}
			s.fn(t, s.key)
		})
	}
}

func testTpl(tb testing.TB, key string) {
	st := getStage(key)
	if st == nil {
		tb.Error("stage not found")
		return
	}

	ctx := NewCtx()
	ctx.Set("user", user, &ins)
	result, err := Render(key, ctx)
	if err != nil {
		tb.Error(err)
	}
	if len(st.expect) == 0 && len(result) != 0 {
		tb.Errorf("%s mismatch", key)
		return
	}
	if !bytes.Equal(result, st.expect) {
		tb.Errorf("%s mismatch", key)
	}
}

func testTplLC(t testing.TB, key string) {
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
	loadStages()

	b.Run("condition", func(b *testing.B) { benchTpl(b) })
	b.Run("conditionHlp", func(b *testing.B) { benchTpl(b) })
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
}

func benchTpl(tb testing.TB) {
	b := interface{}(tb).(*testing.B)
	key := getTBName(b)

	st := getStage(key)
	if st == nil {
		tb.Error("stage not found")
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

func benchTplLC(tb testing.TB) {
	b := interface{}(tb).(*testing.B)
	key := getTBName(b)

	st := getStage(key)
	if st == nil {
		tb.Error("stage not found")
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
