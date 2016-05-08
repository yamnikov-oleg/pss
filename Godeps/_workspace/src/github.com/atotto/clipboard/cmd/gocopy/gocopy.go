package main

import (
	"io/ioutil"
	"os"

	"github.com/yamnikov-oleg/pss/Godeps/_workspace/src/github.com/atotto/clipboard"
)

func main() {

	out, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	if err := clipboard.WriteAll(string(out)); err != nil {
		panic(err)
	}
}
