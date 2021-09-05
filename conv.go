package dyntpl

// Collection of conversion functions.

import (
	"strconv"

	"github.com/koykov/bytebuf"
	"github.com/koykov/fastconv"
)

type intConverter interface {
	Int() (int64, error)
}

// Try to convert value to integer.
func ConvInt(val interface{}) (i int64, ok bool) {
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

// Try to convert value to uint.
func ConvUint(val interface{}) (u uint64, ok bool) {
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

// Try to convert value ti float.
func ConvFloat(val interface{}) (f float64, ok bool) {
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

// Try to convert value to bytes.
func ConvBytes(val interface{}) (b []byte, ok bool) {
	ok = true
	switch val.(type) {
	case []byte:
		b = val.([]byte)
	case *[]byte:
		b = *val.(*[]byte)
	case bytebuf.ChainBuf:
		b = val.(bytebuf.ChainBuf)
	case *bytebuf.ChainBuf:
		b = *val.(*bytebuf.ChainBuf)
	default:
		ok = false
	}
	return
}

// Try to convert value to slice of byte slice.
func ConvBytesSlice(val interface{}) (b [][]byte, ok bool) {
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

// Try to convert value to string.
func ConvStr(val interface{}) (s string, ok bool) {
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

// Try to convert value to string slice.
func ConvStrSlice(val interface{}) (s []string, ok bool) {
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

// Try to convert value ti boolean.
func ConvBool(val interface{}) (b bool, ok bool) {
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
func if2int(raw interface{}) (r int64, ok bool) {
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
			r, _ = strconv.ParseInt(fastconv.B2S(raw.([]byte)), 0, 0)
		}
	case *[]byte:
		if len(*raw.(*[]byte)) > 0 {
			r, _ = strconv.ParseInt(fastconv.B2S(*raw.(*[]byte)), 0, 0)
		}
	case string:
		if len(raw.(string)) > 0 {
			r, _ = strconv.ParseInt(raw.(string), 0, 0)
		}
	case *string:
		if len(*raw.(*string)) > 0 {
			r, _ = strconv.ParseInt(*raw.(*string), 0, 0)
		}
	case *bytebuf.ChainBuf:
		if (*raw.(*bytebuf.ChainBuf)).Len() > 0 {
			r, _ = strconv.ParseInt((*raw.(*bytebuf.ChainBuf)).String(), 0, 0)
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
