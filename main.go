package main

import (
	"fmt"
	"licensezero.com/licensezero/subcommands"
	"os"
	"sort"
)

// Rev represents the current build revision.  Set via ldflags.
var Rev string

var commands = map[string]*subcommands.Subcommand{
	"backup":   subcommands.Backup,
	"bugs":     subcommands.Bugs,
	"identify": subcommands.Identify,
	"import":   subcommands.Import,
	"latest":   subcommands.Latest,
	"version":  subcommands.Version,
	"whoami":   subcommands.WhoAmI,
}

func main() {
	paths := subcommands.NewPaths()
	arguments := os.Args
	if len(arguments) > 1 {
		subcommand := os.Args[1]
		if value, ok := commands[subcommand]; ok {
			if subcommand == "version" || subcommand == "latest" {
				value.Handler([]string{Rev}, paths)
			} else {
				value.Handler(os.Args[2:], paths)
			}
		} else {
			showUsage()
			os.Exit(1)
		}
	} else {
		showUsage()
		os.Exit(0)
	}
}

func showUsage() {
	os.Stdout.WriteString("Manage License Zero projects and dependencies.\n\nSubcommands:\n")
	buyer := map[string]*subcommands.Subcommand{}
	seller := map[string]*subcommands.Subcommand{}
	misc := map[string]*subcommands.Subcommand{}
	for key, value := range commands {
		switch value.Tag {
		case "buyer":
			buyer[key] = value
		case "seller":
			seller[key] = value
		default:
			misc[key] = value
		}
	}
	listSubcommands("For Buyers", buyer)
	listSubcommands("For Sellers", seller)
	listSubcommands("Miscellaneous", misc)
}

func listSubcommands(header string, list map[string]*subcommands.Subcommand) {
	os.Stdout.WriteString("\n  " + header + ":\n\n")
	longestSubcommand := 0
	var names []string
	for name := range list {
		length := len(name) + 1
		if length > longestSubcommand {
			longestSubcommand = length
		}
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		info := list[name]
		fmt.Printf("  %-"+fmt.Sprintf("%d", longestSubcommand)+"s %s\n", name, info.Description)
	}
}
