package main

import (
	"flag"
	"fmt"

	"github.com/TomOnTime/tomutils/cidrtools/ipimbo"
)

func main() {
	flag.Parse()

	ipimboHandle := ipimbo.New()
	var first bool = true
	var head ipimbo.Imbo
	ipimbo.IterateArgs(ipimboHandle, func(v ipimbo.Imbo) {
		if first {
			head = v
			first = false
		} else {
			head = process(
				head,
				v,
				func(v ipimbo.Imbo) {
					fmt.Println(v.String())
				})
		}
	})
	fmt.Println(head.String())
}
