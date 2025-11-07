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
		format = ctx.BufAcc.StakeOut().WriteX(args[0]).StakedString()
	}
	dt, ok := dateConv(val)
	if !ok {
		return
	}

	if len(args) > 1 {
		if locName := ctx.BufAcc.StakeOut().WriteX(args[1]).StakedString(); len(locName) > 0 {
			if loc, err1 := clock.LoadLocation(locName); err1 == nil {
				_ = dt.In(loc)
			}
		}
	}

	ctx.BufAcc.StakeOut().WriteTime(format, dt)
	ctx.BufModOut(buf, ctx.BufAcc.StakedBytes())

	return
}

func modDateAdd(ctx *Ctx, buf *any, val any, args []any) (err error) {
	if len(args) == 0 {
		return ErrModPoorArgs
	}
	raw := ctx.BufAcc.StakeOut().WriteX(args[0]).StakedString()
	if len(raw) == 0 {
		return
	}
	if raw[0] == '+' {
		raw = raw[1:]
	}
	var d time.Duration
	if d, err = clock.Relative(raw); err != nil {
		return
	}
	t, ok := dateConv(val)
	if !ok {
		return
	}
	ctx.BufT = t.Add(d)
	*buf = &ctx.BufT
	return
}

func dateConv(val any) (t time.Time, ok bool) {
	ok = true
	switch x := val.(type) {
	case time.Time:
		t = x
	case *time.Time:
		t = *x
	case int:
		t = time.Unix(int64(x), 0)
	case int8:
		t = time.Unix(int64(x), 0)
	case int16:
		t = time.Unix(int64(x), 0)
	case int32:
		t = time.Unix(int64(x), 0)
	case int64:
		t = time.Unix(x, 0)
	case uint:
		t = time.Unix(int64(x), 0)
	case uint8:
		t = time.Unix(int64(x), 0)
	case uint16:
		t = time.Unix(int64(x), 0)
	case uint32:
		t = time.Unix(int64(x), 0)
	case uint64:
		t = time.Unix(int64(x), 0)
	case *int:
		t = time.Unix(int64(*x), 0)
	case *int8:
		t = time.Unix(int64(*x), 0)
	case *int16:
		t = time.Unix(int64(*x), 0)
	case *int32:
		t = time.Unix(int64(*x), 0)
	case *int64:
		t = time.Unix(*x, 0)
	case *uint:
		t = time.Unix(int64(*x), 0)
	case *uint8:
		t = time.Unix(int64(*x), 0)
	case *uint16:
		t = time.Unix(int64(*x), 0)
	case *uint32:
		t = time.Unix(int64(*x), 0)
	case *uint64:
		t = time.Unix(int64(*x), 0)
	default:
		ok = false
		return
	}
	return
}
