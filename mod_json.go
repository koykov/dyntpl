package dyntpl

import (
	"github.com/koykov/bytebuf"
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
	var b []byte
	if err := modJsonEscape(ctx, buf, val, nil); err == nil {
		b = ctx.AccBuf.StakeOut().
			WriteByte(jqQd).
			Write(ctx.OutBuf.Bytes()).
			WriteByte(jqQd).StakedBytes()
	}
	*buf = ctx.OutBuf.Reset().Write(b)

	return nil
}

// JSON escape of string value.
func modJsonEscape(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) error {
	// Get count of encode iterations (cases: jj=, jjj=, ...).
	itr := printIterations(args)

	if ctx.AccBuf.StakeOut().WriteX(val).Error() != nil {
		return ErrModNoStr
	}
	b := ctx.AccBuf.StakedBytes()
	if l := len(b); l == 0 {
		return nil
	}
	for c := 0; c < itr; c++ {
		ctx.AccBuf.StakeOut()
		jsonEscape(b, &ctx.AccBuf)
		b = ctx.AccBuf.StakedBytes()
	}
	if ctx.chJQ {
		// Double escape when "jsonquote" bonds found.
		ctx.AccBuf.StakeOut()
		jsonEscape(b, &ctx.AccBuf)
		b = ctx.AccBuf.StakedBytes()
	}
	*buf = ctx.OutBuf.Reset().Write(b)

	return nil
}

// Internal JSON escape helper.
func jsonEscape(b []byte, buf *bytebuf.AccumulativeBuf) *bytebuf.AccumulativeBuf {
	var o int
	l := len(b)
	if l == 0 {
		return buf
	}
	// buf.Reset()
	_ = b[l-1]
	for i := 0; i < l; i++ {
		c := b[i]
		if c == jqQd {
			buf.Write(b[o:i]).Write(jqQdR)
			o = i + 1
		}
		if c == jqSl {
			buf.Write(b[o:i]).Write(jqSlR)
			o = i + 1
		}
		if c == jqNl {
			buf.Write(b[o:i]).Write(jqNlR)
			o = i + 1
		}
		if c == jqCr {
			buf.Write(b[o:i]).Write(jqCrR)
			o = i + 1
		}
		if c == jqT {
			buf.Write(b[o:i]).Write(jqTR)
			o = i + 1
		}
		if c == jqFf {
			buf.Write(b[o:i]).Write(jqFfR)
			o = i + 1
		}
		if c == jqBs {
			buf.Write(b[o:i]).Write(jqBsR)
			o = i + 1
		}
		if c == jqLt {
			buf.Write(b[o:i]).Write(jqLtR)
			o = i + 1
		}
		if c == jqQs {
			buf.Write(b[o:i]).Write(jqQsR)
			o = i + 1
		}
		if c == jqZ {
			buf.Write(b[o:i]).Write(jqZR)
			o = i + 1
		}
	}
	buf.Write(b[o:])
	return buf
}
