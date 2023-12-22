package dyntpl

import "testing"

func TestLocalDB(t *testing.T) {
	t.Run("io", func(t *testing.T) {
		db_ := localDB{}
		db_.SetLocal("i18n.locale", "ru-RU")
		db_.SetLocal("i18n.dir", "ltr")
		loc, dir := db_.GetLocal("i18n.locale"), db_.GetLocal("i18n.dir")
		if loc.(string) != "ru-RU" || dir.(string) != "ltr" {
			t.FailNow()
		}
	})
}

func BenchmarkLocalDB(b *testing.B) {
	b.Run("io", func(b *testing.B) {
		b.ReportAllocs()
		db_ := localDB{}
		for i := 0; i < b.N; i++ {
			db_.SetLocal("i18n.locale", "ru-RU")
			loc := db_.GetLocal("i18n.locale")
			if loc.(string) != "ru-RU" {
				b.FailNow()
			}
			db_.reset()
		}
	})
}
