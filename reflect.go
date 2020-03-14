package cbytetpl

import (
	"fmt"
	"reflect"
)

type ReflectInspector struct{}

func (i ReflectInspector) Get(src interface{}, path ...interface{}) interface{} {
	var (
		r interface{}
		c int
		k interface{}
	)
	r = src
	for c, k = range path {
		r = i.inspect(r, k)
	}
	if c < len(path)-1 {
		r = nil
	}
	return r
}

func (i ReflectInspector) Set(dst, value interface{}, path ...interface{}) {
	// Empty method, there is no way to update data using reflection.
}

func (i ReflectInspector) inspect(node interface{}, key interface{}) interface{} {
	v := reflect.ValueOf(node)
	switch v.Kind() {
	case reflect.Ptr:
		if elem := v.Elem(); elem.IsValid() && elem.CanInterface() {
			node = elem.Interface()
			return i.inspect(node, key)
		}
	case reflect.Map:
		kv := reflect.ValueOf(key)
		_ = kv
		for _, f := range v.MapKeys() {
			fv := f.Interface()
			if fvs, ok := fv.(string); ok {
				if fvs == key {
					mv := v.MapIndex(f)
					if mv.IsValid() && mv.CanInterface() {
						return mv.Interface()
					}
					return nil
				}
			}
			if fmt.Sprintf("%v", fv) == key {
				mv := v.MapIndex(f)
				if mv.IsValid() && mv.CanInterface() {
					return mv.Interface()
				}
			}
		}
	case reflect.Struct:
		f := v.FieldByName(key.(string))
		if f.IsValid() && f.CanInterface() {
			return f.Interface()
		}
	case reflect.Slice:
		if bytes, ok := node.([]byte); ok {
			return bytes
		}
		sv := v.Index(key.(int))
		if sv.IsValid() && sv.CanInterface() {
			return sv.Interface()
		}
	}
	return nil
}
