package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/core"
	"os"
)

// Package cmd define the command line instructions

// command is the basic struct pass between api server and client
// It can be marshal to a json string and unmarshal at server.
type command struct {
	cmd string
	args []string
}

// buildCommand build a command from command line parameters
func buildCommand() *command {
	res := &command{
		cmd: "",
		args: make([]string,1),
	}
	for i, p := range os.Args {
		if i==1 {
			res.cmd = p
		} else if i>1{
			res.args = append(res.args, p)
		}
	}
	return res
}

func marshalCommand(tmpcmd *command) (string, error) {
	fmt.Printf("Try to marshal %s command\n", tmpcmd.cmd)
	js, err := json.Marshal(tmpcmd)
	if err != nil {
		return "", err
	}
	return string(js), nil
}


// Run() start the software
//		- parse command parameters
//		- call the functions in core package
func Run() error {
	numArgs := len(os.Args)
	if numArgs <= 1 {
		fmt.Printf("Not enough parameter.\n")
		return nil
	}
	tmpcmd := buildCommand()
	strcmd, err := marshalCommand(tmpcmd)
	if err != nil {
		fmt.Printf("Error occurs when marshal json command:\n%s\n", err.Error())
		return err
	}

	// call functions in core

	switch tmpcmd.cmd {
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
		err = NewServer(textile).Listen()
		return err
        // hang forever

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

		err = SendCmd(strcmd)
		if err != nil {
			fmt.Printf("Error occurs when send json command to api server\n%s\n", err.Error())
			return err
		}

	default:
		fmt.Printf("Unknown command %s\n", os.Args[1])
	}

	return nil
}
