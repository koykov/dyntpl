package dyntpl

// Global describes value of global variable.
type Global any

var (
	globIdx = map[string]int{}
	globBuf []any
)

// RegisterGlobal registers new global variable.
func RegisterGlobal(name, alias string, val Global) {
	globBuf = append(globBuf, val)
	globIdx[name] = len(globBuf) - 1
	if len(alias) > 0 {
		globIdx[alias] = len(globBuf) - 1
	}
}

// RegisterGlobalNS registers new global variable in given namespace.
func RegisterGlobalNS(namespace, name, alias string, val Global) {
	if len(namespace) == 0 {
		RegisterGlobal(name, alias, val)
		return
	}
	name = namespace + "::" + name
	if len(alias) > 0 {
		alias = namespace + "::" + alias
	}
	RegisterGlobal(name, alias, val)
}

// GetGlobal returns global variable by given name.
func GetGlobal(name string) Global {
	if idx, ok := globIdx[name]; ok && idx >= 0 && idx < len(globBuf) {
		return globBuf[idx]
	}
	return nil
}

var _, _ = RegisterGlobalNS, GetGlobal
