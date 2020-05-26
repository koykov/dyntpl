package dyntpl

import "github.com/koykov/cbytealg"

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

func (b *ByteBuf) Write(p []byte) {
	*b = append(*b, p...)
}

// Use suffix "B" instead of "Byte" to avoid warning of incorrect naming of canonical methods:
// > Method 'WriteByte(byte)' should have signature 'WriteByte(byte) error'
// For more info see https://github.com/golang/tools/blob/master/go/analysis/passes/stdmethods/stdmethods.go#L63
func (b *ByteBuf) WriteB(p byte) {
	*b = append(*b, p)
}

func (b *ByteBuf) WriteStr(s string) {
	*b = append(*b, s...)
}

func (b *ByteBuf) WriteInt(i int64) {
	*b, _ = cbytealg.IntToBytes(*b, i)
}

func (b *ByteBuf) WriteFloat(f float64) {
	*b, _ = cbytealg.FloatToBytes(*b, f)
}

func (b *ByteBuf) WriteBool(v bool) {
	*b, _ = cbytealg.BoolToBytes(*b, v)
}
