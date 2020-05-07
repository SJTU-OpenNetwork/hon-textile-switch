package cmd

import "errors"

// api_connect define how the api server handle "connect" command
func (s *Server)api_whitelist(params []string) error {
	if params[0] == "add" {
        return s.AddW
    }
    if params[0] == "remove" {
        return s.node.
    } 
	return errors.New("Invalid whitelist command")
}
