package dyntpl

import (
	"github.com/koykov/fastconv"
)

var (
	// Symbols to replace.
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

	// Replacements.
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

// JSON quote of string value - '"' + JSON escape + '"'.
func modJsonQuote(ctx *Ctx, buf *interface{}, val interface{}, _ []interface{}) error {
	ctx.Buf.Reset().WriteB(jqQd)
	err := modJsonEscape(ctx, buf, val, nil)
	if err == nil {
		ctx.Buf.Write(ctx.Buf1)
	}
	ctx.Buf.WriteB(jqQd)
	*buf = &ctx.Buf
	return nil
}

// JSON escape of string value.
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
	ctx.Buf1.Reset()
	_ = b[l-1]
	for i := 0; i < l; i++ {
		switch b[i] {
		case jqQd:
			ctx.Buf1.Write(b[o:i]).Write(jqQdR)
			o = i + 1
		case jqSl:
			ctx.Buf1.Write(b[o:i]).Write(jqSlR)
			o = i + 1
		case jqNl:
			ctx.Buf1.Write(b[o:i]).Write(jqNlR)
			o = i + 1
		case jqCr:
			ctx.Buf1.Write(b[o:i]).Write(jqCrR)
			o = i + 1
		case jqT:
			ctx.Buf1.Write(b[o:i]).Write(jqTR)
			o = i + 1
		case jqFf:
			ctx.Buf1.Write(b[o:i]).Write(jqFfR)
			o = i + 1
		case jqBs:
			ctx.Buf1.Write(b[o:i]).Write(jqBsR)
			o = i + 1
		case jqLt:
			ctx.Buf1.Write(b[o:i]).Write(jqLtR)
			o = i + 1
		case jqQs:
			ctx.Buf1.Write(b[o:i]).Write(jqQsR)
			o = i + 1
		case jqZ:
			ctx.Buf1.Write(b[o:i]).Write(jqZR)
			o = i + 1
		}
	}
	ctx.Buf1.Write(b[o:])
	*buf = &ctx.Buf1

	return nil
}
