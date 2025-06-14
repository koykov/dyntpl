package dyntpl

import "github.com/koykov/byteconv"

func modFmtFormat(ctx *Ctx, buf *any, _ any, args []any) error {
	if len(args) == 0 {
		return ErrModPoorArgs
	}
	var sfmt string
	switch x := args[0].(type) {
	case string:
		sfmt = x
	case *string:
		sfmt = *x
	case []byte:
		sfmt = byteconv.B2S(x)
	case *[]byte:
		sfmt = byteconv.B2S(*x)
	default:
		return nil
	}
	ctx.BufAcc.StakeOut().
		WriteFormat(sfmt, args[1:]...)
	ctx.BufModOut(buf, ctx.BufAcc.StakedBytes())
	return nil
}
