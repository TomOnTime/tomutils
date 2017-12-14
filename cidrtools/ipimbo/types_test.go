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
			t.Errorf("bad #%d: expected: %v got: %v", i, tt.e, x)
		}
	}
}
