package cmd

import (
	//"github.com/SJTU-OpenNetwork/hon-textile-switch/core"
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"strconv"
	"strings"
)
// api.go implements a simple tcp c/s framework to transport msg between textile-shadow routine.

const ApiPort = ":40101" //Avoid conflicting with original textile api port
const ApiLocal = "localhost"

// openRw open a bufio.ReadWriter to addr
func openRw(addr string) (*bufio.ReadWriter, error) {
	fmt.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

// Call send the cmd to api server and get the result from api server.
func SendCmd(cmd string) error {
	rw, err := openRw(ApiLocal+ApiPort)
	if err != nil {
		fmt.Printf("Api client can not connect to %s", ApiLocal+ApiPort)
		return err
	}
	n, err := rw.WriteString(cmd+"\n")
	if err != nil {
		return errors.Wrap(err, "Could not send the STRING request ("+strconv.Itoa(n)+" bytes written)\n")
	}
	err = rw.Flush()
	if err != nil {
		fmt.Printf("Error occur when flush rw buffer.\n")
		return err
	}

	response, err := rw.ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "Client: Failed to read the reply.\n")
	}
	response = strings.Trim(response, "\n")
	fmt.Printf("Get response: %s", response)
	return nil
}

