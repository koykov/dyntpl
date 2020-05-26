package dyntpl

import (
	"github.com/koykov/fastconv"
)

var (
	jqQd = byte('"')
	jqSl = byte('\\')
	jqNl = byte('\n')
	jqCr = byte('\r')
	jqT  = byte('\t')
	jqFf = byte('\f')
	jqBs = byte('\f')
	jqLt = byte('<')
	jqQs = byte('\'')
	jqZ  = byte(0)

	jqQdR = []byte(`\"`)
	jqSlR = []byte("\\")
	jqNlR = []byte("\n")
	jqCrR = []byte("\r")
	jqTR  = []byte("\t")
	jqFfR = []byte("\u000c")
	jqBsR = []byte("\u0008")
	jqLtR = []byte("\u003c")
	jqQsR = []byte("\u0027")
	jqZR  = []byte("\u0000")
)

func modJsonQuote(ctx *Ctx, buf *interface{}, val interface{}, _ []interface{}) error {
	ctx.Bbuf.ResetWriteByte(jqQd)
	err := modJsonEscape(ctx, buf, val, nil)
	if err == nil {
		ctx.Bbuf.Write(ctx.Bbuf1)
	}
	ctx.Bbuf.WriteByte(jqQd)
	*buf = &ctx.Bbuf
	return nil
}

func modJsonEscape(ctx *Ctx, buf *interface{}, val interface{}, _ []interface{}) error {
	var (
		b    []byte
		l, o int
	)
	if p, ok := ConvBytes(val); ok {
		b = p
	} else if s, ok := ConvStr(val); ok {
		b = fastconv.S2B(s)
	} else {
		return ErrModNoStr
	}
	l = len(b)
	if l == 0 {
		return nil
	}
	ctx.Bbuf1.Reset()
	_ = b[l-1]
	for i := 0; i < l; i++ {
		switch b[i] {
		case jqQd:
			ctx.Bbuf1.Write(b[o:i])
			ctx.Bbuf1.Write(jqQdR)
			o = i + 1
		case jqSl:
			ctx.Bbuf1.Write(b[o:i])
			ctx.Bbuf1.Write(jqSlR)
			o = i + 1
		case jqNl:
			ctx.Bbuf1.Write(b[o:i])
			ctx.Bbuf1.Write(jqNlR)
			o = i + 1
		case jqCr:
			ctx.Bbuf1.Write(b[o:i])
			ctx.Bbuf1.Write(jqCrR)
			o = i + 1
		case jqT:
			ctx.Bbuf1.Write(b[o:i])
			ctx.Bbuf1.Write(jqTR)
			o = i + 1
		case jqFf:
			ctx.Bbuf1.Write(b[o:i])
			ctx.Bbuf1.Write(jqFfR)
			o = i + 1
		case jqBs:
			ctx.Bbuf1.Write(b[o:i])
			ctx.Bbuf1.Write(jqBsR)
			o = i + 1
		case jqLt:
			ctx.Bbuf1.Write(b[o:i])
			ctx.Bbuf1.Write(jqLtR)
			o = i + 1
		case jqQs:
			ctx.Bbuf1.Write(b[o:i])
			ctx.Bbuf1.Write(jqQsR)
			o = i + 1
		case jqZ:
			ctx.Bbuf1.Write(b[o:i])
			ctx.Bbuf1.Write(jqZR)
			o = i + 1
		}
	}
	ctx.Bbuf1.Write(b[o:])
	*buf = &ctx.Bbuf1

	return nil
}
