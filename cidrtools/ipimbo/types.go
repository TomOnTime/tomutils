package ipimbo

import (
	"fmt"
	"net"
)

type Imbo struct {
	version    int
	addr       net.IP
	prefixlen  uint8
	isIP       bool
	isZeroAddr bool
}

func (b *Imbo) String() string {
	if b.isIP {
		return b.addr.String()
	}
	return fmt.Sprintf("%s/%d", b.addr, b.prefixlen)
}
