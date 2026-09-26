package greeter

import (
	"strings"
	"testing"
	"testing/quick"
	"unicode"
)

// These are PROPERTY-BASED tests. Instead of asserting on a few fixed
// examples, each test states a general rule (a "property") that must hold for
// ANY input, and testing/quick generates hundreds of random inputs to try to
// falsify it. This mirrors the README: enforce a general rule, run many cases
// against the high-level intent.

// Property 1: For any input, Greet output starts with "Hello, " and ends "!".
func TestProperty_GreetHasPrefixAndSuffix(t *testing.T) {
	f := func(name string) bool {
		out := Greet(name)
		return strings.HasPrefix(out, prefix) && strings.HasSuffix(out, suffix)
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 500}); err != nil {
		t.Error(err)
	}
}

// Property 2: For any input, the normalized name appears inside the output,
// exactly between the prefix and the suffix.
func TestProperty_GreetContainsNormalizedName(t *testing.T) {
	f := func(name string) bool {
		out := Greet(name)
		want := NormalizeName(name)
		inner := strings.TrimSuffix(strings.TrimPrefix(out, prefix), suffix)
		return inner == want
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 500}); err != nil {
		t.Error(err)
	}
}

// Property 3: NormalizeName is idempotent — normalizing twice equals once.
func TestProperty_NormalizeIsIdempotent(t *testing.T) {
	f := func(name string) bool {
		once := NormalizeName(name)
		twice := NormalizeName(once)
		return once == twice
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 500}); err != nil {
		t.Error(err)
	}
}

// Property 4: Any all-whitespace (or empty) name greets "World".
func TestProperty_BlankNameGreetsWorld(t *testing.T) {
	f := func(n uint8) bool {
		// Build a string of n spaces/tabs/newlines only.
		blanks := []rune{' ', '\t', '\n', '\r'}
		var b strings.Builder
		for i := 0; i < int(n); i++ {
			b.WriteRune(blanks[i%len(blanks)])
		}
		return Greet(b.String()) == "Hello, World!"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 300}); err != nil {
		t.Error(err)
	}
}

// Property 5: The normalized name never has leading/trailing whitespace.
func TestProperty_NormalizedNameHasNoOuterSpace(t *testing.T) {
	f := func(name string) bool {
		got := NormalizeName(name)
		if got == "" {
			return false // must never be empty (falls back to World)
		}
		firstOK := !unicode.IsSpace([]rune(got)[0])
		r := []rune(got)
		lastOK := !unicode.IsSpace(r[len(r)-1])
		return firstOK && lastOK
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 500}); err != nil {
		t.Error(err)
	}
}

// A couple of concrete example tests to complement the properties and document
// intent explicitly.
func TestExample_Basic(t *testing.T) {
	cases := map[string]string{
		"":        "Hello, World!",
		"   ":     "Hello, World!",
		"Alice":   "Hello, Alice!",
		"  Bob  ": "Hello, Bob!",
	}
	for in, want := range cases {
		if got := Greet(in); got != want {
			t.Errorf("Greet(%q) = %q; want %q", in, got, want)
		}
	}
}
