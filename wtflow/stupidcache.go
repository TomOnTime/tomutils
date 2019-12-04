package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

var CACHEFILE = "/tmp/dnscache.gob"

// Implement a very stupid cache.
var stupidcache = map[string]string{}

// https://socketloop.com/tutorials/golang-saving-and-reading-file-with-gob

func stupidBegin(fname string) {

	// open data file
	dataFile, err := os.Open(fname)
	if err != nil {
		fmt.Printf("FYI: Cache file couldn't be read: %v\n", err)
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&stupidcache)
	if err != nil {
		fmt.Printf("FYI: Cache couldn't be decoded: %v\n", err)
	}

	dataFile.Close()
}

func stupidEnd(fname string) {
	// create a file
	dataFile, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("FYI: Cache file couldn't be written: %v\n", err)
	}

	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(stupidcache)

	dataFile.Close()
}
