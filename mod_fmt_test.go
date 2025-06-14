package dyntpl

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/koykov/bytebuf"
)

var (
	NaN    = math.NaN()
	posInf = math.Inf(1)
	negInf = math.Inf(-1)

	intVar = 0

	array  = [5]int{1, 2, 3, 4, 5}
	iarray = [4]any{1, "hello", 2.5, nil}
	slice  = array[:]
	islice = iarray[:]

	barray = [5]uint8{1, 2, 3, 4, 5}
	bslice = barray[:]

	fmtStages = []any{
		12345,
		12345,
		true,

		// basic string
		"abc",
		"abc",
		"abc",
		"\xff\xf0\x0f\xff",
		"\xff\xf0\x0f\xff",
		"",
		"",
		"",
		"",
		"xyz",
		"xyz",
		"xyz",
		"xyz",
		"xyz",
		"xyz",
		"xyz",
		"xyz",

		// basic bytes
		[]byte("abc"),
		[3]byte{'a', 'b', 'c'},
		&[3]byte{'a', 'b', 'c'},
		[]byte("abc"),
		[]byte("abc"),
		[]byte("\xff\xf0\x0f\xff"),
		[]byte("\xff\xf0\x0f\xff"),
		[]byte(""),
		[]byte(""),
		[]byte(""),
		[]byte(""),
		[]byte("xyz"),
		[]byte("xyz"),
		[]byte("xyz"),
		[]byte("xyz"),
		[]byte("xyz"),
		[]byte("xyz"),
		[]byte("xyz"),
		[]byte("xyz"),

		// escaped strings
		"",
		"",
		"\"",
		"\"",
		"`",
		"`",
		"\n",
		"\n",
		`\n`,
		`\n`,
		"abc",
		"abc",
		"日本語",
		"日本語",
		"日本語",
		"日本語",
		"\a\b\f\n\r\t\v\"\\",
		"\a\b\f\n\r\t\v\"\\",
		"\a\b\f\n\r\t\v\"\\",
		"\a\b\f\n\r\t\v\"\\",
		"☺",
		"☺", // The space modifier should have no effect.
		"☺",
		"☺",
		"☺",
		"⌘",
		"⌘",
		"⌘",
		"⌘",
		"⌘",
		"⌘",
		"⌘", // 0 has no effect when - is present.
		"⌘",
		"\n",
		"\r",
		"\t",
		"\b",
		"abc\xffdef",
		"abc\xffdef",
		"abc\xffdef",
		"abc\xffdef",
		// Runes that are not printable.
		"\U0010ffff",
		"\U0010ffff",
		"\U0010ffff",
		"\U0010ffff",
		// Runes that are not valid.
		string(rune(0x110000)),
		string(rune(0x110000)),
		string(rune(0x110000)),
		string(rune(0x110000)),

		// characters
		uint('x'),
		0xe4,
		0x672c,
		'日',
		'⌘', // Specifying precision should have no effect.
		'⌘',
		'⌘',
		uint64(0x100000000),
		// Runes that are not printable.
		'\U00000e00',
		'\U0010ffff',
		// Runes that are not valid.
		-1,
		0xDC80,
		rune(0x110000),
		int64(0xFFFFFFFFF),
		uint64(0xFFFFFFFFF),

		// escaped characters
		uint(0),
		uint(0),
		'"',
		'"',
		'\'',
		'\'',
		'`',
		'`',
		'x',
		'x',
		'ÿ',
		'ÿ',
		'\n',
		'\n',
		'☺',
		'☺',
		'☺', // The space modifier should have no effect.
		'☺', // Specifying precision should have no effect.
		'⌘',
		'⌘',
		'⌘',
		'⌘',
		'⌘',
		'⌘',
		'⌘', // 0 has no effect when - is present.
		'⌘',
		// Runes that are not printable.
		'\U00000e00',
		'\U0010ffff',
		// Runes that are not valid.
		int32(-1),
		0xDC80,
		rune(0x110000),
		int64(0xFFFFFFFFF),
		uint64(0xFFFFFFFFF),

		// width
		"abc",
		[]byte("abc"),
		"\u263a",
		[]byte("\u263a"),
		"abc",
		[]byte("abc"),
		"abc",
		[]byte("abc"),
		"abcdefghijklmnopqrstuvwxyz",
		[]byte("abcdefghijklmnopqrstuvwxyz"),
		"abcdefghijklmnopqrstuvwxyz",
		[]byte("abcdefghijklmnopqrstuvwxyz"),
		"日本語日本語",
		[]byte("日本語日本語"),
		"日本語日本語",
		[]byte("日本語日本語"),
		"日本語日本語",
		[]byte("日本語日本語"),
		"abc",
		[]byte("abc"),
		"abc",
		[]byte("abc"),
		"abcdefghijklmnopqrstuvwxyz",
		[]byte("abcdefghijklmnopqrstuvwxyz"),
		"abcdefghijklmnopqrstuvwxyz",
		[]byte("abcdefghijklmnopqrstuvwxyz"),
		"日本語日本語",
		[]byte("日本語日本語"),
		"日本語",
		[]byte("日本語"),
		"日本語",
		[]byte("日本語"),
		"日本語日本語",
		[]byte("日本語日本語"),
		nil,
		nil,

		// integers
		uint(12345),
		int(-12345),
		^uint8(0),
		^uint16(0),
		^uint32(0),
		^uint64(0),
		int8(-1 << 7),
		int16(-1 << 15),
		int32(-1 << 31),
		int64(-1 << 63),
		0,
		0,
		0,
		0,
		12345,
		12345,
		-12345,
		7,
		-6,
		7,
		-6,
		^uint32(0),
		^uint64(0),
		int64(-1 << 63),
		01234,
		-01234,
		01234,
		-01234,
		01234,
		-01234,
		^uint32(0),
		^uint64(0),
		0,
		0x12abcdef,
		0x12abcdef,
		^uint32(0),
		^uint64(0),
		7,
		12345,
		-12345,
		12345,
		12345,
		-12345,
		1234,
		-1234,
		1234,
		-1234,
		1234,
		-1234,
		0x1234abc,
		0x1234abc,
		01234,

		// Test correct f.intbuf overflow checks.
		1,
		-1,
		42,
		-42,
		42,
		42,
		42,

		// unicode format
		0,
		-1,
		'\n',
		'\n',
		'x', // Plus flag should have no effect.
		'x', // Space flag should have no effect.
		'x', // Precisions below 4 should print 4 digits.
		'\u263a',
		'\u263a',
		'\U0001D6C2',
		'\U0001D6C2',
		'⌘',
		'⌘',
		'⌘',
		'⌘',
		uint(42),
		'日',

		// floats
		0.0,
		1.0,
		0.0,
		1.0,
		-1.0,
		-1.0,
		float32(-1.0),
		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,
		+1.0,
		-1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		0.0,
		1.0,
		-1.0,
		-1.0,
		1.0,
		float32(1.0),
		1.0,
		// Test sharp flag used with floats.
		1e-323,
		-1.0,
		1.1,
		123456.0,
		1234567.0,
		1230000.0,
		1000000.0,
		1.0,
		1.0,
		1.0,
		1.0,
		1100000.0,
		1.0,
		1.0,
		1.0,
		1.0,
		100000.0,
		1.234,
		0.1234,
		1.23,
		0.123,
		1.2,
		0.12,
		10.2,
		0.0,
		0.012,
		123.0,
		123.0,
		123.0,
		123.0,
		123.0,
		123.0,
		123.0,
		123.0,
		123000.0,
		1.0,
		// The sharp flag has no effect for binary float format.
		1.0,
		// Precision has no effect for binary float format.
		float32(1.0),
		-1.0,
		// Test correct f.intbuf boundary checks.
		1.0,
		-1.0,
		// float infinites and NaNs
		posInf,
		negInf,
		NaN,
		posInf,
		posInf,
		negInf,
		negInf,
		negInf,
		negInf,
		negInf,
		posInf,
		NaN,
		NaN,
		NaN,
		NaN,
		NaN,
		NaN,
		// Zero padding does not apply to infinities and NaN.
		posInf,
		posInf,
		negInf,
		NaN,
		NaN,

		1.0,
		1234.5678e3,
		1234.5678e-8,
		-7.0,
		-1e-9,
		1234.5678e3,
		1234.5678e-8,
		-7.0,
		-1e-9,
		1234.5678e3,
		float32(1234.5678e3),
		1234.5678e-8,
		-7.0,
		-1e-9,
		float32(-1e-9),
		1.0,
		1234.5678e3,
		1234.5678e-8,
		-7.0,
		-1e-9,
		1234.5678e3,
		float32(1234.5678e3),
		1234.5678e-8,
		-7.0,
		-1e-9,
		float32(-1e-9),
		"qwertyuiop",
		"qwertyuiop",
		"qwertyuiop",
		'x',
		'x',
		1.2345e3,
		1.2345e-3,
		1.2345e3,
		1.2345e-3,
		1.2345e3,
		1.23456789e3,
		1.23456789e-3,
		12345678901.23456789,
		1.23456789e3,
		1.23456789e3,
		1.23456789e-3,
		1.23456789e3,
		1.23456789e-3,
		1.23456789e20,

		// arrays
		array,
		iarray,
		barray,
		&array,
		&iarray,
		&barray,

		// slices
		slice,
		islice,
		bslice,
		&slice,
		&islice,
		&bslice,

		// byte arrays and slices with %b,%c,%d,%o,%U and %v
		[3]byte{65, 66, 67},
		[3]byte{65, 66, 67},
		[3]byte{65, 66, 67},
		[3]byte{65, 66, 67},
		[3]byte{65, 66, 67},
		[3]byte{65, 66, 67},
		[1]byte{123},
		[]byte{},
		[]byte{},
		[]byte{1, 11, 111},
		[]byte{1, 11, 111},
		[]byte{1, 11, 111},
		[]byte{1, 11, 111},
		[]byte{1, 11, 111},
		[]byte{1, 11, 111},
		[]byte{1, 11, 111},
		[]byte{1, 11, 111},
		[]byte{1, 11, 111},
		// f.space should and f.plus should not have an effect with %v.
		[]byte{1, 11, 111},
		[3]byte{1, 11, 111},
		[]byte{1, 11, 111},
		[3]byte{1, 11, 111},
		// f.space and f.plus should have an effect with %d.
		[]byte{1, 11, 111},
		[3]byte{1, 11, 111},
		[]byte{1, 11, 111},
		[3]byte{1, 11, 111},

		// floates with %v
		1.2345678,
		float32(1.2345678),

		// go syntax
		new(byte),
		make(chan int),
		uint64(1<<64 - 1),
		1000000000,
		map[string]int{"a": 1},
		[]string{"a", "b"},
		[]int(nil),
		[]int{},
		array,
		&array,
		iarray,
		&iarray,
		map[int]byte(nil),
		map[int]byte{},
		"foo",
		barray,
		bslice,
		[]int32(nil),
		1.2345678,
		float32(1.2345678),

		// Whole number floats are printed without decimals. See Issue 27634.
		1.0,
		1000000.0,
		float32(1.0),
		float32(1000000.0),

		// Only print []byte and []uint8 as type []byte if they appear at the top level.
		[]byte(nil),
		[]uint8(nil),
		[]byte{},
		[]uint8{},
		reflect.ValueOf([]byte{}),
		reflect.ValueOf([]uint8{}),
		&[]byte{},
		&[]byte{},
		[3]byte{},
		[3]uint8{},

		// slices with other formats
		[]int{1, 2, 15},
		[]int{1, 2, 15},
		[]int{1, 2, 15},
		[]byte{1, 2, 15},
		[]string{"a", "b"},
		[]byte{1},
		[]byte{1, 2, 3},

		// Padding with byte slices.
		[]byte{},
		[]byte{},
		[]byte{},
		[]byte{},
		[]byte{},
		[]byte{},
		[]byte{0xab},
		[]byte{0xab},
		[]byte{0xab},
		[]byte{0xab},
		[]byte{0xab},
		[]byte{0xab},
		[]byte{0xab},
		[]byte{0xab},
		[]byte{0xab, 0xcd},
		[]byte{0xab, 0xcd},
		[]byte{0xab, 0xcd},
		[]byte{0xab, 0xcd},
		[]byte{0xab, 0xcd},
		[]byte{0xab, 0xcd},
		[]byte{0xab, 0xcd},
		[]byte{0xab, 0xcd},
		[]byte{0xab},
		[]byte{0xab},
		[]byte{0xab, 0xcd},
		[]byte{0xab, 0xcd},
		// Same for strings
		"",
		"",
		"",
		"",
		"",
		"",
		"\xab",
		"\xab",
		"\xab",
		"\xab",
		"\xab",
		"\xab",
		"\xab",
		"\xab",
		"\xab\xcd",
		"\xab\xcd",
		"\xab\xcd",
		"\xab\xcd",
		"\xab\xcd",
		"\xab\xcd",
		"\xab\xcd",
		"\xab\xcd",
		"\xab",
		"\xab",
		"\xab\xcd",
		"\xab\xcd",

		// %T
		byte(0),
		intVar,
		&intVar,
		nil,
		nil,

		// %p with pointers
		(*int)(nil),
		(*int)(nil),
		&intVar,
		&intVar,
		&array,
		&slice,
		(*int)(nil),
		&intVar,
		// %p on non-pointers
		make(chan int),
		make(map[int]int),
		func() {},
		27,  // not a pointer at all
		nil, // nil on its own has no type ...
		nil, // ... and hence is not a pointer type.
		// pointers with specified base
		&intVar,
		&intVar,
		&intVar,
		&intVar,
		&intVar,
		// %v on pointers
		nil,
		nil,
		(*int)(nil),
		(*int)(nil),
		&intVar,
		&intVar,
		(*int)(nil),
		&intVar,

		// %d on Stringer should give integer if possible
		time.Time{}.Month(),
		time.Time{}.Month(),

		// erroneous things
		nil,
		2,
		"hello",
		"hello",
		"hello",
		0,
		0,
		// Extra argument errors should format without flags set.
		"12345",

		// Test that maps with non-reflexive keys print all keys and values.
		map[float64]int{NaN: 1, NaN: 1},

		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,
		1.0,
		-1.0,

		// Use spaces instead of zero if padding to the right.
		"abc",
		1.0,

		// integer formatting should not alter padding for other elements.
		[]any{1, 2.0, "x"},
		[]any{0, 2.0, "x"},
	}
)

func TestModFmt(t *testing.T) {
	for i := 0; i < len(fmtStages); i++ {
		t.Run(fmt.Sprintf("fmt%d", i), func(t *testing.T) { testModWA(t, modArgs{"fmtVar": fmtStages[i]}) })
	}
}

func BenchmarkModFmt(b *testing.B) {
	args := modArgs{"fmtVar": fmtStages[0]}
	var bb bytebuf.Chain
	b.ReportAllocs()
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		key := bb.Reset().WriteString("fmt").WriteInt(int64(j % len(fmtStages))).String()
		st := getStage(key)
		if st == nil {
			b.Error("stage not found")
			return
		}

		ctx := AcquireCtx()
		args["fmtVar"] = fmtStages[j%len(fmtStages)]
		for k, v := range args {
			ctx.SetStatic(k, v)
		}
		buf.Reset()
		err := Write(&buf, key, ctx)
		if err != nil {
			b.Error(err)
		}
		ReleaseCtx(ctx)
	}
}
