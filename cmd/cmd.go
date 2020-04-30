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
        cfg := core.RunConfig{
            RepoPath: os.Args[2],
        }
        textile, err := core.NewTextile(cfg)
        if err != nil {
            return err
        }
        textile.Start()
		//err = NewServer().Listen()
		return err
        // hang forever
/*
	case "connect":
		if numArgs <= 2 {
			fmt.Printf("Not enough parameter\n")
			return nil
		}
		if os.Args[2] == "help" {
			fmt.Printf("shadow connect <peer id> <address>")
			return nil
		} else if numArgs <= 3 {
			fmt.Printf("Not enough parameter\n")
			return nil
		}

		peerId := os.Args[2]
		address := os.Args[3]
		err := NewApi().Call(peerId + " " + address)
		if err != nil {
			fmt.Printf("Error occurs when calling api.\n")
			fmt.Printf("Error:%s\n", err.Error())
		}
*/
	default:
		fmt.Printf("Unknown command %s\n", os.Args[1])
	}
	return nil
}
