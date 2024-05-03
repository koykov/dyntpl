package dyntpl

// JS escape.
func modJSEscape(ctx *Ctx, buf *any, val any, args []any) (err error) {
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
				ctx.BufAcc.WriteString("\\\\")
			case '/':
				ctx.BufAcc.WriteString("\\/")
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
						ctx.BufAcc.WriteString("\\u0000")
						off := ctx.BufAcc.Len()
						ctx.BufAcc.WriteIntBase(int64(r), 16)
						delta := ctx.BufAcc.Len() - off
						hex := ctx.BufAcc.Bytes()[ctx.BufAcc.Len()-delta:]
						ctx.BufAcc.Reduce(delta * 2)
						ctx.BufAcc.Write(hex)
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
