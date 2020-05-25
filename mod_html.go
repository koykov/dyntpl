package dyntpl

import "github.com/koykov/fastconv"

var (
	heLt  = byte('<')
	heGt  = byte('>')
	heQd  = byte('"')
	heQs  = byte('\'')
	heAmp = byte('&')

	heLtR  = []byte("&lt;")
	heGtR  = []byte("&gt;")
	heQdR  = []byte("&quot;")
	heQsR  = []byte("&#39;")
	heAmpR = []byte("&amp;")
)

func modHtmlEscape(ctx *Ctx, buf *interface{}, val interface{}, _ []interface{}) error {
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
	ctx.Bbuf = ctx.Bbuf[:0]
	_ = b[l-1]
	for i := 0; i < l; i++ {
		switch b[i] {
		case heLt:
			ctx.Bbuf = append(ctx.Bbuf, b[o:i]...)
			ctx.Bbuf = append(ctx.Bbuf, heLtR...)
			o = i + 1
		case heGt:
			ctx.Bbuf = append(ctx.Bbuf, b[o:i]...)
			ctx.Bbuf = append(ctx.Bbuf, heGtR...)
			o = i + 1
		case heQd:
			ctx.Bbuf = append(ctx.Bbuf, b[o:i]...)
			ctx.Bbuf = append(ctx.Bbuf, heQdR...)
			o = i + 1
		case heQs:
			ctx.Bbuf = append(ctx.Bbuf, b[o:i]...)
			ctx.Bbuf = append(ctx.Bbuf, heQsR...)
			o = i + 1
		case heAmp:
			ctx.Bbuf = append(ctx.Bbuf, b[o:i]...)
			ctx.Bbuf = append(ctx.Bbuf, heAmpR...)
			o = i + 1
		}
	}
	ctx.Bbuf = append(ctx.Bbuf, b[o:]...)
	*buf = &ctx.Bbuf

	return nil
}
