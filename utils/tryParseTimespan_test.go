package utils_test

import (
	"testing"

	"github.com/swoga/mikrotik-exporter/utils"
)

var (
	validTimespans = []validTimespanTest{
		{"", 0},
		{"1s", 1},
		{"1m", 60},
		{"1h", 60 * 60},
		{"1d", 24 * 60 * 60},
		{"1w", 7 * 24 * 60 * 60},
		{"1w20h1s", 7*24*60*60 + 20*60*60 + 1},
		{"1w20h1s100ms", 7*24*60*60 + 20*60*60 + 1 + 0.1},
	}
	invalidTimespans = []invalidTimespanTest{
		{"1y", "no regex match"},
		{"1y1w", "no regex match"},
	}
)

type validTimespanTest struct {
	in  string
	out float64
}

type invalidTimespanTest struct {
	in  string
	err string
}

func TestParseTimespanValid(t *testing.T) {
	for _, ts := range validTimespans {
		has, err := utils.TryParseTimespan(ts.in)
		if err != nil {
			t.Fatalf("unexpected error (input: %v), got: %v, want: %v", ts.in, err, ts.out)
		}
		if ts.out != has {
			t.Fatalf("unexpected output (input: %v), got: %v, want: %v", ts.in, has, ts.out)
		}
	}
}

func TestParseTimespanInvalid(t *testing.T) {
	for _, ts := range invalidTimespans {
		has, err := utils.TryParseTimespan(ts.in)
		if err == nil {
			t.Fatalf("expected error (input: %v), got: %v, want: %v", ts.in, has, ts.err)
		}
		if err.Error() != ts.err {
			t.Fatalf("unexpected error (input: %v), got: %v, want: %v", ts.in, err, ts.err)
		}
	}
}
