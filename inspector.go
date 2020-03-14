package cbytetpl

type Inspector interface {
	Get(src interface{}, path ...interface{}) interface{}
	Set(dst, value interface{}, path ...interface{})
}
