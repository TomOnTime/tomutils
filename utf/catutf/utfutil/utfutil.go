// Package ioutil implements some I/O utility functions that
// are UTF-encoding agnostic.
package utfutil

// These functions autodetect UTF-8/UTF-16/LE/BE and return UTF-8.
// You can use them as replacements for os.Open() and ioutil.ReadFile()
// when the encoding of the file is unknown.

// This technique is recommended by the W3C for use in HTML 5:
// "For compatibility with deployed content, the byte order
// mark (also known as BOM) is considered more authoritative
// than anything else." http://www.w3.org/TR/encoding/#specification-hooks

import (
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func UTFNewReader(rd io.Reader) io.Reader {
	// Make an tranformer that converts MS-Windows (16LE) to UTF8:
	winutf := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// Make a transformer that is like winutf, but abides by BOM if found:
	utf16bom := unicode.BOMOverride(winutf.NewDecoder())

	// Make a transformer that is like , but abides by BOM if found:
	//utf16bom := unicode.BOMOverride(unicode.UTF8)

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(rd, utf16bom)
	return unicodeReader
}

// OpenUTF is the equivalent of os.Open().
func Open(name string) (io.Reader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return UTFNewReader(f), nil
}

func ReadFile(name string) ([]byte, error) {
	file, err := Open(name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}
