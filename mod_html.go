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

func modHtmlEscape(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) error {
	var l, o int

	// Get count of encode iterations (cases: hh=, hhh=, ...).
	itr := printIterations(args)

	b := ctx.AccBuf.StakeOut().WriteX(val).StakedBytes()
	if ctx.AccBuf.Error() != nil {
		return ErrModNoStr
	}
	if l = len(b); l == 0 {
		return nil
	}
	for c := 0; c < itr; c++ {
		ctx.AccBuf.StakeOut()
		_ = b[l-1]
		for i := 0; i < l; i++ {
			c := b[i]
			if c == heLt {
				ctx.AccBuf.Write(b[o:i]).Write(heLtR)
				o = i + 1
			}
			if c == heGt {
				ctx.AccBuf.Write(b[o:i]).Write(heGtR)
				o = i + 1
			}
			if c == heQd {
				ctx.AccBuf.Write(b[o:i]).Write(heQdR)
				o = i + 1
			}
			if c == heQs {
				ctx.AccBuf.Write(b[o:i]).Write(heQsR)
				o = i + 1
			}
			if c == heAmp {
				ctx.AccBuf.Write(b[o:i]).Write(heAmpR)
				o = i + 1
			}
		}
		ctx.AccBuf.Write(b[o:])

		b = ctx.AccBuf.StakedBytes()
		l = len(b)
	}
	ctx.OutBuf.Reset().Write(b)
	*buf = &ctx.OutBuf

	return nil
}
