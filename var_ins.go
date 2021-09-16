package dyntpl

import (
	"github.com/koykov/inspector"
)

var (
	// Registry of variable-inspector mappings.
	varInsRegistry = map[string]inspector.Inspector{}
	// Suppress go vet warning.
	_ = RegisterVarInsPair
)

// RegisterVarInsPair registers new variable-inspector pair.
func RegisterVarInsPair(varName string, ins inspector.Inspector) {
	varInsRegistry[varName] = ins
}

// GetInsByVarName gets inspector by variable name.
func GetInsByVarName(varName string) (inspector.Inspector, bool) {
	if ins, ok := varInsRegistry[varName]; ok {
		return ins, true
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
