package dyntpl

import (
	"github.com/koykov/any2bytes"
)

var (
	// Symbols to replace.
	heLt  = byte('<')
	heGt  = byte('>')
	heQd  = byte('"')
	heQs  = byte('\'')
	heAmp = byte('&')

	// Replacements.
	heLtR  = []byte("&lt;")
	heGtR  = []byte("&gt;")
	heQdR  = []byte("&quot;")
	heQsR  = []byte("&#39;")
	heAmpR = []byte("&amp;")
)

func modHtmlEscape(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) error {
	var (
		l, o int
		err  error
	)

	// Get count of encode iterations (cases: hh=, hhh=, ...).
	itr := printIterations(args)

	ctx.Buf2.Reset()
	if p, ok := ConvBytes(val); ok {
		ctx.buf = append(ctx.buf[:0], p...)
	} else if s, ok := ConvStr(val); ok {
		ctx.buf = append(ctx.buf[:0], s...)
	} else if ctx.Buf2, err = any2bytes.AnyToBytes(ctx.Buf2, val); err == nil {
		ctx.buf = append(ctx.buf[:0], ctx.Buf2...)
	} else {
		return ErrModNoStr
	}
	l = len(ctx.buf)
	if l == 0 {
		return nil
	}
	ctx.Buf.Reset()
	for c := 0; c < itr; c++ {
		_ = ctx.buf[l-1]
		for i := 0; i < l; i++ {
			c := ctx.buf[i]
			if c == heLt {
				ctx.Buf.Write(ctx.buf[o:i]).Write(heLtR)
				o = i + 1
			}
			if c == heGt {
				ctx.Buf.Write(ctx.buf[o:i]).Write(heGtR)
				o = i + 1
			}
			if c == heQd {
				ctx.Buf.Write(ctx.buf[o:i]).Write(heQdR)
				o = i + 1
			}
			if c == heQs {
				ctx.Buf.Write(ctx.buf[o:i]).Write(heQsR)
				o = i + 1
			}
			if c == heAmp {
				ctx.Buf.Write(ctx.buf[o:i]).Write(heAmpR)
				o = i + 1
			}
		}
		ctx.Buf.Write(ctx.buf[o:])

		ctx.buf = append(ctx.buf[:0], ctx.Buf...)
		l = ctx.Buf.Len()
	}
	*buf = &ctx.Buf

	return nil
}
