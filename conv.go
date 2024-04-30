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
	switch val.(type) {
	case int:
		i = int64(val.(int))
	case *int:
		i = int64(*val.(*int))
	case int8:
		i = int64(val.(int8))
	case *int8:
		i = int64(*val.(*int8))
	case int16:
		i = int64(val.(int16))
	case *int16:
		i = int64(*val.(*int16))
	case int32:
		i = int64(val.(int32))
	case *int32:
		i = int64(*val.(*int32))
	case int64:
		i = val.(int64)
	case *int64:
		i = *val.(*int64)
	default:
		ok = false
	}
	return
}

// ConvUint tries to convert value to uint.
func ConvUint(val any) (u uint64, ok bool) {
	ok = true
	switch val.(type) {
	case uint:
		u = uint64(val.(uint))
	case *uint:
		u = uint64(*val.(*uint))
	case uint8:
		u = uint64(val.(uint8))
	case *uint8:
		u = uint64(*val.(*uint8))
	case uint16:
		u = uint64(val.(uint16))
	case *uint16:
		u = uint64(*val.(*uint16))
	case uint32:
		u = uint64(val.(uint32))
	case *uint32:
		u = uint64(*val.(*uint32))
	case uint64:
		u = val.(uint64)
	case *uint64:
		u = *val.(*uint64)
	default:
		ok = false
	}
	return
}

// ConvFloat tries to convert value to float.
func ConvFloat(val any) (f float64, ok bool) {
	ok = true
	switch val.(type) {
	case float32:
		f = float64(val.(float32))
	case *float32:
		f = float64(*val.(*float32))
	case float64:
		f = val.(float64)
	case *float64:
		f = *val.(*float64)
	default:
		ok = false
	}
	return
}

// ConvBytes tries to convert value to bytes.
func ConvBytes(val any) (b []byte, ok bool) {
	ok = true
	switch val.(type) {
	case []byte:
		b = val.([]byte)
	case *[]byte:
		b = *val.(*[]byte)
	case bytebuf.Chain:
		b = val.(bytebuf.Chain)
	case *bytebuf.Chain:
		b = *val.(*bytebuf.Chain)
	default:
		ok = false
	}
	return
}

// ConvBytesSlice tries to convert value to slice of bytes.
func ConvBytesSlice(val any) (b [][]byte, ok bool) {
	ok = true
	switch val.(type) {
	case [][]byte:
		b = val.([][]byte)
	case *[][]byte:
		b = *val.(*[][]byte)
	default:
		ok = false
	}
	return
}

// ConvStr tries to convert value to string.
func ConvStr(val any) (s string, ok bool) {
	ok = true
	switch val.(type) {
	case string:
		s = val.(string)
	case *string:
		s = *val.(*string)
	default:
		ok = false
	}
	return
}

// ConvStrSlice tries to convert value to string slice.
func ConvStrSlice(val any) (s []string, ok bool) {
	ok = true
	switch val.(type) {
	case []string:
		s = val.([]string)
	case *[]string:
		s = *val.(*[]string)
	default:
		ok = false
	}
	return
}

// ConvBool tries to convert value ti boolean.
func ConvBool(val any) (b bool, ok bool) {
	ok = true
	switch val.(type) {
	case bool:
		b = val.(bool)
	case *bool:
		b = *val.(*bool)
	default:
		ok = false
	}
	return
}

// Convert interface value with arbitrary underlying type to integer value.
func if2int(raw any) (r int64, ok bool) {
	ok = true
	switch raw.(type) {
	case int:
		r = int64(raw.(int))
	case *int:
		r = int64(*raw.(*int))
	case int8:
		r = int64(raw.(int8))
	case *int8:
		r = int64(*raw.(*int8))
	case int16:
		r = int64(raw.(int16))
	case *int16:
		r = int64(*raw.(*int16))
	case int32:
		r = int64(raw.(int32))
	case *int32:
		r = int64(*raw.(*int32))
	case int64:
		r = raw.(int64)
	case *int64:
		r = *raw.(*int64)
	case uint:
		r = int64(raw.(uint))
	case *uint:
		r = int64(*raw.(*uint))
	case uint8:
		r = int64(raw.(uint8))
	case *uint8:
		r = int64(*raw.(*uint8))
	case uint16:
		r = int64(raw.(uint16))
	case *uint16:
		r = int64(*raw.(*uint16))
	case uint32:
		r = int64(raw.(uint32))
	case *uint32:
		r = int64(*raw.(*uint32))
	case uint64:
		r = int64(raw.(uint64))
	case *uint64:
		r = int64(*raw.(*uint64))
	case []byte:
		if len(raw.([]byte)) > 0 {
			r, _ = strconv.ParseInt(byteconv.B2S(raw.([]byte)), 0, 0)
		}
	case *[]byte:
		if len(*raw.(*[]byte)) > 0 {
			r, _ = strconv.ParseInt(byteconv.B2S(*raw.(*[]byte)), 0, 0)
		}
	case string:
		if len(raw.(string)) > 0 {
			r, _ = strconv.ParseInt(raw.(string), 0, 0)
		}
	case *string:
		if len(*raw.(*string)) > 0 {
			r, _ = strconv.ParseInt(*raw.(*string), 0, 0)
		}
	case *bytebuf.Chain:
		if (*raw.(*bytebuf.Chain)).Len() > 0 {
			r, _ = strconv.ParseInt((*raw.(*bytebuf.Chain)).String(), 0, 0)
		}
	case intConverter:
		if i, err := raw.(intConverter).Int(); err == nil {
			return i, true
		}
	default:
		ok = false
	}
	return
}
