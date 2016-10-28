package main

import "fmt"

func init() {
	registerSubcommand(&Subcommand{
		Name:    "modify",
		Aliases: []string{"mod", "m", "update", "upd", "u"},
		Usage:   "<website> [username]",
		Hint:    "Modify a record in the storage.",
		Handler: cmdModify,
	})
}

func cmdModify(args []string) bool {
	mod, stg, masterpwd, ok := openAndFindRecord(args)
	if !ok {
		return false
	}

	fmt.Println("Modifying record for:")
	printRecord(mod)

	var (
		website  string
		username string
	)

	fmt.Printf("New value for website (leave blank to keep old value: %v):\n", mod.Website)
	if _, err := fmt.Scanln(&website); err == nil {
		mod.Website = website
	}

	fmt.Printf("New value for username (leave blank to keep old value: %v):\n", mod.Username)
	if _, err := fmt.Scanln(&username); err == nil {
		mod.Username = username
	}

	pwd := promtPwd("New value for password (leave blank to keep old value)", true)
	if pwd != "" {
		mod.Password = pwd
	}

	if !saveStorage(stg, masterpwd) {
		return false
	}

	fmt.Println("Record has been successfully updated:")
	printRecord(mod)
	return true
}
