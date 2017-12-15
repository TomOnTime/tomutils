package ipimbo

import (
	"net"
	"testing"
)

func TestParseline(t *testing.T) {
	tests := []struct {
		d string
		e Imbo
	}{
		{"1.2.3.0/24", Imbo{version: 4, addr: net.ParseIP("1.2.3.0"), prefixlen: 24, isIP: false, isZeroAddr: true}},
		{"1.2.3.1/24", Imbo{version: 4, addr: net.ParseIP("1.2.3.1"), prefixlen: 24, isIP: false, isZeroAddr: false}},
		{"1.2.3.1/32", Imbo{version: 4, addr: net.ParseIP("1.2.3.1"), prefixlen: 32, isIP: false, isZeroAddr: true}},
		{"1.2.3.0", Imbo{version: 4, addr: net.ParseIP("1.2.3.0"), prefixlen: 32, isIP: true, isZeroAddr: true}},
		{"1.2.3.1", Imbo{version: 4, addr: net.ParseIP("1.2.3.1"), prefixlen: 32, isIP: true, isZeroAddr: true}},
		{"fe80::5054:ff:fef3:3410", Imbo{version: 6, addr: net.ParseIP("fe80::5054:ff:fef3:3410"), prefixlen: 128, isIP: true, isZeroAddr: true}},
		{"fe80::5054:ff:fef3:3410/56", Imbo{version: 6, addr: net.ParseIP("fe80::5054:ff:fef3:3410"), prefixlen: 56, isIP: false, isZeroAddr: false}},
		{"::1", Imbo{version: 6, addr: net.ParseIP("::1"), prefixlen: 128, isIP: true, isZeroAddr: true}},
	}

	//for _, comment := range []string{"", " ", " foo", "\t", "\tbar", " \tbaz"} {
	for _, comment := range []string{""} {
		for i, tt := range tests {
			//fmt.Printf("TEST #%d\n", i)
			x, err := parseline(tt.d + comment)
			if err != nil {
				t.Errorf("bad #%d: unexpected err: %v", i, err)
			}
			if !x.EqualIP(tt.e) {
				t.Errorf("bad #%d: %v\nexpected: %v \n     got: %v", i, tt.d+comment, tt.e.DebugString(), x.DebugString())
			}
		}
	}
}
