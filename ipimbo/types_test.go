package ipimbo

import (
	"net"
	"testing"
)

func TestIpimboString(t *testing.T) {
	ipb := Imbo{version: 6, addr: net.ParseIP("1.2.3.0"), prefixlen: 24, isIP: false, isIP: true}
	x := ipb.String()
	if x != "1.2.3.0/24" {
		t.Errorf("bad output: expected: %v got: %v", "1.2.3.0/24", x)
	}
}
