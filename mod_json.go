package dyntpl

import "github.com/koykov/fastconv"

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
	ctx.Bbuf = ctx.Bbuf[:0]
	ctx.Bbuf = append(ctx.Bbuf, jqQd)
	err := modJsonEscape(ctx, buf, val, nil)
	if err == nil {
		ctx.Bbuf = append(ctx.Bbuf, ctx.Bbuf1...)
	}
	ctx.Bbuf = append(ctx.Bbuf, jqQd)
	*buf = &ctx.Bbuf
	return nil
}

func modJsonEscape(ctx *Ctx, buf *interface{}, val interface{}, _ []interface{}) error {
	var (
		b    []byte
		l, o int
	)
	switch val.(type) {
	case []byte:
		b = val.([]byte)
	case *[]byte:
		b = *val.(*[]byte)
	case string:
		b = fastconv.S2B(val.(string))
	case *string:
		b = fastconv.S2B(*val.(*string))
	default:
		return ErrModNoStr
	}
	l = len(b)
	if l == 0 {
		return nil
	}
	ctx.Bbuf1 = ctx.Bbuf1[:0]
	_ = b[l-1]
	for i := 0; i < l; i++ {
		switch b[i] {
		case jqQd:
			ctx.Bbuf1 = append(ctx.Bbuf1, b[o:i]...)
			ctx.Bbuf1 = append(ctx.Bbuf1, jqQdR...)
			o = i + 1
		case jqSl:
			ctx.Bbuf1 = append(ctx.Bbuf1, b[o:i]...)
			ctx.Bbuf1 = append(ctx.Bbuf1, jqSlR...)
			o = i + 1
		case jqNl:
			ctx.Bbuf1 = append(ctx.Bbuf1, b[o:i]...)
			ctx.Bbuf1 = append(ctx.Bbuf1, jqNlR...)
			o = i + 1
		case jqCr:
			ctx.Bbuf1 = append(ctx.Bbuf1, b[o:i]...)
			ctx.Bbuf1 = append(ctx.Bbuf1, jqCrR...)
			o = i + 1
		case jqT:
			ctx.Bbuf1 = append(ctx.Bbuf1, b[o:i]...)
			ctx.Bbuf1 = append(ctx.Bbuf1, jqTR...)
			o = i + 1
		case jqFf:
			ctx.Bbuf1 = append(ctx.Bbuf1, b[o:i]...)
			ctx.Bbuf1 = append(ctx.Bbuf1, jqFfR...)
			o = i + 1
		case jqBs:
			ctx.Bbuf1 = append(ctx.Bbuf1, b[o:i]...)
			ctx.Bbuf1 = append(ctx.Bbuf1, jqBsR...)
			o = i + 1
		case jqLt:
			ctx.Bbuf1 = append(ctx.Bbuf1, b[o:i]...)
			ctx.Bbuf1 = append(ctx.Bbuf1, jqLtR...)
			o = i + 1
		case jqQs:
			ctx.Bbuf1 = append(ctx.Bbuf1, b[o:i]...)
			ctx.Bbuf1 = append(ctx.Bbuf1, jqQsR...)
			o = i + 1
		case jqZ:
			ctx.Bbuf1 = append(ctx.Bbuf1, b[o:i]...)
			ctx.Bbuf1 = append(ctx.Bbuf1, jqZR...)
			o = i + 1
		}
	}
	ctx.Bbuf1 = append(ctx.Bbuf1, b[o:]...)
	*buf = &ctx.Bbuf1

	return nil
}
