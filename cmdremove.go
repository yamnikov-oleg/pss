package main

import "fmt"

func init() {
	registerSubcommand(&Subcommand{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "<website> [username]",
		Hint:    "Remove specific record from storage.",
		Handler: cmdRemove,
	})
}

func cmdRemove(args []string) bool {
	rem, stg, masterpwd, ok := openAndFindRecord(args)
	if !ok {
		return false
	}

	var ans string
	printRecord(rem)
	fmt.Println("Are you sure you want to delete this record?")
	fmt.Println("This action cannot be rolled back.")
	fmt.Print("(y/N): ")
	fmt.Scanln(&ans)

	if ans != "y" && ans != "Y" {
		fmt.Println("Deletion has been canceled.")
		return true
	}

	var ind = -1
	for i := range stg {
		if stg[i] == rem {
			ind = i
			break
		}
	}
	stg = append(stg[:ind], stg[ind+1:]...)
	if !saveStorage(stg, masterpwd) {
		return false
	}
	fmt.Println("Record has been deleted.")
	return true
}
