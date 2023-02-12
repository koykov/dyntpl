package dyntpl

import (
	"strconv"
)

// JS escape.
func modJSEscape(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	// Get count of encode iterations (cases: aa=, aaa=, AA=, AAA=, ...).
	itr := printIterations(args)

	// Get string to escape.
	if ctx.BufAcc.StakeOut().WriteX(val).Error() != nil {
		return ErrModNoStr
	}
	b := ctx.BufAcc.StakedString()

	// Apply escape.
	for i := 0; i < itr; i++ {
		ctx.BufAcc.StakeOut()
		for _, r := range b {
			switch r {
			case '\\':
				ctx.BufAcc.WriteStr("\\\\")
			case '/':
				ctx.BufAcc.WriteStr("\\/")
			case '\x08':
				ctx.BufAcc.WriteByte('\b')
			case '\x0C':
				ctx.BufAcc.WriteByte('\f')
			case '\x0A':
				ctx.BufAcc.WriteByte('\n')
			case '\x0D':
				ctx.BufAcc.WriteByte('\r')
			case '\x09':
				ctx.BufAcc.WriteByte('\t')
			default:
				if r != ',' && r != '.' && r != '_' && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
					wr := func(r int32) {
						ctx.BufAcc.WriteStr("\\u")
						ctx.Buf = strconv.AppendInt(*ctx.Buf.Reset(), int64(r), 16)
						if ctx.Buf.Len() == 1 {
							ctx.BufAcc.WriteStr("000")
						} else if ctx.Buf.Len() == 2 {
							ctx.BufAcc.WriteStr("00")
						} else if ctx.Buf.Len() == 3 {
							ctx.BufAcc.WriteByte('0')
						}
						ctx.BufAcc.Write(ctx.Buf)
					}
					if r < 0x10000 {
						wr(r)
					} else {
						u := r - 0x10000
						hi := 0xD800 | (u >> 10)
						lo := 0xDC00 | (u & 0x3FF)
						wr(hi)
						wr(lo)
					}
				} else {
					ctx.BufAcc.WriteByte(byte(r))
				}
			}
		}
		b = ctx.BufAcc.StakedString()
	}
	ctx.BufModStrOut(buf, b)

	return
}
