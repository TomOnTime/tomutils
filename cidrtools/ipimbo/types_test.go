package ipimbo

import (
	"net"
	"testing"
)

func TestIpimboString(t *testing.T) {
	tests := []struct {
		d Imbo
		e string
	}{
		{Imbo{version: 4, addr: net.ParseIP("1.2.3.0"), prefixlen: 24, isIP: false, isZeroAddr: true},
			"1.2.3.0/24"},
		{Imbo{version: 4, addr: net.ParseIP("1.2.3.1"), prefixlen: 24, isIP: false, isZeroAddr: false},
			"1.2.3.1/24"},
		{Imbo{version: 6, addr: net.ParseIP("fe80::5054:ff:fef3:3410"), prefixlen: 24, isIP: false, isZeroAddr: false},
			"fe80::5054:ff:fef3:3410/24"},
	}

	for i, tt := range tests {
		x := tt.d.String()
		if x != tt.e {
			t.Errorf("test #%d: expected: %v got: %v", i, tt.e, x)
		}
	}
}

func TestLess(t *testing.T) {
	tests := []struct {
		a, b Imbo
	}{
		// a should be less than b
		{
			Imbo{version: 4, addr: net.ParseIP("1.2.3.0"), prefixlen: 24, isIP: true},
			Imbo{version: 4, addr: net.ParseIP("1.2.3.1"), prefixlen: 24, isIP: true},
		},
		{
			Imbo{version: 4, addr: net.ParseIP("1.2.3.0"), prefixlen: 32, isIP: false},
			Imbo{version: 4, addr: net.ParseIP("1.2.3.0"), prefixlen: 24, isIP: true},
		},
		{
			Imbo{version: 4, addr: net.ParseIP("1.2.3.0"), prefixlen: 24, isIP: false},
			Imbo{version: 4, addr: net.ParseIP("1.2.3.0"), prefixlen: 32, isIP: false},
		},
		{
			Imbo{version: 4, addr: net.ParseIP("10.2.3.0"), prefixlen: 24, isIP: true},
			Imbo{version: 6, addr: net.ParseIP("1::1"), prefixlen: 24, isIP: true},
		},
		{
			Imbo{version: 4, addr: net.ParseIP("10.2.3.0"), prefixlen: 24, isIP: true},
			Imbo{version: 6, addr: net.ParseIP("11::1"), prefixlen: 24, isIP: true},
		},
		{
			Imbo{version: 4, addr: net.ParseIP("10.2.3.0"), prefixlen: 32, isIP: true},
			Imbo{version: 4, addr: net.ParseIP("100.2.3.0"), prefixlen: 32, isIP: true},
		},
		{
			Imbo{version: 4, addr: net.ParseIP("0.0.0.0"), prefixlen: 32, isIP: true},
			Imbo{version: 6, addr: net.ParseIP("::1"), prefixlen: 128, isIP: true},
		},
		{
			Imbo{version: 4, addr: net.ParseIP("10.0.0.0"), prefixlen: 32, isIP: true},
			Imbo{version: 6, addr: net.ParseIP("::1"), prefixlen: 128, isIP: true},
		},
		{
			Imbo{version: 4, addr: net.ParseIP("10.0.0.0"), prefixlen: 32, isIP: true},
			Imbo{version: 6, addr: net.ParseIP("::1"), prefixlen: 128, isIP: true},
		},
	}

	for i, tt := range tests {
		if !tt.b.less(tt.a) {
			t.Errorf("test #%d: %s expected to be less than %s", i, tt.a.String(), tt.b.String())
		}
		if tt.a.less(tt.b) {
			t.Errorf("test #%d: %s expected to be less than %s", i, tt.b.String(), tt.a.String())
		}
	}
}
