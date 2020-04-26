package cmd

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/core"
	"os"
)

// Package cmd define the command line instructions

type cmdsMap map[string]func(string) error	// Store which func for each cmd

// parseArg does two tasks:
//		- parse command parameters
//		- call the functions in core package
func parseArg() error{
	numArgs := len(os.Args)
	if numArgs <= 1 {
		fmt.Printf("Not enough parameter.\n")
		return nil
	}

	// call functions in core
	switch os.Args[1] {
	case "init":
		if numArgs <= 2{
			fmt.Printf("Not enough parameter.\n")
			return nil
		}
		if os.Args[2] == "help" {
			fmt.Printf("shadow init <repo path>\n")
			return nil
		}
		cfg := core.InitConfig{RepoPath:os.Args[2]}
		return core.InitRepo(cfg)
	case "start":
		if numArgs <= 2{
			fmt.Printf("Not enough parameter.\n")
			return nil
		}
		if os.Args[2] == "help" {
			fmt.Printf("shadow start <repo path>\n")
			return nil
		}
		//cfg := core.InitConfig{RepoPath:os.Args[2]}
		//return core.InitRepo(cfg)

	default:
		fmt.Printf("Unknown command %s\n", os.Args[1])
	}
	return nil
}

// Run() start the software
func Run() error {
    //TODO: do not use template
	cmds := make(cmdsMap)
	cmds["init"] = func(path string) error {
		cfg := core.InitConfig{RepoPath: path}
		return core.InitRepo(cfg)
	}
	cmd := os.Args[1]
    path := os.Args[2]
	for key, value := range cmds {
		if key == cmd {
			return value(path)
		}
	}

	return nil
}
