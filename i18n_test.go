package dyntpl

import (
	"bytes"
	"testing"

	"github.com/koykov/hash/fnv"
	"github.com/koykov/i18n"
)

type i18nStage struct {
	key string
	fn  func(tb testing.TB, st *i18nStage, db *i18n.DB)
}

func TestI18n(t *testing.T) {
	loadStages()

	i18nDB, _ := i18n.New(fnv.Hasher{})
	i18nDB.Set("en.messages.welcome", "Welcome, !user!")
	i18nDB.Set("ru.messages.welcome", "Привет, !user!")
	i18nDB.Set("en.pc.cpu", "no|yes")
	i18nDB.Set("en.me.age", "{0} you just born|[1,10] you're a child|[10,18] you're teenager|[18,40] you're adult|[40,80] you're old|[80,*] you're dead")

	stages := []i18nStage{
		{key: "i18n"},
		{key: "i18nPlural"},
		{key: "i18nPluralExt"},
		{key: "i18nSetLocale"},
	}

	for _, s := range stages {
		t.Run(s.key, func(t *testing.T) {
			if s.fn == nil {
				s.fn = testI18n
			}
			s.fn(t, &s, i18nDB)
		})
	}
}

func testI18n(tb testing.TB, st *i18nStage, db *i18n.DB) {
	st1 := getStage(st.key)
	if st1 == nil {
		tb.Error("stage not found")
		return
	}

	ctx := NewCtx()
	ctx.I18n("en", db)
	ctx.Set("user", user, &ins)
	ctx.SetStatic("cores", 4)
	ctx.SetStatic("years", 90)
	result, err := Render(st.key, ctx)
	if err != nil {
		tb.Error(err)
	}
	if !bytes.Equal(result, st1.expect) {
		tb.Errorf("%s mismatch", st.key)
	}
}

func BenchmarkI18n(b *testing.B) {
	loadStages()

	i18nDB, _ := i18n.New(fnv.Hasher{})
	i18nDB.Set("en.messages.welcome", "Welcome, !user!")
	i18nDB.Set("ru.messages.welcome", "Привет, !user!")
	i18nDB.Set("en.pc.cpu", "no|yes")
	i18nDB.Set("en.me.age", "{0} you just born|[1,10] you're a child|[10,18] you're teenager|[18,40] you're adult|[40,80] you're old|[80,*] you're dead")

	b.Run("i18n", func(b *testing.B) { benchI18n(b, i18nDB) })
	b.Run("i18nPlural", func(b *testing.B) { benchI18n(b, i18nDB) })
	b.Run("i18nPluralExt", func(b *testing.B) { benchI18n(b, i18nDB) })
	b.Run("i18nSetLocale", func(b *testing.B) { benchI18n(b, i18nDB) })
}

func benchI18n(tb testing.TB, db *i18n.DB) {
	b := interface{}(tb).(*testing.B)
	key := getTBName(tb)

	st := getStage(key)
	if st == nil {
		tb.Error("stage not found")
		return
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ctx := AcquireCtx()
		ctx.I18n("en", db)
		ctx.Set("user", user, &ins)
		ctx.SetStatic("cores", 4)
		ctx.SetStatic("years", 90)
		buf.Reset()
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
