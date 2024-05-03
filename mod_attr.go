package dyntpl

import (
	"unicode/utf8"
)

// Attribute escape.
func modAttrEscape(ctx *Ctx, buf *any, val any, args []any) (err error) {
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
				ctx.BufAcc.WriteString("&amp;")
			case '<':
				ctx.BufAcc.WriteString("&lt;")
			case '>':
				ctx.BufAcc.WriteString("&gt;")
			case '"':
				ctx.BufAcc.WriteString("&quot;")
			default:
				if r != ',' && r != '.' && r != '-' && r != '_' && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
					if (r < 0x1f && r != '\t' && r != '\n' && r != '\r') || (r >= 0x7f && r <= 0x9f) {
						ctx.BufAcc.WriteString("&#xFFFD;")
					} else if utf8.RuneLen(r) == 1 {
						ctx.BufAcc.WriteString("&#x00")
						off := ctx.BufAcc.Len()
						ctx.BufAcc.WriteIntBase(int64(r), 16)
						delta := ctx.BufAcc.Len() - off
						hex := ctx.BufAcc.Bytes()[ctx.BufAcc.Len()-delta:]
						if delta < 2 {
							delta *= 2
						} else {
							delta += 2
						}
						ctx.BufAcc.Reduce(delta)
						ctx.BufAcc.Write(hex)
						ctx.BufAcc.WriteByte(';')
					} else {
						ctx.BufAcc.WriteString("&#x0000")
						off := ctx.BufAcc.Len()
						ctx.BufAcc.WriteIntBase(int64(r), 16)
						delta := ctx.BufAcc.Len() - off
						hex := ctx.BufAcc.Bytes()[ctx.BufAcc.Len()-delta:]
						if delta < 4 {
							delta *= 2
						} else {
							delta += 4
						}
						ctx.BufAcc.Reduce(delta)
						ctx.BufAcc.Write(hex)
						ctx.BufAcc.WriteByte(';')
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
