package filehash

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Info struct {
	Filename  string
	Signature string
}

func FromFile(filename string) (r []Info, err error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		before, after, found := strings.Cut(line, " ")
		if !found {
			// print "SKIPPING"
			fmt.Printf("SKIPPING: %v\n", line)
		} else {
			r = append(r, Info{Signature: before, Filename: after})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return r, nil
}
