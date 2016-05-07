package main

import (
	"fmt"
	"os"

	"github.com/yamnikov-oleg/pss/lib/pss"
)

func init() {
	registerSubcommand(&Subcommand{
		Name:    "init",
		Hint:    "Make a new password storage.",
		Handler: cmdInit,
	})
}

func cmdInit([]string) bool {
	_, err := os.Stat(pss.StoragePath)
	if err == nil || !os.IsNotExist(err) {
		fmt.Println("Error: file already exists.")
		fmt.Printf("Please remove the storage at %q before creating a new one\n", pss.StoragePath)
		return false
	}

	pwd := promtPwd("Enter master-password", false)
	if !saveStorage(pss.Storage{}, pwd) {
		return false
	}
	fmt.Printf("Successfully created new password storage at %v\n", pss.StoragePath)
	return true
}
