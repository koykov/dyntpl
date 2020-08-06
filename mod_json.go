package dyntpl

import (
	"github.com/koykov/any2bytes"
	"github.com/koykov/bytealg"
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
	jqBs = byte('\b')
	jqLt = byte('<')
	jqQs = byte('\'')
	jqZ  = byte(0)

	// Replacements.
	jqQdR = []byte(`\"`)
	jqSlR = []byte(`\\`)
	jqNlR = []byte(`\n`)
	jqCrR = []byte(`\r`)
	jqTR  = []byte(`\t`)
	jqFfR = []byte(`\u000c`)
	jqBsR = []byte(`\u0008`)
	jqLtR = []byte(`\u003c`)
	jqQsR = []byte(`\u0027`)
	jqZR  = []byte(`\u0000`)
)

// JSON quote of string value - '"' + JSON escape + '"'.
func modJsonQuote(ctx *Ctx, buf *interface{}, val interface{}, _ []interface{}) error {
	ctx.Buf.Reset().WriteByte(jqQd)
	err := modJsonEscape(ctx, buf, val, nil)
	if err == nil {
		ctx.Buf.Write(ctx.Buf1)
	}
	ctx.Buf.WriteByte(jqQd)
	*buf = &ctx.Buf
	return nil
}

// JSON escape of string value.
func modJsonEscape(ctx *Ctx, buf *interface{}, val interface{}, _ []interface{}) error {
	var (
		b   []byte
		err error
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
	ctx.Buf1.Reset()
	ctx.Buf1 = jsonEscape(b, ctx.Buf1)
	if ctx.chJQ {
		// Double escape when "jsonquote" bonds found.
		ctx.Buf2.Reset()
		ctx.Buf2 = jsonEscape(ctx.Buf1.Bytes(), ctx.Buf2)
		*buf = &ctx.Buf2
	} else {
		*buf = &ctx.Buf1
	}

	return nil
}

// Internal JSON escape helper.
func jsonEscape(b []byte, buf bytealg.ChainBuf) bytealg.ChainBuf {
	var o int
	l := len(b)
	if l == 0 {
		return buf
	}
	buf.Reset()
	_ = b[l-1]
	for i := 0; i < l; i++ {
		switch b[i] {
		case jqQd:
			buf.Write(b[o:i]).Write(jqQdR)
			o = i + 1
		case jqSl:
			buf.Write(b[o:i]).Write(jqSlR)
			o = i + 1
		case jqNl:
			buf.Write(b[o:i]).Write(jqNlR)
			o = i + 1
		case jqCr:
			buf.Write(b[o:i]).Write(jqCrR)
			o = i + 1
		case jqT:
			buf.Write(b[o:i]).Write(jqTR)
			o = i + 1
		case jqFf:
			buf.Write(b[o:i]).Write(jqFfR)
			o = i + 1
		case jqBs:
			buf.Write(b[o:i]).Write(jqBsR)
			o = i + 1
		case jqLt:
			buf.Write(b[o:i]).Write(jqLtR)
			o = i + 1
		case jqQs:
			buf.Write(b[o:i]).Write(jqQsR)
			o = i + 1
		case jqZ:
			buf.Write(b[o:i]).Write(jqZR)
			o = i + 1
		}
	}
	buf.Write(b[o:])
	return buf
}
