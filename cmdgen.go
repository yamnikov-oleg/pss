package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	registerSubcommand(&Subcommand{
		Name:    "gen",
		Aliases: []string{"g"},
		Usage:   "[length]",
		Hint:    "Generate random password of given length (default 12).",
		Handler: cmdGen,
	})
}

func cmdGen(args []string) bool {
	var lng = 12
	if len(args) > 0 {
		// We don't care if it fails. If so, it just keeps the default value.
		fmt.Sscan(args[0], &lng)
	}

	buf := make([]byte, lng)
	rand.Seed(time.Now().UnixNano())
	for i := range buf {
		buf[i] = byte(rand.Intn(26 + 26 + 10))
		switch {
		case buf[i] < 10:
			buf[i] += '0' // Digit
		case buf[i] < 10+26:
			buf[i] += 'A' - 10 // Uppercase letter
		default:
			buf[i] += 'a' - 26 - 10 // Lowercase letter
		}
	}

	fmt.Println(string(buf))
	return true
}
