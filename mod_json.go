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
	ctx.bbuf = ctx.bbuf[:0]
	_ = b[l-1]
	for i := 0; i < l; i++ {
		switch b[i] {
		case jqQd:
			ctx.bbuf = append(ctx.bbuf, b[o:i]...)
			ctx.bbuf = append(ctx.bbuf, jqQdR...)
			o = i + 1
		case jqSl:
			ctx.bbuf = append(ctx.bbuf, b[o:i]...)
			ctx.bbuf = append(ctx.bbuf, jqSlR...)
			o = i + 1
		case jqNl:
			ctx.bbuf = append(ctx.bbuf, b[o:i]...)
			ctx.bbuf = append(ctx.bbuf, jqNlR...)
			o = i + 1
		case jqCr:
			ctx.bbuf = append(ctx.bbuf, b[o:i]...)
			ctx.bbuf = append(ctx.bbuf, jqCrR...)
			o = i + 1
		case jqT:
			ctx.bbuf = append(ctx.bbuf, b[o:i]...)
			ctx.bbuf = append(ctx.bbuf, jqTR...)
			o = i + 1
		case jqFf:
			ctx.bbuf = append(ctx.bbuf, b[o:i]...)
			ctx.bbuf = append(ctx.bbuf, jqFfR...)
			o = i + 1
		case jqBs:
			ctx.bbuf = append(ctx.bbuf, b[o:i]...)
			ctx.bbuf = append(ctx.bbuf, jqBsR...)
			o = i + 1
		case jqLt:
			ctx.bbuf = append(ctx.bbuf, b[o:i]...)
			ctx.bbuf = append(ctx.bbuf, jqLtR...)
			o = i + 1
		case jqQs:
			ctx.bbuf = append(ctx.bbuf, b[o:i]...)
			ctx.bbuf = append(ctx.bbuf, jqQsR...)
			o = i + 1
		case jqZ:
			ctx.bbuf = append(ctx.bbuf, b[o:i]...)
			ctx.bbuf = append(ctx.bbuf, jqZR...)
			o = i + 1
		}
	}
	ctx.bbuf = append(ctx.bbuf, b[o:]...)
	*buf = &ctx.bbuf

	return nil
}
