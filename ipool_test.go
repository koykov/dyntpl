package dyntpl

import (
	"bytes"
	"sync"
	"testing"
)

const __ipoolS = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris facilisis, massa a luctus feugiat, " +
	"urna nulla accumsan orci, vitae consequat lectus mi et nisl. Sed vehicula elit quam, vel luctus lacus " +
	"pellentesque sed. Aliquam ac erat quis quam interdum molestie vitae eget sapien. Nam sit amet metus turpis. " +
	"Ut gravida diam a lacus feugiat placerat vel sit amet diam. Duis hendrerit, quam vitae maximus ultricies, " +
	"augue erat lacinia ex, ac semper sem eros quis ligula. Morbi sit amet aliquam ligula. Phasellus bibendum " +
	"ipsum id risus iaculis porttitor at in libero. Mauris accumsan pellentesque sapien, eget laoreet urna posuere " +
	"vitae. In vitae pharetra orci, sed interdum metus. Donec facilisis lectus orci, sit amet placerat dolor tempor " +
	"feugiat."

type testBytebuf struct {
	p sync.Pool
}

func (p *testBytebuf) Get() any {
	raw := p.p.Get()
	if raw == nil {
		b := make([]byte, 0, 64)
		return &b
	}
	return raw.(*[]byte)
}

func (p *testBytebuf) Reset(x any) {
	b := x.(*[]byte)
	*b = (*b)[:0]
}

func (p *testBytebuf) Put(x any) {
	p.p.Put(x)
}

func init() {
	_ = RegisterPool("test_bytebuf", &testBytebuf{})
	RegisterModFn("__testCopyUsePool", "", func(ctx *Ctx, buf *any, val any, _ []any) error {
		raw, err := ctx.AcquireFrom("test_bytebuf")
		if err != nil {
			return err
		}
		bbuf := raw.(*[]byte)
		data := val.(string)
		*bbuf = append(*bbuf, data...)
		*buf = bbuf
		return nil
	})
	RegisterModFn("__testCopyNoPool", "", func(ctx *Ctx, buf *any, val any, _ []any) error {
		data := val.(string)
		bbuf := append([]byte(nil), data...)
		*buf = &bbuf
		return nil
	})

	registerTestStages("ipool")
}

func TestInternalPool(t *testing.T) {
	fn := func(t *testing.T) {
		key := getTBName(t)
		st := getStage(key)
		if st == nil {
			t.Error("stage not found")
			return
		}
		ctx := NewCtx()
		ctx.SetStatic("myVar", __ipoolS)
		result, err := Render(key, ctx)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(result, st.expect) {
			t.Errorf("%s mismatch", key)
		}
		ctx.Reset()
	}
	t.Run("ipoolUsePool", fn)
	t.Run("ipoolNoPool", fn)
}

func BenchmarkInternalPool(b *testing.B) {
	fn := func(b *testing.B) {
		key := getTBName(b)
		st := getStage(key)
		if st == nil {
			b.Error("stage not found")
			return
		}
		b.ReportAllocs()
		var buf bytes.Buffer
		for i := 0; i < b.N; i++ {
			ctx := AcquireCtx()
			ctx.SetStatic("myVar", __ipoolS)
			err := Write(&buf, key, ctx)
			if err != nil {
				b.Error(err)
			}
			if !bytes.Equal(buf.Bytes(), st.expect) {
				b.Errorf("%s mismatch", key)
			}
			ReleaseCtx(ctx)
			buf.Reset()
		}
	}
	b.Run("ipoolUsePool", fn)
	b.Run("ipoolNoPool", fn)
}
