package main

import "testing"

func TestSubcmdAliasDuplicates(t *testing.T) {
	nmap := map[string]*Subcommand{}
	for _, cmd := range commands {
		names := []string{cmd.Name}
		names = append(names, cmd.Aliases...)
		for _, n := range names {
			if _, ok := nmap[n]; ok {
				t.Errorf("Name %q is taken by both command %v and %v",
					n, cmd.Name, nmap[n].Name)
				continue
			}
			nmap[n] = cmd
		}
	}
}
