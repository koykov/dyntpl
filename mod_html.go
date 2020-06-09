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
	ctx.Buf.Reset()
	_ = b[l-1]
	for i := 0; i < l; i++ {
		switch b[i] {
		case heLt:
			ctx.Buf.Write(b[o:i]).Write(heLtR)
			o = i + 1
		case heGt:
			ctx.Buf.Write(b[o:i]).Write(heGtR)
			o = i + 1
		case heQd:
			ctx.Buf.Write(b[o:i]).Write(heQdR)
			o = i + 1
		case heQs:
			ctx.Buf.Write(b[o:i]).Write(heQsR)
			o = i + 1
		case heAmp:
			ctx.Buf.Write(b[o:i]).Write(heAmpR)
			o = i + 1
		}
	}
	ctx.Buf.Write(b[o:])
	*buf = &ctx.Buf

	return nil
}
