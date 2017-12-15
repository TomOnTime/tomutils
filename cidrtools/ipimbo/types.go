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
	return fmt.Sprintf("x %+v", *b)
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

func (b *Imbo) less(a Imbo) bool {
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

type ImboList []Imbo

func (ipdb ImboList) Sort() {
	sort.SliceStable(ipdb[:], func(i, j int) bool {
		return ipdb[j].less(ipdb[i])
	})
}
