package dyntpl

import "regexp"

var (
	// Regexp to check is argument is static value.
	isStaticRE = regexp.MustCompile(`^\d+\.*\d*|true|false|nil|"[^"]*"|'[^']*'$`)
)

// Check if arg is static value.
func isStatic(arg []byte) bool {
	return isStaticRE.Match(arg)
}
