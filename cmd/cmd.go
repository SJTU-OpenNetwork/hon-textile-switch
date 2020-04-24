package cmd

import (
	"github.com/SJTU-OpenNetwork/hon-textile-switch/core"
	"os"
)

// Package cmd define the command line instructions

type cmdsMap map[string]func(string) error	// Store which func for each cmd

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
