package main

import (
	"flag"
	"fmt"
	"ipimbo"
	"os"
)

func main() {

	flag.Parse()

	var fh *os.File
	var err error

	ipdb := ipimbo.New()

	if flag.NArg() == 0 {
		for v := range ipdb.Parse(os.Stdin) {
			fmt.Println(v)
		}
	} else {
		for _, fname := range flag.Args() {
			fh, err = os.Open(fname)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading file: ", err)
				os.Exit(1)
			}
			defer fh.Close()

			for v := range ipdb.Parse(fh) {
				fmt.Println(v)
			}

			fh.Close() // Needed?

		}
	}

}
