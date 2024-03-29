package chat

import (
	"fmt"
	"io"
)

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	go io.Copy(a, b)
	io.Copy(b, a)
}

var partner = make(chan io.ReadWriteCloser)

func Match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")

	select {
	case partner <- c:
		fmt.Fprint(c, "second partner")
		// now handled by the other goroutine
	case p := <-partner:
		fmt.Fprint(c, "first partner")
		chat(p, c)
	}

}
