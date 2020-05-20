package dyntpl

type ModFn func(ctx *Ctx, buf *interface{}, val interface{}, args []interface{}) error

type mod struct {
	id  []byte
	fn  *ModFn
	arg []*modArg
}

type modArg struct {
	val    []byte
	static bool
}

var (
	modRegistry = map[string]ModFn{}
)

func RegisterModFn(name, alias string, mod ModFn) {
	modRegistry[name] = mod
	if len(alias) > 0 {
		modRegistry[alias] = mod
	}
}

func GetModFn(name string) *ModFn {
	if fn, ok := modRegistry[name]; ok {
		return &fn
	}
	return nil
}
