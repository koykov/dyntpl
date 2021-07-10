package dyntpl

import (
	"math"
)

const (
	// Types of round.
	round = iota
	roundPrec
	ceil
	ceilPrec
	floor
	floorPrec
)

// If var is empty, the given default value (first in args) will print instead.
func modDefault(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) (err error) {
	if len(args) == 0 {
		err = ErrModNoArgs
		return
	}
	// Check is value is empty.
	if EmptyCheck(ctx, val) {
		*buf = args[0]
	}
	return
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

func modTestNameOf(_ *Ctx, _ *interface{}, _ interface{}, _ []interface{}) (err error) {
	// do nothing
	return
}
