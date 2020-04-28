package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// Server handle request from other routine
type Server struct {
	listener net.Listener
	//taskQueue chan net.Conn
}

func NewServer() *Server {
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
			fmt.Printf("Get command: %s\n", cmd)
			rw.WriteString("ok")
			rw.Flush()
			return
		case err != nil:
			fmt.Println("Read command fail.\nError:%s\n", err.Error())
			return
		}

		//cmd = strings.Trim(cmd, "\n ")

	}
}
