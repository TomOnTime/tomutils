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
		if !tt.b.Less(tt.a) {
			t.Errorf("test #%d: %s expected to be less than %s", i, tt.a.String(), tt.b.String())
		}
		if tt.a.Less(tt.b) {
			t.Errorf("test #%d: %s expected to be less than %s", i, tt.b.String(), tt.a.String())
		}
	}
}

func TestHostAllZerosv(t *testing.T) {
	tests := []struct {
		d string
		e bool
	}{
		{"1.2.3.0/24", true},
		{"1.2.3.1/24", false},
		{"1.2.3.1/32", true},
		{"1.2.3.0", true},
		{"1.2.3.1", true},
		{"fe80::5054:ff:fef3:3410", true},
		{"fe80::5054:ff:fef3:3410/56", false},
		{"fe80:9999::/56", true},
		{"fe80:9999::1/56", false},
	}

	for i, tt := range tests {
		x, err := parseline(tt.d)
		if err != nil {
			t.Errorf("bad #%d: unexpected err: %v", i, err)
		}
		r := x.HostAllZeros()
		if r != tt.e {
			t.Errorf("bad #%d: %s expected: %v got: %v", i, tt.d, tt.e, r)
		}
	}
}

func TestZeroOutHostBits(t *testing.T) {
	tests := []struct {
		d string
		e string
	}{
		{"1.2.3.0/24", "1.2.3.0"},
		{"1.2.3.1/24", "1.2.3.0"},
		{"1.2.3.1/32", "1.2.3.1"},
		{"255.255.255.255/17", "255.255.128.0"},
		{"1.2.3.0", "1.2.3.0"},
		{"1.2.3.1", "1.2.3.1"},
		{"fe80::5054:ff:fef3:3410", "fe80::5054:ff:fef3:3410"},
		{"fe80::5054:ff:fef3:3410/56", "fe80::"},
		{"fe80:9999::/56", "fe80:9999::"},
		{"fe80:9999::1/56", "fe80:9999::"},
	}

	for i, tt := range tests {
		x, err := parseline(tt.d)
		if err != nil {
			t.Errorf("bad #%d: unexpected err: %v", i, err)
		}
		r := x.ZeroOutHostBits().String()
		if r != tt.e {
			t.Errorf("bad #%d: %s expected: %v got: %v", i, tt.d, tt.e, r)
		}
	}
}
func TestContains(t *testing.T) {
	tests := []struct {
		a, b string
		e    bool
	}{
		{"10.0.0.0/24", "10.10.10.0/24", false},
		{"10.0.0.0/8", "10.10.10.0/24", true},
		{"24.24.24.0/24", "24.24.24.24", true},
		{"1.2.3.0/24", "1.2.3.1", true},
		{"9.1.1.0/24", "1.1.1.1", false},
		{"9.1.1.0/24", "9.2.2.2", false},
		{"255.255.255.0/24", "255.255.254.0/23", false},
		{"255.255.255.0/24", "255.255.255.0/24", true},
		{"255.255.255.0/24", "255.255.255.0/25", true},
		{"255.255.255.0/24", "255.255.255.0/26", true},
		{"abcd:1234::/32", "abcd:1234:ff::1", true},
		{"abcd:1234::/32", "2222:1234:ff::1", false},
		{"ffff:ffff::/32", "ffff:fffe::/31", false},
		{"ffff:ffff::/32", "ffff:ffff::/32", true},
		{"ffff:ffff::/32", "ffff:ffff:8::/33", true},
		{"ffff:ffff::/32", "ffff:ffff:c::/34", true},
	}

	for i, tt := range tests {
		pa, err := parseline(tt.a)
		if err != nil {
			t.Errorf("bad #%d: unexpected err: %v", i, err)
		}
		pb, err := parseline(tt.b)
		if err != nil {
			t.Errorf("bad #%d: unexpected err: %v", i, err)
		}
		r := pa.Contains(pb)
		if r != tt.e {
			t.Errorf("bad #%d: (%s : %s) expected: %v got: %v", i, tt.a, tt.b, tt.e, r)
		}
	}
}
