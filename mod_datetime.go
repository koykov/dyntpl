package dyntpl

import (
	"time"

	"github.com/koykov/clock"
)

func modNow(ctx *Ctx, buf *any, _ any, args []any) (err error) {
	if len(args) > 0 {
		a := ctx.BufAcc.StakeOut().WriteX(args[0]).StakedString()
		if a == "stuck" {
			ctx.BufT = time.Date(2020, 2, 23, 0, 0, 0, 0, time.UTC) // dyntpl birthdate
			*buf = &ctx.BufT
			return
		}
	}
	ctx.BufT = time.Now()
	*buf = &ctx.BufT
	return
}

func modDate(ctx *Ctx, buf *any, val any, args []any) (err error) {
	format := clock.Layout
	if len(args) > 0 {
		format = args[0].(string)
	}
	var dt time.Time
	switch x := val.(type) {
	case time.Time:
		dt = x
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
	case uint:
	case uint8:
	case uint16:
	case uint32:
	case uint64:
		dt = time.Unix(int64(x), 0)
	default:
		return
	}

	ctx.BufAcc.StakeOut().WriteTime(format, dt)
	ctx.BufModOut(buf, ctx.BufAcc.StakedBytes())

	return
}
