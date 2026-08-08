package dyntpl

import (
	"encoding/binary"
	"math"

	"github.com/koykov/byteconv"
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

func modBin(ctx *Ctx, buf *any, val any, args []any) error {
	var a any
	switch {
	case val != nil:
		a = val
	case len(args) > 0:
		a = args[0]
		args = args[1:]
	default:
		return ErrModNoArgs
	}

	var order binary.ByteOrder
	order = binary.LittleEndian
	if len(args) > 1 {
		var orderRaw string
		if s, ok := ConvStr(args[1]); ok {
			orderRaw = s
		} else if b, ok := ConvBytes(args[1]); ok {
			orderRaw = byteconv.B2S(b)
		}
		if len(orderRaw) > 0 {
			switch orderRaw {
			case "be", "BE", "big_endian", "bigEndian", "BigEndian":
				order = binary.BigEndian
			case "le", "LE", "little_endian", "littleEndian", "LittleEndian":
				fallthrough
			default:
				order = binary.LittleEndian
			}
		}
	}

	ctx.BufAcc.StakeOut().WriteBinary(order, a)
	ctx.BufModOut(buf, ctx.BufAcc.StakedBytes())

	return nil
}
