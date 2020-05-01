package cbytetpl

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
					152354345634,
					14.345241,
					[]byte("pay for domain"),
				},
				{
					153465345246,
					-3.0000342543,
					[]byte("got refund"),
				},
				{
					156436535640,
					2325242534.35324523,
					[]byte("maintenance"),
				},
			},
		},
	}
)

func TestCtx_Get(t *testing.T) {
	var (
		ins testobj_ins.TestObjectInspector
		raw interface{}
	)
	ctx := NewCtx()
	ctx.Set("obj", testO, &ins)

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
}

// func BenchmarkCtx_Get(b *testing.B) {
// 	var (
// 		ins testobj_ins.TestObjectInspector
// 		raw interface{}
// 	)
// 	ctx := NewCtx()
// 	ctx.Set("obj", testO, &ins)
//
// 	b.ReportAllocs()
// 	for i := 0; i < b.N; i++ {
// 		raw = ctx.Get("obj.Id")
// 		if ctx.Err != nil {
// 			b.Error("ctx get error", ctx.Err)
// 		}
// 		if *raw.(*string) != "foo" {
// 			b.Error("ctx get mismatch: obj.Id")
// 		}
//
// 		raw = ctx.Get("obj.Finance.Balance")
// 		if ctx.Err != nil {
// 			b.Error("ctx get error", ctx.Err)
// 		}
// 		if *raw.(*float64) != 9000 {
// 			b.Error("ctx get mismatch: obj.Finance.Balance")
// 		}
// 	}
// }
//
// func BenchmarkCtxPool_Get(b *testing.B) {
// 	var (
// 		ins testobj_ins.TestObjectInspector
// 		raw interface{}
// 	)
//
// 	b.ReportAllocs()
// 	for i := 0; i < b.N; i++ {
// 		ctx := CP.Get()
// 		ctx.Set("obj", testO, &ins)
//
// 		raw = ctx.Get("obj.Id")
// 		if ctx.Err != nil {
// 			b.Error("ctx get error", ctx.Err)
// 		}
// 		if *raw.(*string) != "foo" {
// 			b.Error("ctx get mismatch: obj.Id")
// 		}
//
// 		raw = ctx.Get("obj.Finance.Balance")
// 		if ctx.Err != nil {
// 			b.Error("ctx get error", ctx.Err)
// 		}
// 		if *raw.(*float64) != 9000 {
// 			b.Error("ctx get mismatch: obj.Finance.Balance")
// 		}
//
// 		CP.Put(ctx)
// 	}
// }
