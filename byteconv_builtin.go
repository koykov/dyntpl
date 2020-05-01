package cbytetpl

import (
	"strconv"

	"github.com/koykov/fastconv"
)

func byteConvBytes(buf []byte, val interface{}) ([]byte, error) {
	switch val.(type) {
	case *[]byte:
		buf = append(buf, *val.(*[]byte)...)
	case []byte:
		buf = append(buf, val.([]byte)...)
	default:
		return buf, ErrUnknownType
	}

	return buf, nil
}

func byteConvStr(buf []byte, val interface{}) ([]byte, error) {
	switch val.(type) {
	case *string:
		buf = append(buf, fastconv.S2B(*val.(*string))...)
	case []byte:
		buf = append(buf, fastconv.S2B(val.(string))...)
	default:
		return buf, ErrUnknownType
	}

	return buf, nil
}

func byteConvInt(buf []byte, val interface{}) ([]byte, error) {
	var i int64
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
		return buf, ErrUnknownType
	}

	buf = strconv.AppendInt(buf, i, 10)
	return buf, nil
}

func byteConvUint(buf []byte, val interface{}) ([]byte, error) {
	var i uint64
	switch val.(type) {
	case uint:
		i = uint64(val.(uint))
	case *uint:
		i = uint64(*val.(*uint))
	case uint8:
		i = uint64(val.(uint8))
	case *uint8:
		i = uint64(*val.(*uint8))
	case uint16:
		i = uint64(val.(uint16))
	case *uint16:
		i = uint64(*val.(*uint16))
	case uint32:
		i = uint64(val.(uint32))
	case *uint32:
		i = uint64(*val.(*uint32))
	case uint64:
		i = val.(uint64)
	case *uint64:
		i = *val.(*uint64)
	default:
		return buf, ErrUnknownType
	}

	buf = strconv.AppendUint(buf, i, 10)
	return buf, nil
}

func byteConvFloat(buf []byte, val interface{}) ([]byte, error) {
	var f float64
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
		return buf, ErrUnknownType
	}

	buf = strconv.AppendFloat(buf, f, 'f', -1, 64)
	return buf, nil
}
