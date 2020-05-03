package cmd

import "errors"

// api_connect define how the api server handle "connect" command
func (s *Server)api_connect(params []string) error {
	if len(params) <= 1 {
		return errors.New("Not enough parameter")
	}
	return s.node.Connect(params[0], params[1])
}
