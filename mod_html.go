package dyntpl

var (
	// Symbols to replace.
	heLt  = byte('<')
	heGt  = byte('>')
	heQd  = byte('"')
	heQs  = byte('\'')
	heAmp = byte('&')

	// Replacements.
	heLtR  = []byte("&lt;")
	heGtR  = []byte("&gt;")
	heQdR  = []byte("&quot;")
	heQsR  = []byte("&#39;")
	heAmpR = []byte("&amp;")
)

func modHTMLEscape(ctx *Ctx, buf *any, val any, args []any) error {
	var l, o int

	// Get count of encode iterations (cases: hh=, hhh=, ...).
	itr := printIterations(args)

	if ctx.BufAcc.StakeOut().WriteX(val).Error() != nil {
		return ErrModNoStr
	}
	b := ctx.BufAcc.StakedBytes()
	if l = len(b); l == 0 {
		return nil
	}
	for c := 0; c < itr; c++ {
		ctx.BufAcc.StakeOut()
		_ = b[l-1]
		for i := 0; i < l; i++ {
			c := b[i]
			if c == heLt {
				ctx.BufAcc.Write(b[o:i]).Write(heLtR)
				o = i + 1
			}
			if c == heGt {
				ctx.BufAcc.Write(b[o:i]).Write(heGtR)
				o = i + 1
			}
			if c == heQd {
				ctx.BufAcc.Write(b[o:i]).Write(heQdR)
				o = i + 1
			}
			if c == heQs {
				ctx.BufAcc.Write(b[o:i]).Write(heQsR)
				o = i + 1
			}
			if c == heAmp {
				ctx.BufAcc.Write(b[o:i]).Write(heAmpR)
				o = i + 1
			}
		}
		ctx.BufAcc.Write(b[o:])

		b = ctx.BufAcc.StakedBytes()
		l = len(b)
	}
	ctx.BufModOut(buf, b)

	return nil
}
