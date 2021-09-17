package dyntpl

import (
	"bytes"
	"testing"

	"github.com/koykov/hash/fnv"
	"github.com/koykov/i18n"
)

type i18nStage struct {
	key string
	fn  func(t *testing.T, st *i18nStage, db *i18n.DB)
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

func testI18n(t *testing.T, st *i18nStage, db *i18n.DB) {
	st1 := getStage(st.key)
	if st1 == nil {
		t.Error("stage not found")
		return
	}

	ctx := NewCtx()
	ctx.I18n("en", db)
	ctx.Set("user", user, &ins)
	ctx.SetStatic("cores", 4)
	ctx.SetStatic("years", 90)
	result, err := Render(st.key, ctx)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(result, st1.expect) {
		t.Errorf("%s mismatch", st.key)
	}
}

func BenchmarkI18n(b *testing.B) {
	pretest()
	fn := func(b *testing.B, tplName string, expect []byte, errMsg string) {
		db, _ := i18n.New(fnv.Hasher{})
		db.Set("en.messages.welcome", "Welcome, !user!")
		db.Set("ru.messages.welcome", "Привет, !user!")
		db.Set("en.pc.cpu", "no|yes")
		db.Set("en.me.age", "{0} you just born|[1,10] you're a child|[10,18] you're teenager|[18,40] you're adult|[40,80] you're old|[80,*] you're dead")

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			ctx := AcquireCtx()
			ctx.I18n("en", db)
			ctx.Set("user", user, &ins)
			ctx.SetStatic("cores", 4)
			ctx.SetStatic("years", 90)
			buf.Reset()
			err := Write(&buf, tplName, ctx)
			if err != nil {
				b.Error(err)
			}
			if !bytes.Equal(buf.Bytes(), expect) {
				b.Error(errMsg)
			}
			ReleaseCtx(ctx)
		}
	}

	b.Run("tplI18n", func(b *testing.B) { fn(b, "tplI18n", expectI18n, "tplI18n mismatch") })
	b.Run("tplI18nPlural", func(b *testing.B) { fn(b, "tplI18nPlural", expectI18nPlural, "tplI18nPlural mismatch") })
	b.Run("tplI18nPluralExt", func(b *testing.B) { fn(b, "tplI18nPluralExt", expectI18nPluralExt, "tplI18nPluralExt mismatch") })
	b.Run("tplI18nSetLocale", func(b *testing.B) { fn(b, "tplI18nSetLocale", expectI18nSetLocale, "tplI18nSetLocale mismatch") })
}
