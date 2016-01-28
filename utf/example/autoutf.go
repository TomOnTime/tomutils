package main

// Functions that autodetect UTF-8 and UTF-16/LE/BE but return UTF-8.
// The goal is to create functions that "just do the right thing"
// no matter what UTF encoding is used.
//
// OpenUTF: Similar to os.Open()
// ReadFileUTF: Similar to ioutil.ReadFile()
// NewReaderUTF: Similar to NewReader() in other libraries.

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

//func OpenUTF(filename string) (io.Reader, error) {
//	f, err := os.Open(filename)
//	if err != nil {
//		return nil, err
//	}
//	return NewReaderUTF(f)
//}

func NewReaderUTF(r io.Reader) io.Reader {
	// Make an tranformer that converts MS-Windows (16LE) to UTF8:
	winutf := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// Make a transformer that is like winutf, but abides by BOM if found:
	utf16bom := unicode.BOMOverride(winutf.NewDecoder())

	// Make a transformer that is like , but abides by BOM if found:
	//utf16bom := unicode.BOMOverride(unicode.UTF8)

	// This technique is recommended by the W3C for use in HTML 5:
	// "For compatibility with deployed content, the byte order
	// mark (also known as BOM) is considered more authoritative
	// than anything else." http://www.w3.org/TR/encoding/#specification-hooks

	// Make a Reader that uses utf16bom:
	return transform.NewReader(r, utf16bom)
}

func ReadFileUTF(filename string) ([]byte, error) {

	// decode and print:
	decoded, err := ioutil.ReadAll(NewReaderUTF(filename))
	return decoded, err
}

func main() {
	data, err := ReadFileUTF("inputfile.txt")
	if err != nil {
		log.Fatal(err)
	}
	final := strings.Replace(string(data), "\r\n", "\n", -1)
	fmt.Println(final)

}
