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
	fmt.Printf("EqualIP\n%+v\n%+v == ", a, *b)
	if a.version != b.version {
		fmt.Println("false (VERSION)")
		return false
	}
	if a.prefixlen != b.prefixlen {
		fmt.Println("false (PREFIX)")
		return false
	}
	if a.isIP != b.isIP {
		fmt.Println("false (ISIP)")
		return false
	}
	if a.isZeroAddr != b.isZeroAddr {
		fmt.Println("false (ZERO)")
		return false
	}
	r := a.addr.Equal(b.addr)
	fmt.Printf("%v (EQUAL %+v %+v)\n", r, a.addr, b.addr)
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
