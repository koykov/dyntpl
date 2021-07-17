package dyntpl

import "github.com/koykov/x2bytes"

const (
	// Hex digits in upper case.
	hexUp = "0123456789ABCDEF"
)

// Link escape string value.
func modLinkEscape(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	// Get count of encode iterations (cases: ll=, lll=, ...).
	itr := printIterations(args)

	// Get the source.
	if p, ok := ConvBytes(val); ok {
		ctx.buf = append(ctx.buf[:0], p...)
	} else if s, ok := ConvStr(val); ok {
		ctx.buf = append(ctx.buf[:0], s...)
	} else if ctx.Buf2, err = x2bytes.ToBytesWR(ctx.Buf2, val); err == nil {
		ctx.buf = append(ctx.buf[:0], ctx.Buf2...)
	} else {
		return ErrModNoStr
	}
	l := len(ctx.buf)
	if l == 0 {
		return ErrModEmptyStr
	}
	for c := 0; c < itr; c++ {
		ctx.Buf.Reset()
		_ = ctx.buf[l-1]
		for i := 0; i < l; i++ {
			if ctx.buf[i] == '"' {
				ctx.Buf.WriteStr(`\"`)
			} else if ctx.buf[i] == ' ' {
				ctx.Buf.WriteByte('+')
			} else {
				ctx.Buf.WriteByte(ctx.buf[i])
			}
		}
		ctx.buf = append(ctx.buf[:0], ctx.Buf...)
		l = ctx.Buf.Len()
	}
	*buf = &ctx.Buf
	return
}

// URL encode string value.
//
// see https://golang.org/src/net/url/url.go#L100
func modUrlEncode(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	// Get count of encode iterations (cases: uu=, uuu=, ...).
	itr := printIterations(args)

	// Get the source.
	if p, ok := ConvBytes(val); ok {
		ctx.buf = append(ctx.buf[:0], p...)
	} else if s, ok := ConvStr(val); ok {
		ctx.buf = append(ctx.buf[:0], s...)
	} else if ctx.Buf2, err = x2bytes.ToBytesWR(ctx.Buf2, val); err == nil {
		ctx.buf = append(ctx.buf[:0], ctx.Buf2...)
	} else {
		return ErrModNoStr
	}
	l := len(ctx.buf)
	if l == 0 {
		return ErrModEmptyStr
	}
	for c := 0; c < itr; c++ {
		ctx.Buf.Reset()
		_ = ctx.buf[l-1]
		for i := 0; i < l; i++ {
			if ctx.buf[i] >= 'a' && ctx.buf[i] <= 'z' || ctx.buf[i] >= 'A' && ctx.buf[i] <= 'Z' ||
				ctx.buf[i] >= '0' && ctx.buf[i] <= '9' || ctx.buf[i] == '-' || ctx.buf[i] == '.' || ctx.buf[i] == '_' {
				ctx.Buf.WriteByte(ctx.buf[i])
			} else if ctx.buf[i] == ' ' {
				ctx.Buf.WriteByte('+')
			} else {
				ctx.Buf.WriteByte('%').WriteByte(hexUp[ctx.buf[i]>>4]).WriteByte(hexUp[ctx.buf[i]&15])
			}
		}
		ctx.buf = append(ctx.buf[:0], ctx.Buf...)
		l = ctx.Buf.Len()
	}
	*buf = &ctx.Buf
	return
}
