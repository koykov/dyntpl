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

func (b *ByteBuf) Reset() {
	*b = (*b)[:0]
}

func (b *ByteBuf) Write(p []byte) {
	*b = append(*b, p...)
}

func (b *ByteBuf) WriteByte(p byte) {
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

func (b *ByteBuf) ResetWrite(p []byte) {
	*b = append((*b)[:0], p...)
}

func (b *ByteBuf) ResetWriteByte(p byte) {
	*b = append((*b)[:0], p)
}

func (b *ByteBuf) ResetWriteStr(s string) {
	*b = append((*b)[:0], s...)
}

func (b *ByteBuf) ResetWriteInt(i int64) {
	*b, _ = cbytealg.IntToBytes((*b)[:0], i)
}

func (b *ByteBuf) ResetWriteFloat(f float64) {
	*b, _ = cbytealg.FloatToBytes((*b)[:0], f)
}

func (b *ByteBuf) ResetWriteBool(v bool) {
	*b, _ = cbytealg.BoolToBytes((*b)[:0], v)
}
