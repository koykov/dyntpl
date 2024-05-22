package dyntpl

import (
	"strconv"

	"github.com/koykov/byteconv"
)

func modAbs(ctx *Ctx, buf *any, val any, _ []any) (err error) {
	f, ok := floatConv(val)
	if !ok {
		return
	}
	if f < 0 {
		f = -f
	}
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func modInc(ctx *Ctx, buf *any, val any, _ []any) (err error) {
	f, ok := floatConv(val)
	if !ok {
		return
	}
	f += 1
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func modDec(ctx *Ctx, buf *any, val any, _ []any) (err error) {
	f, ok := floatConv(val)
	if !ok {
		return
	}
	f -= 1
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func modMathAdd(ctx *Ctx, buf *any, val any, args []any) (err error) {
	var (
		f, d float64
		ok   bool
	)
	if f, d, err, ok = mathConv2(val, args); !ok {
		return
	}
	f += d
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func modMathSub(ctx *Ctx, buf *any, val any, args []any) (err error) {
	var (
		f, d float64
		ok   bool
	)
	if f, d, err, ok = mathConv2(val, args); !ok {
		return
	}
	f -= d
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func mathConv2(val any, args []any) (float64, float64, error, bool) {
	if len(args) == 0 {
		return 0, 0, ErrModPoorArgs, false
	}
	f, ok := floatConv(val)
	if !ok {
		return 0, 0, nil, false
	}
	d, ok := floatConv(args[0])
	if !ok {
		return 0, 0, nil, false
	}
	return f, d, nil, true
}

func floatConv(val any) (f float64, ok bool) {
	if val == nil {
		return 0, false
	}
	ok = true
	switch x := val.(type) {
	case int:
		f = float64(x)
	case int8:
		f = float64(x)
	case int16:
		f = float64(x)
	case int32:
		f = float64(x)
	case int64:
		f = float64(x)
	case uint:
		f = float64(x)
	case uint8:
		f = float64(x)
	case uint16:
		f = float64(x)
	case uint32:
		f = float64(x)
	case uint64:
		f = float64(x)
	case float32:
		f = float64(x)
	case *int:
		f = float64(*x)
	case *int8:
		f = float64(*x)
	case *int16:
		f = float64(*x)
	case *int32:
		f = float64(*x)
	case *int64:
		f = float64(*x)
	case *uint:
		f = float64(*x)
	case *uint8:
		f = float64(*x)
	case *uint16:
		f = float64(*x)
	case *uint32:
		f = float64(*x)
	case *uint64:
		f = float64(*x)
	case *float32:
		f = float64(*x)
	case float64:
		f = x
	case *float64:
		f = *x
	case string:
		f1, err := strconv.ParseFloat(x, 64)
		if ok = err == nil; ok {
			f = f1
		}
	case *string:
		f1, err := strconv.ParseFloat(*x, 64)
		if ok = err == nil; ok {
			f = f1
		}
	case []byte:
		f1, err := strconv.ParseFloat(byteconv.B2S(x), 64)
		if ok = err == nil; ok {
			f = f1
		}
	case *[]byte:
		f1, err := strconv.ParseFloat(byteconv.B2S(*x), 64)
		if ok = err == nil; ok {
			f = f1
		}
	default:
		ok = false
		return
	}
	return
}
