package main

import (
	"testing"

	"github.com/TomOnTime/tomutils/cidrtools/ipimbo"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		a, b string
		e    bool
	}{
		// last, next, skip next
		{"1.2.3.0/24", "1.2.3.1", true},
		{"1.2.3.0/24", "1.2.5.1", false},
		{"10.0.0.0/8", "10.0.0.0/8", true},
		{"10.0.0.0/8", "10.10.10.0/24", true},
		{"10.10.10.0/24", "24.24.24.0/24", false},
		{"24.24.24.0/24", "24.24.24.24/32", true},
	}

	for i, tt := range tests {
		pa, err := ipimbo.Parseline(tt.a)
		if err != nil {
			t.Errorf("bad #%d: unexpected err: %v", i, err)
		}
		pb, err := ipimbo.Parseline(tt.b)
		if err != nil {
			t.Errorf("bad #%d: unexpected err: %v", i, err)
		}
		r := process(pa, pb, func(v ipimbo.Imbo) {})
		if (r.String() == tt.a) != tt.e {
			t.Errorf("bad #%d: (%s : %s) expected: %v got: %v", i, tt.a, tt.b, tt.e, r.String())
		}
	}
}

func TestJudge(t *testing.T) {
	tests := []struct {
		a, b string
		e1   int
		e2   bool
	}{
		{"10.0.0.0/8", "10.0.0.0/8", CONTAINS, true},
		{"10.0.0.0/8", "10.10.10.0/24", CONTAINS, true},
		{"10.0.0.0/8", "24.24.24.0/24", DISJOINT, true},
		{"24.24.24.0/24", "24.24.24.24/32", CONTAINS, true},
		{"24.24.24.0/24", "24.24.24.24", CONTAINS, true},
		{"24.24.24.0/24", "33.33.33.0/24", DISJOINT, true},
		{"33.33.33.0/24", "33.33.99.1", DISJOINT, true},
	}

	for i, tt := range tests {
		pa, err := ipimbo.Parseline(tt.a)
		if err != nil {
			t.Errorf("bad #%d: unexpected err: %v", i, err)
		}
		pb, err := ipimbo.Parseline(tt.b)
		if err != nil {
			t.Errorf("bad #%d: unexpected err: %v", i, err)
		}
		r, err := judge(pa, pb)
		if r != tt.e1 {
			t.Errorf("bad #%d: judge(%s : %s) expected: %v got: %v", i, tt.a, tt.b, tt.e1, r)
		}
	}
}
