package main

import (
	"encoding/csv"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/TomOnTime/tomutils/vidnamer/filehash"
	"github.com/TomOnTime/tomutils/vidnamer/filminventory"
)

/*

1. Read a directory of files
2. For each file:
    ParseFilename
		generate CSV line.

*/

func gatherFilenames(names []string) []string {
	//fmt.Printf("DEBUG: gatherFilenames called with %+v\n", names)
	var result []string

	for _, root := range names {

		fileSystem := os.DirFS(root)

		fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if d.IsDir() {
				return nil // Skip directories
			}
			if strings.HasPrefix(d.Name(), ".") {
				return nil // Skip "dot files"
			}

			//fmt.Printf("DEBUG: WalkDir: %q\n", path)
			result = append(result, path)

			return nil
		})
	}

	return result
}

func parseHeaders(line string) []string {
	return strings.Split(line, ",")
}

func writeHeader(w *csv.Writer, headers []string) {
	w.Write(headers)
}

func writeData(w *csv.Writer, filmDB []filminventory.Film, sigDB *filehash.DB) {

	for _, film := range filmDB {
		//fmt.Printf("% 4d: %q\n", i, film.Title)
		items := film.ToCSV(sigDB)
		//fmt.Printf("DEBUG: items=%+v\n", items)
		w.Write(items)
	}
}

func main() {

	// Read the hash file
	sigDB, err := filehash.Initialize("../sha256.list")
	if err != nil {
		log.Fatalf("Can't filehash.Initalize(): %v", err)
	}

	w := csv.NewWriter(os.Stdout)

	headers := parseHeaders(filminventory.CSVHeader)
	writeHeader(w, headers)

	fileNames := gatherFilenames(os.Args[1:])
	//fmt.Printf("DEBUG: fileNames=%+v\n", fileNames)
	filmDB := filminventory.FromManyFilepaths(fileNames, sigDB)
	writeData(w, filmDB, sigDB)

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

}
