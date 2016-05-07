package main

import (
	"fmt"
	"strings"
)

func init() {
	registerSubcommand(&Subcommand{
		Name:    "list",
		Aliases: []string{"ls", "l"},
		Usage:   "[search query]",
		Hint:    "Print list of saved records, filtred by optional search query.",
		Handler: cmdList,
	})
}

func cmdList(args []string) bool {
	stg, _, ok := openStorage()
	if !ok {
		return false
	}

	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	fmt.Println("Contents of the storage:")
	for _, r := range stg {
		if strings.Contains(r.Website, query) || strings.Contains(r.Username, query) {
			printRecord(r)
		}
	}
	return true
}
