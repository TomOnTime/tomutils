package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/yl2chen/cidranger"
)

type IpInfo struct {
	ip   net.IP
	line string
}

func readRoutes(f io.Reader, ipList []IpInfo) ([]IpInfo, error) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse the line. It should start with an IP address.
		trim := strings.TrimSpace(line)
		i := strings.IndexAny(trim, "/ \t")
		if i == 0 { // Skip empty lines.
			continue
		}
		if i == -1 {
			i = len(trim)
		}
		if i < 4 { // Too short to be an IP? error
			return nil, fmt.Errorf("Line has no IP or CIDR: %v", line)
		}
		ip := net.ParseIP(trim[0:i])
		if ip == nil {
			return nil, fmt.Errorf("Line does not start with valid IP: %v", line)
		}

		ipList = append(ipList, IpInfo{ip, line})
	}

	return ipList, nil
}

type RouteEntry struct {
	Ip          net.IPNet
	Destination string
}

func (*RouteEntry) Network() net.IPNet {
	return RouteEntry.Ip
}

func main() {
	flag.Parse()
	var err error

	// Read all input:
	if flag.NArg() != 1 {
		log.Fatal("Missing routetable filename on command line")
	}

	ranger := cidranger.NewPCTrieRanger()

	// Read any/all route files:
	for _, fname := range os.Args[1:] {
		fh, err := os.Open(fname)
		defer fh.Close()
		if err != nil {
			log.Fatal(err)
		}
		err = readRoutes(fh, ranger)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Print it:
	containingNetworks, err = ranger.ContainingNetworks(net.ParseIP("0.0.0.0"))
	for _, x := range containingNetworks {
		fmt.Println(x.Data())
	}

	// Read Stdin, look up each item.

}
