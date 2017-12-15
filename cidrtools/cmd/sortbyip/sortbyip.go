package main

import (
	"flag"
	"fmt"

	"github.com/TomOnTime/tomutils/cidrtools/ipimbo"
)

func main() {
	flag.Parse()

	ipimboHandle := ipimbo.New()
	var ipdb ipimbo.ImboList
	ipimbo.IterateArgs(ipimboHandle, func(v ipimbo.Imbo) {
		ipdb = append(ipdb, v)
	})

	fmt.Println("BEFORE:")
	for _, im := range ipdb {
		fmt.Println(im.DebugString())
	}

	// Sort
	ipdb.Sort()

	// Print
	fmt.Println("AFTER:")
	for _, im := range ipdb {
		fmt.Println(im.DebugString())
	}

}
