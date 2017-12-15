package main

import (
	"fmt"
	"os"

	"github.com/TomOnTime/tomutils/cidrtools/ipimbo"
)

const (
	INVALID  = iota // List is not sorted.
	CONTAINS = iota // b is contained in a.
	DISJOINT = iota // b is not contained in a.
)

func process(last, next ipimbo.Imbo, f func(ipimbo.Imbo)) ipimbo.Imbo {
	decision, err := judge(last, next)
	//fmt.Printf("%v, %v := judge(%s, %s)\n", decision, err, last.String(), next.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
	switch decision {
	case CONTAINS:
		// Skip next.
		return last
	case DISJOINT:
		// b is the new a.
		f(last)
		return next
	default:
		panic("Should not happen (decision)")
	}
}

func judge(a, b ipimbo.Imbo) (int, error) {

	// List was not sorted (list is assumed to be sorted)
	//   10.10.10.0/32
	//   10.10.10.0/24
	//   -> INVALID
	if a.Less(b) {
		return INVALID, fmt.Errorf("List was not sorted: %s, %s", a, b)
	}

	// At the boundary between IPv4 and IPv6, always start fresh.
	//   10.10.10.10/24
	//   ::1
	//   -> keep B
	if a.Version() != b.Version() {
		return DISJOINT, nil
	}

	// If "host portion" is not all zeros, this is an error.
	//   (anything)
	//   10.10.10.9/24
	//   -> INVALID
	if !b.HostAllZeros() {
		return INVALID, fmt.Errorf("The host bits are not all zero: %s", b)
	}

	// If b is contained in a, discard b.
	//   10.10.10.0/24
	//   10.10.10.0/32
	//   -> discard b
	if a.PrefixLen() > b.PrefixLen() {
		return DISJOINT, nil
	}
	if a.Contains(b) {
		return CONTAINS, nil
	}

	//   1.10.10.0/32
	//   10.10.10.0/24
	//   -> keep b
	return DISJOINT, nil
}
