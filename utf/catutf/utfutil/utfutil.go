// Package ioutil implements some I/O utility functions that
// are UTF-encoding agnostic.
package utfutil

// These functions autodetect UTF BOM and return UTF-8. If no
// BOM is found, they assume UTF-8 (Linux) or UTF-16LE (Windows).
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

type Assumption int

const (
	UTF8 Assumption = iota
	UTF16LE
	UTF16BE
	WINDOWS = UTF16LE
	POSIX   = UTF8
	HTML5   = UTF8
)

func UTFNewReader(rd io.Reader, ume Assumption) io.Reader {
	var trans transform.Transformer
	switch ume {
	case UTF8:
		// Make a transformer that assumes UTF-8 but abides by the BOM.
		trans = unicode.BOMOverride(unicode.UTF8.NewDecoder())
	case UTF16LE:
		// Make an tranformer that decodes MS-Windows (16LE) UTF files:
		winutf := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		// Make a transformer that is like winutf, but abides by BOM if found:
		trans = unicode.BOMOverride(winutf.NewDecoder())
	case UTF16BE:
		// Make an tranformer that decodes UTF-16BE files:
		utf16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		// Make a transformer that is like utf16be, but abides by BOM if found:
		trans = unicode.BOMOverride(utf16be.NewDecoder())
	}

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(rd, trans)
	return unicodeReader
}

// OpenUTF is the equivalent of os.Open().
func Open(name string, ume Assumption) (io.Reader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return UTFNewReader(f, ume), nil
}

// ReadFile is the equivalent of ioutil.ReadFile()
func ReadFile(name string, ume Assumption) ([]byte, error) {
	file, err := Open(name, ume)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}
