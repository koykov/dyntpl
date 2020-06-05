package dyntpl

import (
	"strconv"

	"github.com/koykov/fastconv"
)

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

func ConvBytes(val interface{}) (b []byte, ok bool) {
	ok = true
	switch val.(type) {
	case []byte:
		b = val.([]byte)
	case *[]byte:
		b = *val.(*[]byte)
	case *ByteBuf:
		b = *val.(*ByteBuf)
	default:
		ok = false
	}
	return
}

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

func if2int(raw interface{}) (lim int64, ok bool) {
	ok = true
	switch raw.(type) {
	case int:
		lim = int64(raw.(int))
	case *int:
		lim = int64(*raw.(*int))
	case int8:
		lim = int64(raw.(int8))
	case *int8:
		lim = int64(*raw.(*int8))
	case int16:
		lim = int64(raw.(int16))
	case *int16:
		lim = int64(*raw.(*int16))
	case int32:
		lim = int64(raw.(int32))
	case *int32:
		lim = int64(*raw.(*int32))
	case int64:
		lim = raw.(int64)
	case *int64:
		lim = *raw.(*int64)
	case uint:
		lim = int64(raw.(uint))
	case *uint:
		lim = int64(*raw.(*uint))
	case uint8:
		lim = int64(raw.(uint8))
	case *uint8:
		lim = int64(*raw.(*uint8))
	case uint16:
		lim = int64(raw.(uint16))
	case *uint16:
		lim = int64(*raw.(*uint16))
	case uint32:
		lim = int64(raw.(uint32))
	case *uint32:
		lim = int64(*raw.(*uint32))
	case uint64:
		lim = int64(raw.(uint64))
	case *uint64:
		lim = int64(*raw.(*uint64))
	case []byte:
		lim, _ = strconv.ParseInt(fastconv.B2S(raw.([]byte)), 0, 0)
	case *[]byte:
		lim, _ = strconv.ParseInt(fastconv.B2S(*raw.(*[]byte)), 0, 0)
	case string:
		lim, _ = strconv.ParseInt(raw.(string), 0, 0)
	case *string:
		lim, _ = strconv.ParseInt(*raw.(*string), 0, 0)
	case *ByteBuf:
		lim, _ = strconv.ParseInt(fastconv.B2S(*raw.(*ByteBuf)), 0, 0)
	default:
		ok = false
	}
	return
}
