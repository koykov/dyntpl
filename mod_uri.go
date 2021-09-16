package dyntpl

const (
	// Hex digits in upper case.
	hexUp = "0123456789ABCDEF"
)

// Link escape string value.
func modLinkEscape(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	// Get count of encode iterations (cases: ll=, lll=, ...).
	itr := printIterations(args)

	// Get the source.
	if ctx.BufAcc.StakeOut().WriteX(val).Error() != nil {
		return ErrModNoStr
	}
	b := ctx.BufAcc.StakedBytes()
	l := len(b)
	if l == 0 {
		return
	}
	for c := 0; c < itr; c++ {
		ctx.BufAcc.StakeOut()
		_ = b[l-1]
		for i := 0; i < l; i++ {
			if b[i] == '"' {
				ctx.BufAcc.WriteStr(`\"`)
			} else if b[i] == ' ' {
				ctx.BufAcc.WriteByte('+')
			} else {
				ctx.BufAcc.WriteByte(b[i])
			}
		}
		b = ctx.BufAcc.StakedBytes()
		l = len(b)
	}
	ctx.BufModOut(buf, b)

	return
}

// URL encode string value.
//
// see https://golang.org/src/net/url/url.go#L100
func modURLEncode(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	// Get count of encode iterations (cases: uu=, uuu=, ...).
	itr := printIterations(args)

	// Get the source.
	if ctx.BufAcc.StakeOut().WriteX(val).Error() != nil {
		return ErrModNoStr
	}
	b := ctx.BufAcc.StakedBytes()
	l := len(b)
	if l == 0 {
		return
	}
	for c := 0; c < itr; c++ {
		ctx.BufAcc.StakeOut()
		_ = b[l-1]
		for i := 0; i < l; i++ {
			if b[i] >= 'a' && b[i] <= 'z' || b[i] >= 'A' && b[i] <= 'Z' ||
				b[i] >= '0' && b[i] <= '9' || b[i] == '-' || b[i] == '.' || b[i] == '_' {
				ctx.BufAcc.WriteByte(b[i])
			} else if b[i] == ' ' {
				ctx.BufAcc.WriteByte('+')
			} else {
				ctx.BufAcc.WriteByte('%').WriteByte(hexUp[b[i]>>4]).WriteByte(hexUp[b[i]&15])
			}
		}
		b = ctx.BufAcc.StakedBytes()
		l = len(b)
	}
	ctx.BufModOut(buf, b)

	return
}
