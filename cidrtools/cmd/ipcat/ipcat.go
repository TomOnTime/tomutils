package main

import (
	"flag"
	"fmt"

	"github.com/TomOnTime/tomutils/cidrtools/ipimbo"
)

func main() {
	flag.Parse()

	ipimboHandle := ipimbo.New()
	ipimbo.IterateArgs(ipimboHandle, func(v ipimbo.Imbo) {
		fmt.Println(v.String())
	})
}
