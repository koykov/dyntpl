package dyntpl

import (
	"testing"

	"github.com/koykov/inspector/testobj"
	"github.com/koykov/inspector/testobj_ins"
)

var (
	testO = &testobj.TestObject{
		Id:         "foo",
		Name:       []byte("bar"),
		Cost:       12.34,
		Permission: &testobj.TestPermission{15: true, 23: false},
		Flags: testobj.TestFlag{
			"export": 17,
			"ro":     4,
			"rw":     7,
			"Valid":  1,
		},
		Finance: &testobj.TestFinance{
			MoneyIn:  3200,
			MoneyOut: 1500.637657,
			Balance:  9000,
			History: []testobj.TestHistory{
				{
					DateUnix: 152354345634,
					Cost:     14.345241,
					Comment:  []byte("pay for domain"),
				},
				{
					DateUnix: 153465345246,
					Cost:     -3.0000342543,
					Comment:  []byte("got refund"),
				},
				{
					DateUnix: 156436535640,
					Cost:     2325242534.35324523,
					Comment:  []byte("maintenance"),
				},
			},
		},
	}
)

func TestCtx(t *testing.T) {
	var (
		ins testobj_ins.TestObjectInspector
		raw any
	)
	ctx := NewCtx()
	ctx.Set("obj", testO, &ins)

	t.Run("get", func(t *testing.T) {
		raw = ctx.Get("obj.Id")
		if ctx.Err != nil {
			t.Error("ctx get error", ctx.Err)
		}
		if *raw.(*string) != "foo" {
			t.Error("ctx get mismatch: obj.Id")
		}

		raw = ctx.Get("obj.Finance.Balance")
		if ctx.Err != nil {
			t.Error("ctx get error", ctx.Err)
		}
		if *raw.(*float64) != 9000 {
			t.Error("ctx get mismatch: obj.Finance.Balance")
		}
	})
}

func BenchmarkCtx(b *testing.B) {
	var (
		ins testobj_ins.TestObjectInspector
		raw any
	)
	b.Run("get", func(b *testing.B) {
		ctx := NewCtx()
		ctx.Set("obj", testO, &ins)

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchCtx(b, ctx, &raw)
		}
	})
	b.Run("getWithPool", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctx := AcquireCtx()
			ctx.Set("obj", testO, &ins)

			benchCtx(b, ctx, &raw)

			ReleaseCtx(ctx)
		}
	})
}

func benchCtx(b *testing.B, ctx *Ctx, buf *any) {
	raw := *buf
	raw = ctx.Get("obj.Id")
	if ctx.Err != nil {
		b.Error("ctx get error", ctx.Err)
	}
	if *raw.(*string) != "foo" {
		b.Error("ctx get mismatch: obj.Id")
	}

	raw = ctx.Get("obj.Finance.Balance")
	if ctx.Err != nil {
		b.Error("ctx get error", ctx.Err)
	}
	if *raw.(*float64) != 9000 {
		b.Error("ctx get mismatch: obj.Finance.Balance")
	}
}
