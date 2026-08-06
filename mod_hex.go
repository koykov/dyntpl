package dyntpl

import (
	"encoding/binary"
	"math"
)

func modHex(ctx *Ctx, buf *any, val any, args []any) error {
	var a any
	switch {
	case val != nil:
		a = val
	case len(args) > 0:
		a = args[0]
	default:
		return ErrModNoArgs
	}

	ctx.BufAcc.StakeOut()
	if b, ok := ConvBytes(a); ok {
		ctx.BufAcc.WriteHex(b)
	} else if bb, ok := ConvBytesSlice(a); ok {
		for i := 0; i < len(bb); i++ {
			ctx.BufAcc.WriteHex(bb[i])
		}
	} else if s, ok := ConvStr(a); ok {
		ctx.BufAcc.WriteHexString(s)
	} else if ss, ok := ConvStrSlice(a); ok {
		for i := 0; i < len(ss); i++ {
			ctx.BufAcc.WriteHexString(ss[i])
		}
	} else if i, ok := ConvInt(a); ok {
		ctx.BufAcc.WriteIntBase(i, 16)
	} else if u, ok := ConvUint(a); ok {
		ctx.BufAcc.WriteUintBase(u, 16)
	} else if f, ok := ConvFloat(a); ok {
		bits := math.Float64bits(f)
		ctx.BufAcc.WriteUintBase(bits, 16)
	} else {
		raw := ctx.BufAcc.WriteBinary(binary.LittleEndian, a).StakedBytes()
		ctx.BufAcc.StakeOut().WriteHex(raw)
	}
	ctx.BufModOut(buf, ctx.BufAcc.StakedBytes())

	return nil
}
