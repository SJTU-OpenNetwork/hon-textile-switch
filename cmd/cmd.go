package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/core"
	"os"
	"strings"
)

// Package cmd define the command line instructions

// command is the basic struct pass between api server and client
// It can be marshal to a json string and unmarshal at server.
type command struct {
	Type string `json:"type"`
	Args []string `json:"args"`
}

// buildCommand build a command from command line parameters
func buildCommand() *command {
	fmt.Print(os.Args)
	fmt.Printf("length of os.Args: %d\n", len(os.Args))
	res := &command{
		Type: "",
		Args: make([]string,1),
	}
	var tmpCount = 0

	for _, p := range os.Args {
		// Remove the empty args (I have no idea why they appear)
		pTrim := strings.Trim(p, " ")
		if pTrim != "" {
			if tmpCount==1 {
				res.Type = pTrim
			} else if tmpCount > 1 {
				res.Args = append(res.Args, pTrim)
			}
			tmpCount += 1
		}

	}
	return res
}

func marshalCommand(tmpcmd *command) (string, error) {
	fmt.Printf("Try to marshal %s command\n", tmpcmd.Type)
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

	switch tmpcmd.Type {
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
