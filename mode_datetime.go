package dyntpl

import (
	"time"

	"github.com/koykov/clock"
)

func modNow(_ *Ctx, buf *any, _ any, _ []any) (err error) {
	*buf = time.Now()
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

	fmtb := ctx.BufAcc.StakeOut().GrowDelta(128).StakedBytes()
	if fmtb, err = clock.AppendFormat(fmtb[:0], format, dt); err != nil {
		return
	}
	ctx.BufModOut(buf, fmtb)

	return
}
