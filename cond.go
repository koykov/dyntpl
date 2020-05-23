package dyntpl

type CondFn func(ctx *Ctx, args []interface{}) bool

var (
	condRegistry = map[string]CondFn{}
)

func RegisterCondFn(name string, cond CondFn) {
	condRegistry[name] = cond
}

func GetCondFn(name string) *CondFn {
	if fn, ok := condRegistry[name]; ok {
		return &fn
	}
	return nil
}
