package dyntpl

import "regexp"

var (
	isStaticRE = regexp.MustCompile(`^\d+\.*\d*|true|false|nil|"[^"]*"|'[^']*'$`)
)

func isStatic(arg []byte) bool {
	return isStaticRE.Match(arg)
}
