package utils_test

import (
	"testing"

	"github.com/swoga/mikrotik-exporter/utils"
)

var (
	validDoubles = []validDobuleTest{
		{"1", 1},
		{"1x", 1},
	}
	emptyStr       = ""
	xStr           = "x"
	invalidDoubles = []invalidDoubleTest{
		{nil, "input nil"},
		{&emptyStr, "strconv.ParseFloat: parsing \"\": invalid syntax"},
		{&xStr, "strconv.ParseFloat: parsing \"\": invalid syntax"},
	}
)

type validDobuleTest struct {
	in  string
	out float64
}

type invalidDoubleTest struct {
	in  *string
	err string
}

func TestParseDoubleValid(t *testing.T) {
	for _, d := range validDoubles {
		has, err := utils.TryParseDouble(&d.in)
		if err != nil {
			t.Fatalf("unexpected error (input: %v), got: %v, want: %v", d.in, err, d.out)
		}
		if d.out != has {
			t.Fatalf("unexpected output (input: %v), got: %v, want: %v", d.in, has, d.out)
		}
	}
}

func TestParseDoubleInvalid(t *testing.T) {
	for _, d := range invalidDoubles {
		has, err := utils.TryParseDouble(d.in)
		if err == nil {
			t.Fatalf("expected error (input: %v), got: %v, want: %v", d.in, has, d.err)
		}
		if err.Error() != d.err {
			t.Fatalf("unexpected error (input: %v), got: %v, want: %v", d.in, err, d.err)
		}
	}
}
