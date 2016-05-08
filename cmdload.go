package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/yamnikov-oleg/pss/lib/pss"
)

func init() {
	registerSubcommand(&Subcommand{
		Name:    "load",
		Usage:   "<path>",
		Hint:    "Load dump from json file into the storage.",
		Handler: cmdLoad,
	})
}

func cmdLoad(args []string) bool {
	if len(args) != 1 {
		fmt.Println("Please provide single filepath via command line arguments.")
		return false
	}

	file, err := os.Open(args[0])
	if err != nil {
		fmt.Println("Error opening file:")
		fmt.Println(err)
		return false
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:")
		fmt.Println(err)
		return false
	}

	var loaded pss.Storage
	if err := json.Unmarshal(buf, &loaded); err != nil {
		fmt.Println("Error parsing data:")
		fmt.Println(err)
		return false
	}

	stg, pwd, ok := openStorage()
	if !ok {
		return false
	}

	stg = append(stg, loaded...)
	sort.Sort(byWebsite(stg))
	if !saveStorage(stg, pwd) {
		return false
	}

	fmt.Println("Successfully loaded:")
	for _, rec := range loaded {
		printRecord(rec)
	}
	return true
}
