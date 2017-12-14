package ipimbo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func Parse(r io.Reader) <-chan Imbo {
	c := make(chan Imbo)
	go func() {
		defer close(c)

		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			c <- strings.ToUpper(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		// c <- i
	}()
	return c
}

func Parse(r io.Reader) <-chan Imbo {
	fh.Close() // Needed?

}
