package dyntpl

import (
	"github.com/koykov/cbytealg"
	"github.com/koykov/fastconv"
)

type ByteBuf []byte

func ByteBufToBytes(dst []byte, val interface{}) ([]byte, error) {
	if b, ok := val.(*ByteBuf); ok {
		dst = append(dst, *b...)
		return dst, nil
	}
	return dst, cbytealg.ErrUnknownType
}

func (b *ByteBuf) Reset() *ByteBuf {
	*b = (*b)[:0]
	return b
}

func (b *ByteBuf) Bytes() []byte {
	return *b
}

func (b *ByteBuf) String() string {
	return fastconv.B2S(*b)
}

func (b *ByteBuf) Write(p []byte) *ByteBuf {
	*b = append(*b, p...)
	return b
}

// Use suffix "B" instead of "Byte" to avoid warning of incorrect naming of canonical methods:
// > Method 'WriteByte(byte)' should have signature 'WriteByte(byte) error'
// For more info see https://github.com/golang/tools/blob/master/go/analysis/passes/stdmethods/stdmethods.go#L63
func (b *ByteBuf) WriteB(p byte) *ByteBuf {
	*b = append(*b, p)
	return b
}

func (b *ByteBuf) WriteStr(s string) *ByteBuf {
	*b = append(*b, s...)
	return b
}

func (b *ByteBuf) WriteInt(i int64) *ByteBuf {
	*b, _ = cbytealg.IntToBytes(*b, i)
	return b
}

func (b *ByteBuf) WriteFloat(f float64) *ByteBuf {
	*b, _ = cbytealg.FloatToBytes(*b, f)
	return b
}

func (b *ByteBuf) WriteBool(v bool) *ByteBuf {
	*b, _ = cbytealg.BoolToBytes(*b, v)
	return b
}
