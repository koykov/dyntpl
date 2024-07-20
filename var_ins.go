package dyntpl

import (
	"github.com/koykov/inspector"
)

type VarInsTuple struct {
	docgen
	ins inspector.Inspector
}

var (
	// Registry of variable-inspector mappings.
	varInsRegistry = map[string]int{}
	varInsBuf      []VarInsTuple
)

// RegisterVarInsPair registers new variable-inspector pair.
func RegisterVarInsPair(varName string, ins inspector.Inspector) *VarInsTuple {
	if idx, ok := varInsRegistry[varName]; ok && idx >= 0 && idx < len(varInsBuf) {
		return &varInsBuf[idx]
	}
	varInsBuf = append(varInsBuf, VarInsTuple{
		docgen: docgen{name: varName, ins: ins.TypeName()},
		ins:    ins,
	})
	idx := len(varInsBuf) - 1
	varInsRegistry[varName] = idx
	return &varInsBuf[idx]
}

// GetInsByVarName gets inspector by variable name.
func GetInsByVarName(varName string) (inspector.Inspector, bool) {
	if idx, ok := varInsRegistry[varName]; ok && idx >= 0 && idx < len(varInsBuf) {
		return varInsBuf[idx].ins, true
	}
	return nil, false
}

// GetInspector gets inspector by both variable name or inspector name.
func GetInspector(varName, name string) (ins inspector.Inspector, err error) {
	err = inspector.ErrUnknownInspector
	if len(name) > 0 {
		ins, err = inspector.GetInspector(name)
	}
	if err != nil && len(varName) > 0 {
		if ins, ok := GetInsByVarName(varName); ok {
			return ins, nil
		}
	}
	return ins, err
}

var _ = RegisterVarInsPair
