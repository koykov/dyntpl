package dyntpl

import (
	"math"
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

func modMathMul(ctx *Ctx, buf *any, val any, args []any) (err error) {
	var (
		f, d float64
		ok   bool
	)
	if f, d, err, ok = mathConv2(val, args); !ok {
		return
	}
	f *= d
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func modMathDiv(ctx *Ctx, buf *any, val any, args []any) (err error) {
	var (
		f, d float64
		ok   bool
	)
	if f, d, err, ok = mathConv2(val, args); !ok {
		return
	}
	f /= d
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func modMathMod(ctx *Ctx, buf *any, val any, args []any) (err error) {
	var (
		f, d float64
		ok   bool
	)
	if f, d, err, ok = mathConv2(val, args); !ok {
		return
	}
	f1, d1 := int64(f), int64(d)
	f1 = f1 % d1
	ctx.BufI = f1
	*buf = &ctx.BufI
	return
}

func modMathSqrt(ctx *Ctx, buf *any, val any, _ []any) (err error) {
	f, ok := floatConv(val)
	if !ok {
		return
	}
	f = math.Sqrt(f)
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func modMathCbrt(ctx *Ctx, buf *any, val any, _ []any) (err error) {
	f, ok := floatConv(val)
	if !ok {
		return
	}
	f = math.Cbrt(f)
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func modMathRadical(ctx *Ctx, buf *any, val any, args []any) (err error) {
	var (
		f, d, eps float64
		ok        bool
	)
	if f, d, err, ok = mathConv2(val, args); !ok {
		return
	}
	eps = .000001
	if len(args) > 2 {
		if e, ok := floatConv(args[1]); ok {
			eps = e
		}
	}

	// Newton
	root := f / d
	rn := f
	for math.Abs(root-rn) >= eps {
		rn = f
		for i := 1; i < int(d); i++ {
			rn = rn / root
		}
		root = .5 * (rn + root)
	}
	ctx.BufF = root
	*buf = &ctx.BufF
	return
}

func modMathExp(ctx *Ctx, buf *any, val any, _ []any) (err error) {
	f, ok := floatConv(val)
	if !ok {
		return
	}
	f = math.Exp(f)
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func modMathLog(ctx *Ctx, buf *any, val any, _ []any) (err error) {
	f, ok := floatConv(val)
	if !ok {
		return
	}
	f = math.Log(f)
	ctx.BufF = f
	*buf = &ctx.BufF
	return
}

func modMathFact(ctx *Ctx, buf *any, val any, args []any) (err error) {
	var (
		f, d float64
		ok   bool
	)
	if f, d, err, ok = mathConv2(val, args); !ok {
		return
	}
	r := f
	for i := 1; i < int(d); i++ {
		r *= f
	}
	ctx.BufF = r
	*buf = &ctx.BufF
	return
}

func modMathMax(ctx *Ctx, buf *any, _ any, args []any) (err error) {
	var f, d float64
	if f, d, err = mathConvArgs2(args); err != nil {
		return
	}
	ctx.BufF = math.Max(f, d)
	*buf = &ctx.BufF
	return
}

func modMathMin(ctx *Ctx, buf *any, _ any, args []any) (err error) {
	var f, d float64
	if f, d, err = mathConvArgs2(args); err != nil {
		return
	}
	ctx.BufF = math.Min(f, d)
	*buf = &ctx.BufF
	return
}

func modMathPow(ctx *Ctx, buf *any, val any, args []any) (err error) {
	var (
		f, d float64
		ok   bool
	)
	if f, d, err, ok = mathConv2(val, args); !ok {
		return
	}
	ctx.BufF = math.Pow(f, d)
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

func mathConvArgs2(args []any) (float64, float64, error) {
	if len(args) < 2 {
		return 0, 0, ErrModPoorArgs
	}
	d, ok := floatConv(args[0])
	if !ok {
		return 0, 0, nil
	}
	f, ok := floatConv(args[1])
	if !ok {
		return 0, 0, nil
	}
	return f, d, nil
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
