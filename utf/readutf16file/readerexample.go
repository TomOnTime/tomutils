package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Similar to ioutil.ReadFile() but decodes UTF-16.  Useful when
// reading data from MS-Windows systems that generate UTF-16BE files,
// but will do the right thing if other BOMs are found.
func ReadFileUTF16(filename string) ([]byte, error) {

	// Read the file into a []byte:
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16bom)

	// decode and print:
	decoded, err := ioutil.ReadAll(unicodeReader)
	return decoded, err
}

func main() {
	data, err := ReadFileUTF16("inputfile.txt")
	if err != nil {
		log.Fatal(err)
	}
	final := strings.Replace(string(data), "\r\n", "\n", -1)
	fmt.Println(final)

}
