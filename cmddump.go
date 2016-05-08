package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func init() {
	registerSubcommand(&Subcommand{
		Name:    "dump",
		Usage:   "<path>",
		Hint:    "Dump storage contents into json file.",
		Handler: cmdDump,
	})
}

func cmdDump(args []string) bool {
	if len(args) != 1 {
		fmt.Println("Please provide single filepath via command line arguments.")
		return false
	}

	file, err := os.Create(args[0])
	if err != nil {
		fmt.Println("Error creating the file:")
		fmt.Println(err)
		return false
	}
	defer file.Close()

	stg, _, ok := openStorage()
	if !ok {
		return false
	}

	buf, err := json.Marshal(stg)
	if err != nil {
		fmt.Println("Error marshalling the storage:")
		fmt.Println(err)
		return false
	}

	if _, err := file.Write(buf); err != nil {
		fmt.Println("Error writing file:")
		fmt.Println(err)
		return false
	}

	fmt.Printf("Successfully written the dump to %q\n", args[0])
	return true
}
