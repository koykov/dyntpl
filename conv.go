package dyntpl

// Collection of conversion functions.

import (
	"strconv"

	"github.com/koykov/bytebuf"
	"github.com/koykov/byteconv"
)

type intConverter interface {
	Int() (int64, error)
}

// ConvInt tries to convert value to integer.
func ConvInt(val any) (i int64, ok bool) {
	ok = true
	switch x := val.(type) {
	case int:
		i = int64(x)
	case *int:
		i = int64(*x)
	case int8:
		i = int64(x)
	case *int8:
		i = int64(*x)
	case int16:
		i = int64(x)
	case *int16:
		i = int64(*x)
	case int32:
		i = int64(x)
	case *int32:
		i = int64(*x)
	case int64:
		i = x
	case *int64:
		i = *x
	default:
		ok = false
	}
	return
}

// ConvUint tries to convert value to uint.
func ConvUint(val any) (u uint64, ok bool) {
	ok = true
	switch x := val.(type) {
	case uint:
		u = uint64(x)
	case *uint:
		u = uint64(*x)
	case uint8:
		u = uint64(x)
	case *uint8:
		u = uint64(*x)
	case uint16:
		u = uint64(x)
	case *uint16:
		u = uint64(*x)
	case uint32:
		u = uint64(x)
	case *uint32:
		u = uint64(*x)
	case uint64:
		u = x
	case *uint64:
		u = *x
	default:
		ok = false
	}
	return
}

// ConvFloat tries to convert value to float.
func ConvFloat(val any) (f float64, ok bool) {
	ok = true
	switch x := val.(type) {
	case float32:
		f = float64(x)
	case *float32:
		f = float64(*x)
	case float64:
		f = x
	case *float64:
		f = *x
	default:
		ok = false
	}
	return
}

// ConvBytes tries to convert value to bytes.
func ConvBytes(val any) (b []byte, ok bool) {
	ok = true
	switch x := val.(type) {
	case []byte:
		b = x
	case *[]byte:
		b = *x
	case bytebuf.Chain:
		b = x
	case *bytebuf.Chain:
		b = *x
	default:
		ok = false
	}
	return
}

// ConvBytesSlice tries to convert value to slice of bytes.
func ConvBytesSlice(val any) (b [][]byte, ok bool) {
	ok = true
	switch x := val.(type) {
	case [][]byte:
		b = x
	case *[][]byte:
		b = *x
	default:
		ok = false
	}
	return
}

// ConvStr tries to convert value to string.
func ConvStr(val any) (s string, ok bool) {
	ok = true
	switch x := val.(type) {
	case string:
		s = x
	case *string:
		s = *x
	default:
		ok = false
	}
	return
}

// ConvStrSlice tries to convert value to string slice.
func ConvStrSlice(val any) (s []string, ok bool) {
	ok = true
	switch x := val.(type) {
	case []string:
		s = x
	case *[]string:
		s = *x
	default:
		ok = false
	}
	return
}

// ConvBool tries to convert value ti boolean.
func ConvBool(val any) (b bool, ok bool) {
	ok = true
	switch x := val.(type) {
	case bool:
		b = x
	case *bool:
		b = *x
	default:
		ok = false
	}
	return
}

// Convert interface value with arbitrary underlying type to integer value.
func if2int(raw any) (r int64, ok bool) {
	ok = true
	switch x := raw.(type) {
	case int:
		r = int64(x)
	case *int:
		r = int64(*x)
	case int8:
		r = int64(x)
	case *int8:
		r = int64(*x)
	case int16:
		r = int64(x)
	case *int16:
		r = int64(*x)
	case int32:
		r = int64(x)
	case *int32:
		r = int64(*x)
	case int64:
		r = x
	case *int64:
		r = *x
	case uint:
		r = int64(x)
	case *uint:
		r = int64(*x)
	case uint8:
		r = int64(x)
	case *uint8:
		r = int64(*x)
	case uint16:
		r = int64(x)
	case *uint16:
		r = int64(*x)
	case uint32:
		r = int64(x)
	case *uint32:
		r = int64(*x)
	case uint64:
		r = int64(x)
	case *uint64:
		r = int64(*x)
	case []byte:
		if len(x) > 0 {
			r, _ = strconv.ParseInt(byteconv.B2S(x), 0, 0)
		}
	case *[]byte:
		if len(*x) > 0 {
			r, _ = strconv.ParseInt(byteconv.B2S(*x), 0, 0)
		}
	case string:
		if len(x) > 0 {
			r, _ = strconv.ParseInt(x, 0, 0)
		}
	case *string:
		if len(*x) > 0 {
			r, _ = strconv.ParseInt(*x, 0, 0)
		}
	case *bytebuf.Chain:
		if (*x).Len() > 0 {
			r, _ = strconv.ParseInt((*x).String(), 0, 0)
		}
	case intConverter:
		if i, err := x.Int(); err == nil {
			return i, true
		}
	default:
		ok = false
	}
	return
}
