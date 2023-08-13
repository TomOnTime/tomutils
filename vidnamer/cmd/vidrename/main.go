package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TomOnTime/tomutils/vidnamer/filehash"
	"github.com/TomOnTime/tomutils/vidnamer/filminventory"
	"github.com/gocarina/gocsv"
)

func main() {

	// Read sha256 into SigDB

	// Read the hash file
	sigDB, err := filehash.Initialize("../sha256.list")
	if err != nil {
		log.Fatalf("Can't filehash.Initalize(): %v", err)
	}

	// Read CSV one line at a time
	clientsFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	clients := []*filminventory.FilmCSV{}

	if err := gocsv.UnmarshalFile(clientsFile, &clients); err != nil { // Load clients from file
		panic(err)
	}
	for _, client := range clients {
		film := filminventory.CsvToFilm(*client)

		sigFilename := sigDB.SigToFilename(client.Sha256)

		existing := sigFilename
		desired := film.DesiredFilename()
		if existing == desired {
			fmt.Printf("# GOOD: %q\n", existing)
		} else {
			fmt.Printf("mv \\\n    %q\\\n    %q\n", existing, desired)
		}
	}
}
