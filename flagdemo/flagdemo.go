package main

import (
	"fmt"

	"github.com/TomOnTime/tomutils/flag"
)

var fa = flag.Int("a", 1234, "a help")
var fb = flag.Int("b", 66, "b help")
var fs = flag.String("s", "default", "s help")

func main() {
	flag.DefaultsFromFiles("also.flags", "also2.flags")
	flag.DefaultsFromFiles("also3.flags")
	flag.Parse()

	fmt.Printf("a=%d b=%d s=%#v\n", *fa, *fb, *fs)

}
