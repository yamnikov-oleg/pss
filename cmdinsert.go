package main

import (
	"fmt"
	"sort"

	"github.com/yamnikov-oleg/pss/lib/pss"
)

func init() {
	registerSubcommand(&Subcommand{
		Name:    "insert",
		Aliases: []string{"ins", "i"},
		Usage:   "<username> <website>",
		Hint:    "Put a new password record into the storage.",
		Handler: cmdInsert,
	})
}

type byWebsite pss.Storage

func (l byWebsite) Len() int {
	return len(l)
}

func (l byWebsite) Less(i, j int) bool {
	return l[i].Website < l[j].Website
}

func (l byWebsite) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func cmdInsert(args []string) bool {
	if len(args) != 2 {
		fmt.Println("Please, supply website and username via command line.")
		return false
	}
	uname, website := args[0], args[1]

	stg, masterpwd, ok := openStorage()
	if !ok {
		return false
	}

	pwd := promtPwd(fmt.Sprintf("Password for %v at %v", uname, website), false)

	stg = append(stg, &pss.Record{
		Website:  website,
		Username: uname,
		Password: pwd,
	})
	sort.Sort(byWebsite(stg))

	if !saveStorage(stg, masterpwd) {
		return false
	}

	fmt.Println("Successfully saved the record")
	return true
}
