package dyntpl

import (
	"math"

	"github.com/koykov/any2bytes"
	"github.com/koykov/fastconv"
)

const (
	// Types of round.
	round = iota
	roundPrec
	ceil
	ceilPrec
	floor
	floorPrec

	// Hex digits in upper case.
	hexUp = "0123456789ABCDEF"
)

// If var is empty, the given default value (first in args) will print instead.
func modDefault(_ *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if len(args) == 0 {
		err = ErrModNoArgs
		return
	}
	// Consecutive try to assert value to known (builtin) types:
	// * int
	// * uint
	// * float
	// * bytes
	// * string
	// * bool
	// ... and check if value is empty.
	if i, ok := ConvInt(val); ok {
		if i == 0 {
			*buf = args[0]
			return
		}
	}
	if u, ok := ConvUint(val); ok {
		if u == 0 {
			*buf = args[0]
			return
		}
	}
	if f, ok := ConvFloat(val); ok {
		if f == 0 {
			*buf = args[0]
			return
		}
	}
	if b, ok := ConvBytes(val); ok {
		if len(b) == 0 {
			*buf = args[0]
			return
		}
	}
	if s, ok := ConvStr(val); ok {
		if len(s) == 0 {
			*buf = args[0]
			return
		}
	}
	if b, ok := ConvBool(val); ok {
		if !b {
			*buf = args[0]
			return
		}
	}
	return nil
}

// Shorthand replacement of {% if ... %}{%= ... %}{% endif %} statement.
//
// Example of usage: {%= leftVal|ifThen(val) %}, leftVal should be a boolean.
func modIfThen(_ *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if len(args) == 0 {
		err = ErrModNoArgs
		return
	}
	if b, ok := ConvBool(val); ok {
		if b {
			*buf = args[0]
		}
	}
	return
}

// Shorthand replacement of {% if ... %}{%= ... %}{% else %}{%= ... %}{% endif %} statement.
//
// Example of usage: {%= leftVal|ifThenElse(valIfTrue, valIfFalse) %}, leftVal should be a boolean.
// valIfTrue and valIfFalse may has arbitrary types or may be a static values.
func modIfThenElse(_ *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if len(args) < 2 {
		err = ErrModPoorArgs
		return
	}
	if b, ok := ConvBool(val); ok {
		if b {
			*buf = args[0]
		} else {
			*buf = args[1]
		}
	}
	return
}

// Round float val to integer using rounding half away from zero algorithm.
func modRound(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if f, ok := ConvFloat(val); ok {
		ctx.BufF = roundHelper(f, round, args)
		*buf = &ctx.BufF
	}
	return
}

// Round to precision, example: pi|roundPrec(3) will print 3.141
func modRoundPrec(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if f, ok := ConvFloat(val); ok {
		ctx.BufF = roundHelper(f, roundPrec, args)
		*buf = &ctx.BufF
	}
	return
}

// Round to least integer value greater than or equal to val.
func modCeil(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if f, ok := ConvFloat(val); ok {
		ctx.BufF = roundHelper(f, ceil, args)
		*buf = &ctx.BufF
	}
	return
}

// Ceil round to precision, example: 56.68734|ceilPrec will print 56.688
func modCeilPrec(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if f, ok := ConvFloat(val); ok {
		ctx.BufF = roundHelper(f, ceilPrec, args)
		*buf = &ctx.BufF
	}
	return
}

// Round to greatest integer value less than or equal to val.
func modFloor(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if f, ok := ConvFloat(val); ok {
		ctx.BufF = roundHelper(f, floor, args)
		*buf = &ctx.BufF
	}
	return
}

// Float round to precision, example: 20.214999|floorPrec(3) will print 20.214
func modFloorPrec(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if f, ok := ConvFloat(val); ok {
		ctx.BufF = roundHelper(f, floorPrec, args)
		*buf = &ctx.BufF
	}
	return
}

// Universal internal round helper for round modifiers.
func roundHelper(f float64, mode int, args []interface{}) float64 {
	var (
		prec int64
		ok   bool
	)
	if len(args) > 0 {
		if prec, ok = if2int(args[0]); !ok {
			return f
		}
	}
	switch mode {
	case round:
		return math.Round(f)
	case roundPrec:
		if prec == 0 {
			return f
		}
		p := math.Pow10(int(prec))
		return float64(int(f*p)) / p
	case ceil:
		return math.Ceil(f)
	case ceilPrec:
		if prec == 0 {
			return f
		}
		p := math.Pow10(int(prec))
		x := p * f
		return math.Ceil(x) / p
	case floor:
		return math.Floor(f)
	case floorPrec:
		if prec == 0 {
			return f
		}
		p := math.Pow10(int(prec))
		x := p * f
		return math.Floor(x) / p
	}
	return f
}

// URL encode string value.
//
// see https://golang.org/src/net/url/url.go#L100
func modUrlEncode(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	// Get count of encode iterations (cases: uu=, uuu=, ...).
	itr := printIterations(args)

	// Get the source.
	if p, ok := ConvBytes(val); ok {
		ctx.buf = append(ctx.buf[:0], p...)
	} else if s, ok := ConvStr(val); ok {
		ctx.buf = append(ctx.buf[:0], s...)
	} else if ctx.Buf2, err = any2bytes.AnyToBytes(ctx.Buf2, val); err == nil {
		ctx.buf = append(ctx.buf[:0], ctx.Buf2...)
	} else {
		return ErrModNoStr
	}
	l := len(ctx.buf)
	if l == 0 {
		return ErrModEmptyStr
	}
	for c := 0; c < itr; c++ {
		ctx.Buf.Reset()
		_ = ctx.buf[l-1]
		for i := 0; i < l; i++ {
			if ctx.buf[i] >= 'a' && ctx.buf[i] <= 'z' || ctx.buf[i] >= 'A' && ctx.buf[i] <= 'Z' ||
				ctx.buf[i] >= '0' && ctx.buf[i] <= '9' || ctx.buf[i] == '-' || ctx.buf[i] == '.' || ctx.buf[i] == '_' {
				ctx.Buf.WriteByte(ctx.buf[i])
			} else if ctx.buf[i] == ' ' {
				ctx.Buf.WriteByte('+')
			} else {
				ctx.Buf.WriteByte('%').WriteByte(hexUp[ctx.buf[i]>>4]).WriteByte(hexUp[ctx.buf[i]&15])
			}
		}
		ctx.buf = append(ctx.buf[:0], ctx.Buf...)
		l = ctx.Buf.Len()
	}
	*buf = &ctx.Buf
	return
}

func modInclude(ctx *Ctx, buf *interface{}, _ interface{}, args []interface{}) (err error) {
	if len(args) < 0 {
		err = ErrModPoorArgs
		return
	}
	// Get tpl identifier and fallback id.
	var id, fbID string
	if idRaw, ok := args[0].(*[]byte); ok {
		id = fastconv.B2S(*idRaw)
	}
	if len(args) > 1 {
		if idRaw, ok := args[1].(*[]byte); ok {
			fbID = fastconv.B2S(*idRaw)
		}
	}
	// Try to render template using the same context object.
	w := ctx.getW()
	if err = RenderFbTo(w, id, fbID, ctx); err != nil {
		return
	}
	ctx.buf = append(ctx.buf[:0], w.Bytes()...)
	*buf = &ctx.buf
	return
}

func modTestNameOf(_ *Ctx, _ *interface{}, _ interface{}, _ []interface{}) (err error) {
	// do nothing
	return
}
