package cbytetpl

type Inspector interface {
	Get(src interface{}, path ...string) interface{}
	Set(dst, value interface{}, path ...string)
}
