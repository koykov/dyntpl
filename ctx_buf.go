package dyntpl

import (
	"github.com/koykov/any2bytes"
	"github.com/koykov/fastconv"
)

// Primitive byte buffer to use in context objects and modifiers/helpers.
type ByteBuf []byte

// Conversion to bytes function.
func ByteBufToBytes(dst []byte, val interface{}) ([]byte, error) {
	if b, ok := val.(*ByteBuf); ok {
		dst = append(dst, *b...)
		return dst, nil
	}
	return dst, any2bytes.ErrUnknownType
}

// Reset length of the buffer.
func (b *ByteBuf) Reset() *ByteBuf {
	*b = (*b)[:0]
	return b
}

// Get contents of the buffer.
func (b *ByteBuf) Bytes() []byte {
	return *b
}

// Get contents of the buffer as string.
func (b *ByteBuf) String() string {
	return fastconv.B2S(*b)
}

// Write bytes to the buffer.
func (b *ByteBuf) Write(p []byte) *ByteBuf {
	*b = append(*b, p...)
	return b
}

// Write single byte.
//
// Use suffix "B" instead of "Byte" to avoid warning of incorrect naming of canonical methods:
// > Method 'WriteByte(byte)' should have signature 'WriteByte(byte) error'
// For more info see https://github.com/golang/tools/blob/master/go/analysis/passes/stdmethods/stdmethods.go#L63
func (b *ByteBuf) WriteB(p byte) *ByteBuf {
	*b = append(*b, p)
	return b
}

// Write string to the buffer.
func (b *ByteBuf) WriteStr(s string) *ByteBuf {
	*b = append(*b, s...)
	return b
}

// Write integer value to the buffer.
func (b *ByteBuf) WriteInt(i int64) *ByteBuf {
	*b, _ = any2bytes.IntToBytes(*b, i)
	return b
}

// Write float value to the boffer.
func (b *ByteBuf) WriteFloat(f float64) *ByteBuf {
	*b, _ = any2bytes.FloatToBytes(*b, f)
	return b
}

// Write boolean value to the buffer.
func (b *ByteBuf) WriteBool(v bool) *ByteBuf {
	*b, _ = any2bytes.BoolToBytes(*b, v)
	return b
}

// Get length of the buffer.
func (b *ByteBuf) Len() int {
	return len(*b)
}

// Get capacity of the buffer.
func (b *ByteBuf) Cap() int {
	return cap(*b)
}
