package dyntpl

func modFmtFormat(ctx *Ctx, buf *any, _ any, args []any) error {
	if len(args) == 0 {
		return ErrModPoorArgs
	}
	sfmt := args[0].(string)
	ctx.BufAcc.StakeOut().
		WriteFormat(sfmt, args[1:])
	ctx.BufModOut(buf, ctx.BufAcc.StakedBytes())
	return nil
}
