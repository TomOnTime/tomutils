package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	iterateArgs(func(s string) {
		s = strings.TrimSpace(s)
		if s != "" && s[0] != '#' {
			fmt.Println(s)
		}
	})
}

func iterateArgs(f func(s string)) {
	var fh *os.File
	var err error

	if flag.NArg() == 0 {
		for v := range readFile(os.Stdin) {
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
			for v := range readFile(fh) {
				f(v)
			}
			fh.Close()
		}
	}
}

func readFile(r io.Reader) <-chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			c <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}()
	return c
}
