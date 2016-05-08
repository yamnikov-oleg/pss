package main

import (
	"fmt"

	"github.com/yamnikov-oleg/pss/Godeps/_workspace/src/github.com/atotto/clipboard"
)

func main() {
	text, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}

	fmt.Print(text)
}
