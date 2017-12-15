package ipimbo

import (
	"net"
	"strings"

	"github.com/pkg/errors"
)

// parseline parses an input line and returns an Imbo. The line format
// is one of:
//     IPADDRESS
//     CIDRBLOCK
//     IPADDRESS <whitespace> comment
//     CIDRBLOCK <whitespace> comment
// IPv4 and IPv6 are supported.
// This function assumes that the string has been passed through
// strings.TrimSpace(), is non-empty, and is not a comment line.
func parseline(s string) (Imbo, error) {
	ipb := Imbo{}

	// Extract the first field, and the remainder.
	var first string
	if strings.IndexAny(s, " \t") < 0 {
		first = s
		ipb.comment = ""
	} else {
		fields := strings.Fields(s)
		first = fields[0]
		ipb.comment = s[len(first):]
	}

	// Try to parse the first field.
	// Is it a cidr block?
	ipAddr, ipNet, err := net.ParseCIDR(first)
	//var x int
	if err == nil {
		if ipAddr.To4() != nil {
			ipb.version = 4
		} else if ipAddr.To16() != nil {
			ipb.version = 6
		} else {
			return ipb, errors.Errorf("should not happen CIDR: %s", first)
		}
		ipb.addr = ipAddr
		ipb.prefixlen, _ = ipNet.Mask.Size()
		ipb.isIP = false
		ipb.isZeroAddr = ipAddr.Equal(ipNet.IP)
		return ipb, nil
	}
	// Is it an IP address?
	ip := net.ParseIP(first)
	if ip != nil {
		if ip.To4() != nil {
			ipb.version = 4
			ipb.prefixlen = 32
		} else if ip.To16() != nil {
			ipb.version = 6
			ipb.prefixlen = 128
		} else {
			return ipb, errors.Errorf("should not happen IP: %s", first)
		}
		ipb.addr = ip
		ipb.isIP = true
		ipb.isZeroAddr = true
		return ipb, nil
	}

	return ipb, errors.Errorf("invalid: first='%s' line='%v'", first, s)
}
