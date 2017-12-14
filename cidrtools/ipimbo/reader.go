package ipimbo

import (
	"bufio"
	"fmt"
	"io"
)

func Parse(r io.Reader) <-chan Imbo {
	c := make(chan Imbo)
	go func() {
		defer close(c)

		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line = scanner.Text()
			fmt.Println(t)
			c <- t
		}

		if err := scanner.Err(); err != nil {
			return nil
		}

	}()

	return c
}
