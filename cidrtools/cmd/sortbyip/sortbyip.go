package main

import (
	"flag"
	"fmt"

	"github.com/TomOnTime/tomutils/cidrtools/ipimbo"
)

func main() {
	flag.Parse()

	var ipdb ipimbo.ImboList

	ipimboHandle := ipimbo.New()
	ipimbo.IterateArgs(ipimboHandle, func(v ipimbo.Imbo) {
		ipdb = append(ipdb, v)
	})

	// Sort
	ipdb.Sort()

	// Print
	for _, im := range ipdb {
		fmt.Println(im.String())
	}
}
