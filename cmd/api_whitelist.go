package cmd

import (
	"errors"
	"fmt"
)

// api_connect define how the api server handle "connect" command
func (s *Server)api_whitelist(params []string) error {
	if params == nil {
		return errors.New("Nil params when call api_whitelist\n")
	}
	if len(params) < 1 || (len(params)<2 && params[0] != "list") {
		return errors.New("Not enough parameter")
	}
	//return s.node.Connect(params[0], params[1])
	switch params[0] {
	case "add":
		return s.node.WhitelistAddItem(params[1])
	case "remove":
		return s.node.WhitelistRemoveItem(params[1])
	case "list":
		s.node.PrintWhiteList()
		return nil
	default:
		return errors.New(fmt.Sprintf("Unknown sub command %s", params[0]))
	}
}
