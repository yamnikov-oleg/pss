package main

import (
	"os"
	"testing"
)

func TestUsage(t *testing.T) {
	testCases := map[string]struct {
		cmd  string
		subs []*Subcommand
		text string
	}{
		"No subcommands": {
			"mycmd",
			nil,
			`Usage: mycmd
`,
		},
		"Some commands": {
			"go",
			[]*Subcommand{
				&Subcommand{
					"run", []string{"r", "rn"},
					"[file.go] <options>",
					"Run specific source file",
					nil,
				},
				&Subcommand{
					"build", nil,
					"",
					"Build current package",
					nil,
				},
			},
			`Usage: go <command>

Available commands:

build       Build current package
            Usage: build

run         Run specific source file
            Alises: r, rn
            Usage: run [file.go] <options>

`,
		},
	}

	for name, cas := range testCases {
		act := Usage(cas.cmd, cas.subs)
		if act != cas.text {
			t.Errorf("%v. Ожидалось:\n%vПолучено:\n%v", name, cas.text, act)
		}
	}
}

func TestDispatchSubcommand(t *testing.T) {
	// Actual arguments after a call
	var aargs []string
	// Make a handler which save the args and return ret.
	handler := func(ret bool) func(args []string) bool {
		return func(args []string) bool {
			aargs = args
			return ret
		}
	}
	// Suppress printing
	os.Stderr, _ = os.Open(os.DevNull)

	testCases := map[string]struct {
		args  []string
		subs  []*Subcommand
		eret  bool
		eargs []string
	}{
		"No cmds": {
			[]string{"cmd", "sub", "arg1", "arg2"},
			nil,
			false,
			nil,
		},
		"1 cmd": {
			[]string{"cmd", "sub", "arg1", "arg2"},
			[]*Subcommand{
				&Subcommand{
					Name:    "sub",
					Aliases: []string{"s", "ss", "sb"},
					Handler: handler(true),
				},
			},
			true,
			[]string{"arg1", "arg2"},
		},
		"3 cmd": {
			[]string{"go", "build", "source.go"},
			[]*Subcommand{
				&Subcommand{
					Name:    "build",
					Handler: handler(true),
				},
				&Subcommand{
					Name:    "run",
					Aliases: []string{"r"},
					Handler: handler(true),
				},
				&Subcommand{
					Name:    "test",
					Aliases: []string{"t", "tst"},
					Handler: handler(true),
				},
			},
			true,
			[]string{"source.go"},
		},
		"3 cmd, no args": {
			[]string{"go", "build"},
			[]*Subcommand{
				&Subcommand{
					Name:    "build",
					Handler: handler(true),
				},
				&Subcommand{
					Name:    "run",
					Aliases: []string{"r"},
					Handler: handler(true),
				},
				&Subcommand{
					Name:    "test",
					Aliases: []string{"t", "tst"},
					Handler: handler(true),
				},
			},
			true,
			nil,
		},
		"no subcmd supplied": {
			[]string{"go"},
			[]*Subcommand{
				&Subcommand{
					Name:    "build",
					Handler: handler(true),
				},
				&Subcommand{
					Name:    "run",
					Aliases: []string{"r"},
					Handler: handler(true),
				},
			},
			false,
			nil,
		},
		"cmd ret false": {
			[]string{"go", "build"},
			[]*Subcommand{
				&Subcommand{
					Name:    "build",
					Handler: handler(false),
				},
				&Subcommand{
					Name:    "run",
					Aliases: []string{"r"},
					Handler: handler(true),
				},
			},
			false,
			nil,
		},
		"zero args": {
			nil,
			[]*Subcommand{
				&Subcommand{
					Name:    "build",
					Handler: handler(true),
				},
				&Subcommand{
					Name:    "run",
					Aliases: []string{"r"},
					Handler: handler(true),
				},
			},
			false,
			nil,
		},
	}

	for name, cas := range testCases {
		aargs = nil
		ret := DispatchSubcommand(cas.args, cas.subs)
		if ret != cas.eret {
			t.Errorf("%v: возвращено %v, ожидалось %v", name, ret, cas.eret)
		}
		if len(aargs) != len(cas.eargs) {
			t.Errorf("%v: длины аргументов не совпадают: %v и %v",
				name, len(aargs), len(cas.eargs))
		}
		for i := range aargs {
			if aargs[i] != cas.eargs[i] {
				t.Errorf("%v: аргумент %d не совпадает: %v и %v",
					name, i, aargs[i], cas.eargs[i])
			}
		}
	}
}
