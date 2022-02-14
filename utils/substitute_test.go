package utils_test

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/swoga/mikrotik-exporter/utils"
)

var (
	substitutes = []substituteTest{
		{"x", "x"},
		{"a", "a"},
		{"{a}", "1"},
		{"{a}{b}", "12"},
		{"x{a}x{b}x", "x1x2x"},
	}
	substituteVars = map[string]string{
		"a": "1",
		"b": "2",
	}
)

type substituteTest struct {
	in  string
	out string
}

func TestSubstitute(t *testing.T) {
	for _, s := range substitutes {
		has := utils.Substitute(log.Logger, s.in, substituteVars)
		if s.out != has {
			t.Fatalf("unexpected output (input: %v), got: %v, want: %v", s.in, has, s.out)
		}
	}
}
