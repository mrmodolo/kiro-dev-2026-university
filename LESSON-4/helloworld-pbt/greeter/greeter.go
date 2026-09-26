// Package greeter provides a small "Hello, World!" style greeting.
//
// It is intentionally written so that its behavior can be described by
// general PROPERTIES (rules that must always hold for any input), which is
// what property-based testing verifies — instead of a handful of hand-picked
// examples.
package greeter

import (
	"fmt"
	"strings"
)

const (
	prefix       = "Hello, "
	suffix       = "!"
	defaultName  = "World"
)

// NormalizeName trims surrounding whitespace from a name and falls back to
// "World" when the result is empty. Normalizing an already-normalized name
// returns it unchanged (idempotence).
func NormalizeName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return defaultName
	}
	return trimmed
}

// Greet returns a greeting of the form "Hello, <name>!".
//
// Properties that always hold, for ANY input string:
//   - the result starts with "Hello, " and ends with "!"
//   - the normalized name is contained in the result
//   - a blank/empty name greets "World"
func Greet(name string) string {
	return fmt.Sprintf("%s%s%s", prefix, NormalizeName(name), suffix)
}
