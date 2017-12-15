package ipimbo

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func (*Handle) ReadFile(r io.Reader) <-chan Imbo {
	c := make(chan Imbo)
	go func() {
		defer close(c)

		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			//  c <- parseline(scanner.Text())
			im, err := parseline(scanner.Text())
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			c <- im
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}()
	return c
}
