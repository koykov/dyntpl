package dyntpl

import "strconv"

// CSS escape.
func modCSSEscape(ctx *Ctx, buf *any, val any, args []any) (err error) {
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
			case '\r':
				ctx.BufAcc.WriteString("\\D ")
			case '\n':
				ctx.BufAcc.WriteString("\\A ")
			case '\t':
				ctx.BufAcc.WriteString("\\9 ")
			case 0:
				ctx.BufAcc.WriteString("\\0 ")
			case ' ':
				ctx.BufAcc.WriteString("\\20 ")
			default:
				if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
					ctx.BufAcc.WriteString("\\")
					ctx.Buf = strconv.AppendInt(*ctx.Buf.Reset(), int64(r), 16)
					ctx.BufAcc.Write(ctx.Buf).WriteByte(' ')
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
