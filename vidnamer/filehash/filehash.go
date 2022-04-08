package filehash

import (
	"bufio"
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
		parts := strings.Cut(line, " ")
		if len(parts) != 2 {
			// print "SKIPPING"
		} else {
			r = append(r, &Info{Signature: parts[0], Filename: parts[1]})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return r, nil
}
