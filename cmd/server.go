package cmd

import (
	"bufio"
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
	return &Server{}
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
			fmt.Println("Read command successfully.")
			return
		case err != nil:
			fmt.Println("Read command fail.\nError:%s\n", err.Error())
			return
		}
		cmd = strings.Trim(cmd, "\n")
		fmt.Printf("Get command: %s\n", cmd)

		// handle specific command
		cmd_formatted := buildCommand(cmd)
		switch cmd_formatted.cmd {
		case "connect":
			err = s.api_connect(cmd_formatted.args)
		default:
			err = fmt.Errorf("Unknown cmd: %s", cmd_formatted.cmd)
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

// Format the string command
func buildCommand(cmd string) *command {
	params :=  strings.Split(cmd, " ")
	res := &command{
		cmd: "",
		args: make([]string,1),
	}
	for i, p := range params {
		if i==0 {
			res.cmd = p
		} else {
			res.args = append(res.args, p)
		}
	}
	return res
}

type command struct {
	cmd string
	args []string
}
