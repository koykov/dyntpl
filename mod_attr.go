package dyntpl

import (
	"strconv"
	"unicode/utf8"
)

// Attribute escape.
func modAttrEscape(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
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
			case '&':
				ctx.BufAcc.WriteStr("&amp;")
			case '<':
				ctx.BufAcc.WriteStr("&lt;")
			case '>':
				ctx.BufAcc.WriteStr("&gt;")
			case '"':
				ctx.BufAcc.WriteStr("&quot;")
			default:
				if r != ',' && r != '.' && r != '-' && r != '_' && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
					if (r < 0x1f && r != '\t' && r != '\n' && r != '\r') || (r >= 0x7f && r <= 0x9f) {
						ctx.BufAcc.WriteStr("&#xFFFD;")
					} else if utf8.RuneLen(r) == 1 {
						ctx.BufAcc.WriteStr("&#x")
						ctx.Buf = strconv.AppendInt(*ctx.Buf.Reset(), int64(r), 16)
						if ctx.Buf.Len() < 2 {
							ctx.BufAcc.WriteByte('0')
						}
						ctx.BufAcc.Write(ctx.Buf).WriteByte(';')
					} else {
						ctx.BufAcc.WriteStr("&#x")
						ctx.Buf = strconv.AppendInt(*ctx.Buf.Reset(), int64(r), 16)
						if ctx.Buf.Len() < 4 {
							ctx.BufAcc.WriteByte('0')
						}
						ctx.BufAcc.Write(ctx.Buf).WriteByte(';')
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
