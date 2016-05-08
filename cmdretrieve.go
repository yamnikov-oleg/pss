package main

import (
	"fmt"

	"github.com/atotto/clipboard"
)

func init() {
	registerSubcommand(&Subcommand{
		Name:    "retrieve",
		Aliases: []string{"r", "checkout", "co"},
		Usage:   "<website> [username]",
		Hint:    "Load a password from storage to clipboard",
		Handler: cmdRetrieve,
	})
}

func cmdRetrieve(args []string) bool {
	rec, _, _, ok := openAndFindRecord(args)
	if !ok {
		return false
	}

	if err := clipboard.WriteAll(rec.Password); err != nil {
		fmt.Println("Error accessing clipboard:")
		fmt.Println(err)
		return false
	}
	fmt.Println("Password for: ")
	printRecord(rec)
	fmt.Println("has been copied to your clipboard.")
	fmt.Println("Use Ctrl-V or 'Paste' command to use it.")
	return true
}
