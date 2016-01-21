package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type utfScanner interface {
	Read(p []byte) (n int, err error)
}

// Creates a scanner similar to os.Open() but decodes the file as UTF-16.
// Useful when reading data from MS-Windows systems that generate UTF-16BE
// files, but will do the right thing if other BOMs are found.
func NewScannerUTF16(filename string) (utfScanner, error) {

	// Read the file into a []byte:
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(file, utf16bom)
	return unicodeReader, nil
}

func main() {

	s, err := NewScannerUTF16("inputfile.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading inputfile:", err)
	}

}
