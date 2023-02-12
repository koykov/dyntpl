package dyntpl

// Attribute escape.
func modAttrEscape(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	// Get count of encode iterations (cases: aa=, aaa=, AA=, AAA=, ...).
	itr := printIterations(args)

	// Get string to escape.
	if ctx.BufAcc.StakeOut().WriteX(val).Error() != nil {
		return ErrModNoStr
	}
	b := ctx.BufAcc.StakedBytes()

	// Apply escape.
	for i := 0; i < itr; i++ {
		ctx.BufAcc.StakeOut()
		for j := 0; j < len(b); j++ {
			switch b[j] {
			case '&':
				ctx.BufAcc.WriteStr("&amp;")
			case '<':
				ctx.BufAcc.WriteStr("&lt;")
			case '>':
				ctx.BufAcc.WriteStr("&gt;")
			case '"':
				ctx.BufAcc.WriteStr("&quot;")
			case '\'':
				ctx.BufAcc.WriteStr("&#x27;")
			case '`':
				ctx.BufAcc.WriteStr("&#x60;")
			case '!':
				ctx.BufAcc.WriteStr("&#x21;")
			case '@':
				ctx.BufAcc.WriteStr("&#x40;")
			case '$':
				ctx.BufAcc.WriteStr("&#x24;")
			case '%':
				ctx.BufAcc.WriteStr("&#x25;")
			case '(':
				ctx.BufAcc.WriteStr("&#x28;")
			case ')':
				ctx.BufAcc.WriteStr("&#x29;")
			case '=':
				ctx.BufAcc.WriteStr("&#x3D;")
			case '+':
				ctx.BufAcc.WriteStr("&#x2B;")
			case '{':
				ctx.BufAcc.WriteStr("&#x7B;")
			case '}':
				ctx.BufAcc.WriteStr("&#x7D;")
			case '[':
				ctx.BufAcc.WriteStr("&#x5B;")
			case ']':
				ctx.BufAcc.WriteStr("&#x5D;")
			case '#':
				ctx.BufAcc.WriteStr("&#x23;")
			case ';':
				ctx.BufAcc.WriteStr("&#x3B;")
			default:
				ctx.BufAcc.WriteByte(b[j])
			}
		}
		b = ctx.BufAcc.StakedBytes()
	}
	ctx.BufModOut(buf, b)

	return
}
