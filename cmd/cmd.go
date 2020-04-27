package cmd

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/core"
	"os"
)

// Package cmd define the command line instructions


// Run() start the software
//		- parse command parameters
//		- call the functions in core package
func Run() error {
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
