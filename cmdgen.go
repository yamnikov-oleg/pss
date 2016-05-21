package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/yamnikov-oleg/pss/Godeps/_workspace/src/github.com/atotto/clipboard"
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

	if err := clipboard.WriteAll(string(buf)); err != nil {
		fmt.Println("Error accessing clipboard:")
		fmt.Println(err)
		return false
	}

	fmt.Println("Generated password has been copied to your clipboard.")
	fmt.Println("Use Ctrl-V or 'Paste' command to use it.")
	return true
}
