package dyntpl

// Global describes value of global variable.
type Global any

type GlobalTuple struct {
	id, typ, desc string
	val           Global
}

func (t *GlobalTuple) WithType(typ string) *GlobalTuple {
	t.typ = typ
	return t
}

func (t *GlobalTuple) WithDescription(desc string) *GlobalTuple {
	t.desc = desc
	return t
}

var (
	globIdx = map[string]int{}
	globBuf []GlobalTuple
)

// RegisterGlobal registers new global variable.
//
// Caution! Globals registered after template parsing will take no effect.
func RegisterGlobal(name, alias string, val Global) *GlobalTuple {
	if idx, ok := globIdx[name]; ok && idx >= 0 && idx < len(globBuf) {
		return &globBuf[idx]
	}
	globBuf = append(globBuf, GlobalTuple{
		id:  name,
		typ: "any",
		val: val,
	})
	idx := len(globBuf) - 1
	globIdx[name] = idx
	if len(alias) > 0 {
		globIdx[alias] = idx
	}
	return &globBuf[idx]
}

// RegisterGlobalNS registers new global variable in given namespace.
func RegisterGlobalNS(namespace, name, alias string, val Global) *GlobalTuple {
	if len(namespace) == 0 {
		return RegisterGlobal(name, alias, val)
	}
	name = namespace + "::" + name
	if len(alias) > 0 {
		alias = namespace + "::" + alias
	}
	return RegisterGlobal(name, alias, val)
}

// GetGlobal returns global variable by given name.
func GetGlobal(name string) Global {
	if idx, ok := globIdx[name]; ok && idx >= 0 && idx < len(globBuf) {
		return globBuf[idx].val
	}
	return nil
}

var _, _ = RegisterGlobalNS, GetGlobal
