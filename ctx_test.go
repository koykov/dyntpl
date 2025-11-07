package dyntpl

import (
	"sync"
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
	ctx.Set("obj", testO, ins)

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
	t.Run("defer", func(t *testing.T) {
		var c int
		const n = 10
		for i := 0; i < n; i++ {
			ctx.Defer(func() error {
				c++
				return nil
			})
		}
		_ = ctx.defer_()
		if c != n {
			t.Errorf("ctx defer mismatch: need %d, got %d", n, c)
		}
	})
	t.Run("defer pool", func(t *testing.T) {
		type x struct {
			payload []byte
		}
		p := sync.Pool{New: func() any {
			return &x{}
		}}
		var a, r int
		var ctx_ Ctx
		const n = 10
		for i := 0; i < n; i++ {
			raw_ := p.Get()
			a++
			x_ := raw_.(*x)
			ctx_.Defer(func() error {
				x_.payload = x_.payload[:0]
				p.Put(x_)
				r++
				return nil
			})
		}
		_ = ctx_.defer_()
		if a != r {
			t.Errorf("ctx defer pool mismatch: need %d/%d, got %d/%d", n, n, a, r)
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
	b.Run("defer", func(b *testing.B) {
		b.ReportAllocs()
		ctx := NewCtx()
		var c int
		for i := 0; i < b.N; i++ {
			ctx.Reset()
			c = 0
			ctx.Defer(func() error {
				c++
				return nil
			})
			_ = ctx.defer_()
			if c != 1 {
				b.Errorf("ctx defer mismatch: need %d, got %d", 1, c)
			}
		}
	})
	b.Run("defer pool", func(b *testing.B) {
		b.ReportAllocs()
		type x struct {
			payload []byte
		}
		p := sync.Pool{New: func() any {
			return &x{}
		}}
		var ctx_ Ctx
		for i := 0; i < b.N; i++ {
			ctx_.Reset()
			raw_ := p.Get()
			x_ := raw_.(*x)
			ctx_.Defer(func() error {
				x_.payload = x_.payload[:0]
				p.Put(x_)
				return nil
			})
			_ = ctx_.defer_()
		}
	})
}

func benchCtx(b *testing.B, ctx *Ctx, buf *any) {
	_ = *buf
	raw := ctx.Get("obj.Id")
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
