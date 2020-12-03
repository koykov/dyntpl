package dyntpl

import (
	"github.com/koykov/any2bytes"
	"github.com/koykov/fastconv"
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

// HTML escape of string value.
func modHtmlEscape(ctx *Ctx, buf *interface{}, val interface{}, _ []interface{}) error {
	var (
		b    []byte
		l, o int
		err  error
	)
	ctx.Buf2.Reset()
	if p, ok := ConvBytes(val); ok {
		b = p
	} else if s, ok := ConvStr(val); ok {
		b = fastconv.S2B(s)
	} else if ctx.Buf2, err = any2bytes.AnyToBytes(ctx.Buf2, val); err == nil {
		b = ctx.Buf2
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
		c := b[i]
		if c == heLt {
			ctx.Buf.Write(b[o:i]).Write(heLtR)
			o = i + 1
		}
		if c == heGt {
			ctx.Buf.Write(b[o:i]).Write(heGtR)
			o = i + 1
		}
		if c == heQd {
			ctx.Buf.Write(b[o:i]).Write(heQdR)
			o = i + 1
		}
		if c == heQs {
			ctx.Buf.Write(b[o:i]).Write(heQsR)
			o = i + 1
		}
		if c == heAmp {
			ctx.Buf.Write(b[o:i]).Write(heAmpR)
			o = i + 1
		}
	}
	ctx.Buf.Write(b[o:])
	*buf = &ctx.Buf

	return nil
}
