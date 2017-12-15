package ipimbo

import (
	"flag"
	"fmt"
	"os"
)

// args is a helper function that calls f() on every line of the files
// on the command line.
func IterateArgs(ipdb *Handle, f func(imdb Imbo)) {
	var fh *os.File
	var err error

	if flag.NArg() == 0 {
		for v := range ipdb.ReadFile(os.Stdin) {
			f(v)
		}
	} else {
		for _, fname := range flag.Args() {
			fh, err = os.Open(fname)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading file: ", err)
				os.Exit(1)
			}
			defer fh.Close()

			for v := range ipdb.ReadFile(fh) {
				f(v)
			}

			fh.Close() // Needed?

		}
	}

}
