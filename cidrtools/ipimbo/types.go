package ipimbo

import (
	"bytes"
	"fmt"
	"net"
	"sort"
)

type Imbo struct {
	version    int
	addr       net.IP
	prefixlen  int
	isIP       bool
	isZeroAddr bool
	comment    string
}

// IPString prints the IP/CIDR.
func (b *Imbo) IPString() string {
	if b.isIP {
		return b.addr.String()
	}
	return fmt.Sprintf("%s/%d", b.addr, b.prefixlen)
}

// String prints the IP/CIDR followed by the comment (if there is one).
func (b *Imbo) String() string {
	if b.comment == "" {
		return b.IPString()
	} else {
		return b.IPString() + b.comment
	}
}

func (b *Imbo) DebugString() string {
	return fmt.Sprintf("%+v", *b)
}

func (b *Imbo) EqualIP(a Imbo) bool {
	if a.version != b.version {
		return false
	}
	if a.prefixlen != b.prefixlen {
		return false
	}
	if a.isIP != b.isIP {
		return false
	}
	if a.isZeroAddr != b.isZeroAddr {
		return false
	}
	r := a.addr.Equal(b.addr)
	return r
}

func (b *Imbo) Less(a Imbo) bool {
	if a.version != b.version {
		return a.version != 6
	}
	if !a.addr.Equal(b.addr) {
		return bytes.Compare(a.addr, b.addr) < 0
	}
	if a.isIP != b.isIP {
		return a.isIP == false
	}
	return a.prefixlen < b.prefixlen
}

func (a *Imbo) Version() int {
	return a.version
}

func (a *Imbo) PrefixLen() int {
	return a.prefixlen
}

func (a *Imbo) HostAllZeros() bool {
	return a.isZeroAddr
}

func (a *Imbo) AddressLen() int {
	if a.version == 4 {
		return 32
	} else if a.version == 6 {
		return 128
	} else {
		panic("corrupted IP: not IPv4 or IPv6")
	}
}

func (a *Imbo) ZeroOutHostBits() net.IP {
	mask := net.CIDRMask(a.prefixlen, a.AddressLen())
	return a.addr.Mask(mask)
}

// mask returns an IP with all but the first n bits zero'ed out.
func (a *Imbo) mask(n int) net.IP {
	mask := net.CIDRMask(n, a.AddressLen())
	return a.addr.Mask(mask)
}

func (a *Imbo) Contains(b Imbo) bool {
	n := net.IPNet{
		IP:   a.addr,
		Mask: net.CIDRMask(a.prefixlen, a.AddressLen()),
	}
	return n.Contains(b.addr)
}

type ImboList []Imbo

func (ipdb ImboList) Sort() {
	sort.SliceStable(ipdb[:], func(i, j int) bool {
		return ipdb[j].Less(ipdb[i])
	})
}
