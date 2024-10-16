package dyntpl

import (
	"testing"
	"time"
)

var (
	dt97, _ = time.Parse("2006-01-02", "1997-04-19")
	dt0     = time.Unix(1136239445, 123456789).UTC()

	dtRFC3339 = []time.Time{
		time.Date(2008, 9, 17, 20, 4, 26, 0, time.UTC),
		time.Date(1994, 9, 17, 20, 4, 26, 0, time.FixedZone("EST", -18000)),
		time.Date(2000, 12, 26, 1, 15, 6, 0, time.FixedZone("OTO", 15600)),
	}

	loc, _   = time.LoadLocation("Europe/Moscow")
	dtNative = time.Unix(0, 1233810057012345600).In(loc)
	dtTZ     = time.Date(1994, 9, 17, 20, 4, 26, 0, time.FixedZone("EST", -18000))
	dtAdd    = time.Date(2012, 1, 21, 20, 4, 26, 555, time.UTC)
)

func TestModDatetime(t *testing.T) {
	t.Run("now", testMod)

	t.Run("datePercent", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateYearShort", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateYear", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateMonth", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateMonthNameShort", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateMonthName", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateWeekNumberSun", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateWeekNumberMon", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateDay", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateDaySpacePad", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateDayOfYear", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateDayOfWeek", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateDayOfWeekISO", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateDayNameShort", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateDayName", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateHour", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateHourSpacePad", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateHour12", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateHour12SpacePad", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateMinute", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateSecond", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateAM_PM", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("date_am_pm", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("datePreferredTime", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateUnixtime", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateComplex_r", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateComplexR", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateComplexT", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateComplexc", func(t *testing.T) { testModWA(t, modArgs{"date": dt0}) })
	t.Run("dateComplexD", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })
	t.Run("dateComplexF", func(t *testing.T) { testModWA(t, modArgs{"date": dt97}) })

	t.Run("dateRFC3339_0", func(t *testing.T) { testModWA(t, modArgs{"date": dtRFC3339[0]}) })
	t.Run("dateRFC3339_1", func(t *testing.T) { testModWA(t, modArgs{"date": dtRFC3339[1]}) })
	t.Run("dateRFC3339_2", func(t *testing.T) { testModWA(t, modArgs{"date": dtRFC3339[2]}) })

	t.Run("dateLayout", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateANSIC", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateUnixDate", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateRubyDate", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateRFC822", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateRFC822Z", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateRFC850", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateRFC1123", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateRFC1123Z", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateRFC3339", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateRFC3339Nano", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateKitchen", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateStamp", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateStampMilli", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateStampMicro", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateStampNano", func(t *testing.T) { testModWA(t, modArgs{"date": dtNative}) })
	t.Run("dateLayoutTZ", func(t *testing.T) { testModWA(t, modArgs{"date": dtTZ}) })

	t.Run("addNS", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addUS", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addMS", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addS", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addM", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addH", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addD", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addW", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addMM", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addY", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addC", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addMIL", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
	t.Run("addMixed", func(t *testing.T) { testModWA(t, modArgs{"date": dtAdd}) })
}

func BenchmarkModDatetime(b *testing.B) {
	b.Run("now", benchMod)

	b.Run("dateComplexR", func(b *testing.B) { benchModWA(b, modArgs{"date": dt0}) })
	b.Run("dateStampNano", func(b *testing.B) { benchModWA(b, modArgs{"date": dtNative}) })

	b.Run("addMixedBench", func(b *testing.B) { benchModWA(b, modArgs{"date": dtAdd}) })
}
