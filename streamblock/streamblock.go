// This package is used to handle the streamblock
package streamblock

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/pb"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/utils"
	"io/ioutil"
	"os"
)

var ErrPathDoesNotExist = fmt.Errorf("path does not exist")
var ErrDataDoesNotExist = fmt.Errorf("blockdata does not exist")

// Store stores the received streamblock to path
func Store(block pb.StreamBlock,path string) error{
	//create dir if not exist
	if !utils.DirectoryExists(path) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		err :=  ioutil.WriteFile(block.Id,block.data,os.ModePerm)
		if err!=nil {
			return err
		}
	}
	return nil
}

// Get gets the data from path
func Get(path string, filename string) ([]byte,error) {
	if !utils.DirectoryExists(path) {
		return nil, ErrPathDoesNotExist
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file == filename {
			data,err := ioutil.ReadFile(filename)
			if err != nil{
				return data, nil
			}
		}
	}
	return nil, ErrDataDoesNotExist
}