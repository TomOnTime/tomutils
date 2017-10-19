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

type RouteEntry struct {
	Ip   net.IPNet
	dest string
}

func (re RouteEntry) Network() net.IPNet {
	return re.Ip
}

func (re RouteEntry) Dest() string {
	return re.dest
}

func readRoutes(f io.Reader, ranger cidranger.Ranger) error {
	var err error
	var network1 *net.IPNet

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("READ line=", line)

		// Parse the line. It should start with an IP address.
		trim := strings.TrimSpace(line)
		_, network1, err = net.ParseCIDR(trim)
		if err != nil {
			log.Fatal(err)
		}

		err = ranger.Insert(RouteEntry{*network1, line})
		if err != nil {
			log.Fatal(err)
		}

		i := strings.IndexAny(trim, "/ \t")
		if i == 0 { // Skip empty lines.
			continue
		}
		if i == -1 {
			i = len(trim)
		}
		if i < 4 { // Too short to be an IP? error
			return fmt.Errorf("Line has no IP or CIDR: %v", line)
		}
		ip := net.ParseIP(trim[0:i])
		if ip == nil {
			return fmt.Errorf("Line does not start with valid IP: %v", line)
		}

	}
	return nil
}

func main() {
	flag.Parse()

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

	fmt.Println(ranger)

	// Print it:
	containingNetworks, err := ranger.ContainingNetworks(net.ParseIP("0.0.0.0"))
	if err != nil {
		log.Fatal(err)
	}
	for _, x := range containingNetworks {
		fmt.Println(x.(RouteEntry).Dest())
	}

	// Read Stdin, look up each item.

}
