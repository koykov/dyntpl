package dyntpl

import (
	"strconv"

	"github.com/koykov/fastconv"
)

func lim2int(raw interface{}) (lim int64, ok bool) {
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
	default:
		ok = false
	}
	return
}
