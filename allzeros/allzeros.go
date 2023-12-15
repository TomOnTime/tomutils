package main

/*

allzeros takes a list of filenames and outputs the ones where the
entire contents of the file is zero bytes.

This was written when I discovered a disk error had lost a lot of data
and turned many files of the same length, but the contents was all
zero bytes.

Usage:

    allzeros * >/tmp/zeros.txt

*/

import (
	"bufio"
	"fmt"
	"os"
)

func isAllZero(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Use a buffered reader to read the file in chunks
	reader := bufio.NewReader(file)

	// Check if all bytes are zero
	for {
		b, err := reader.ReadByte()
		if err != nil {
			// EOF reached
			break
		}

		if b != 0 {
			return false, nil
		}
	}

	return true, nil
}

func main() {
	// Check if at least one filename is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go file1 file2 ...")
		os.Exit(1)
	}

	// Process each filename provided as a command line argument
	for _, filename := range os.Args[1:] {
		isZero, err := isAllZero(filename)
		if err != nil {
			fmt.Printf("Error checking file %s: %v\n", filename, err)
			continue
		}

		// Print the filename if all bytes are zero
		if isZero {
			fmt.Println(filename)
		}
	}
}
