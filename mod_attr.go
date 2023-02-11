package dyntpl

import "github.com/koykov/fastconv"

// Attribute escape.
func modAttrEscape(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	// Get count of encode iterations (cases: aa=, aaa=, AA=, AAA=, ...).
	itr := printIterations(args)

	// Get string to escape.
	var b []byte
	switch val.(type) {
	case []byte:
		b = val.([]byte)
	case *[]byte:
		b = *val.(*[]byte)
	case string:
		b = fastconv.S2B(val.(string))
	case *string:
		b = fastconv.S2B(*val.(*string))
	default:
		err = ErrModNoStr
		return
	}

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
				ctx.BufAcc.WriteStr("&#039;")
			default:
				ctx.BufAcc.WriteByte(b[j])
			}
		}
		b = ctx.BufAcc.StakedBytes()
	}
	ctx.BufModOut(buf, b)

	return
}
