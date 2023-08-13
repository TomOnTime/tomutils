package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	// Read everything into memory.
	contentsBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	contents := string(contentsBytes)

	// Find all the filenames
	isZero := func(c rune) bool {
		return c == '\000'
	}
	fileNames := strings.FieldsFunc(contents, isZero)

	// Analyze and find all the directories:
	isDir := map[string]bool{}
	for _, name := range fileNames {
		d := filepath.Dir(name)
		isDir[d] = true
	}

	// BUG: Empty directories are listed as isDir == false.

	// scan the list, outputting each item in all.txt, files.txt or dirs.txt
	ao, err := os.Create("all.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer ao.Close()

	fo, err := os.Create("files.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fo.Close()

	fo0, err := os.Create("files0.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fo0.Close()

	do, err := os.Create("dirs.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer do.Close()

	do0, err := os.Create("dirs0.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer do0.Close()

	for _, name := range fileNames {

		fmt.Fprintln(ao, name)

		if isDir[name] {
			fmt.Fprintln(do, name)

			fmt.Fprint(do0, name)
			fmt.Fprint(do0, "\000")

		} else {
			fmt.Fprintln(fo, name)

			fmt.Fprint(fo0, name)
			fmt.Fprint(fo0, "\000")

		}
	}

}
