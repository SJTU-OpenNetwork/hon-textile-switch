package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/core"
	"io"
	"net"
	"strings"
)

// TODO:
//		Use seperate goroutine to handle each command

// Server handle request from other routine
type Server struct {
	listener net.Listener
	node *core.Textile
	//taskQueue chan net.Conn
}

func NewServer(node *core.Textile) *Server {
	return &Server{
		node: node,
	}
}

func (s *Server) Listen() error {
	var err error
	s.listener, err = net.Listen("tcp", ApiPort)
	if err != nil {
		fmt.Printf("Error occurs when try to listen port %s for api.", ApiPort)
		return err
	}
	fmt.Printf("Listen %s for api service.\n", s.listener.Addr().String())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Fail to listen api port.\n")
			continue
		}
		go s.handleMessage(conn)
	}
}

func (s *Server) handleMessage(conn net.Conn) {
	defer conn.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(conn),
		bufio.NewWriter(conn))
	for{
		cmd, err := rw.ReadString('\n')
		switch  {
		case err == io.EOF:
			fmt.Println("Command reader close\n")
			return
		case err != nil:
			fmt.Printf("Read command fail.\nError:%v\n", err)
			return
		}
		cmd = strings.Trim(cmd, "\n")
		fmt.Printf("Get command: %s\n", cmd)

		// handle specific command
		//cmd_formatted := buildCommand(cmd)
		cmd_formatted := &command{}
		err = json.Unmarshal([]byte(cmd), cmd_formatted)
		if err != nil {
			fmt.Printf("Error occurs when unmarshal json command\n%s\n", err.Error())
			return
		}
		switch cmd_formatted.Type {
		case "connect":
			err = s.api_connect(cmd_formatted.Args)
		case "whitelist":
			err = s.api_whitelist(cmd_formatted.Args)
		default:
			err = fmt.Errorf("Unknown cmd: %s", cmd_formatted.Type)
		}
		if err != nil {
			fmt.Printf("Error occurs when execute command.\nError:%s\n", err.Error())
			rw.WriteString("error\n")
		}else {
			rw.WriteString("ok\n")
		}

		rw.Flush()
	}
}

