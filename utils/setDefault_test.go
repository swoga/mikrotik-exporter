package utils_test

import (
	"testing"

	"github.com/swoga/mikrotik-exporter/utils"
)

func TestSetDefault(t *testing.T) {
	testdata := [][]string{
		{"x", "x"},
		{"", "y"},
	}

	def := "y"

	for _, test := range testdata {
		v := test[0]
		want := test[1]

		utils.SetDefaultString(&v, def)
		if v != want {
			t.Fatalf("unexpected output, want: %v, got: %v", want, v)
		}
	}
}
