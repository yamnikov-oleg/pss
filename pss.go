package main

import "os"

var commands []*Subcommand

func registerSubcommand(sub *Subcommand) {
	commands = append(commands, sub)
}

func main() {
	if !DispatchSubcommand(os.Args, commands) {
		os.Exit(1)
	}
}
