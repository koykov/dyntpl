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
	ctx.bbuf = ctx.bbuf[:0]
	ctx.bbuf = append(ctx.bbuf, jqQd)
	err := modJsonEscape(ctx, buf, val, nil)
	if err == nil {
		ctx.bbuf = append(ctx.bbuf, ctx.bbuf1...)
	}
	ctx.bbuf = append(ctx.bbuf, jqQd)
	*buf = &ctx.bbuf
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
	ctx.bbuf1 = ctx.bbuf1[:0]
	_ = b[l-1]
	for i := 0; i < l; i++ {
		switch b[i] {
		case jqQd:
			ctx.bbuf1 = append(ctx.bbuf1, b[o:i]...)
			ctx.bbuf1 = append(ctx.bbuf1, jqQdR...)
			o = i + 1
		case jqSl:
			ctx.bbuf1 = append(ctx.bbuf1, b[o:i]...)
			ctx.bbuf1 = append(ctx.bbuf1, jqSlR...)
			o = i + 1
		case jqNl:
			ctx.bbuf1 = append(ctx.bbuf1, b[o:i]...)
			ctx.bbuf1 = append(ctx.bbuf1, jqNlR...)
			o = i + 1
		case jqCr:
			ctx.bbuf1 = append(ctx.bbuf1, b[o:i]...)
			ctx.bbuf1 = append(ctx.bbuf1, jqCrR...)
			o = i + 1
		case jqT:
			ctx.bbuf1 = append(ctx.bbuf1, b[o:i]...)
			ctx.bbuf1 = append(ctx.bbuf1, jqTR...)
			o = i + 1
		case jqFf:
			ctx.bbuf1 = append(ctx.bbuf1, b[o:i]...)
			ctx.bbuf1 = append(ctx.bbuf1, jqFfR...)
			o = i + 1
		case jqBs:
			ctx.bbuf1 = append(ctx.bbuf1, b[o:i]...)
			ctx.bbuf1 = append(ctx.bbuf1, jqBsR...)
			o = i + 1
		case jqLt:
			ctx.bbuf1 = append(ctx.bbuf1, b[o:i]...)
			ctx.bbuf1 = append(ctx.bbuf1, jqLtR...)
			o = i + 1
		case jqQs:
			ctx.bbuf1 = append(ctx.bbuf1, b[o:i]...)
			ctx.bbuf1 = append(ctx.bbuf1, jqQsR...)
			o = i + 1
		case jqZ:
			ctx.bbuf1 = append(ctx.bbuf1, b[o:i]...)
			ctx.bbuf1 = append(ctx.bbuf1, jqZR...)
			o = i + 1
		}
	}
	ctx.bbuf1 = append(ctx.bbuf1, b[o:]...)
	*buf = &ctx.bbuf1

	return nil
}
