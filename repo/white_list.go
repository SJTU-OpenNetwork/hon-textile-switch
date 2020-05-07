package repo

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/util"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const writeRetry = 3	// Time to retry when write whitelist file failed.

type ErrFetchLockFail struct {
	dirPath string
}

func (e *ErrFetchLockFail) Error() string {
	return fmt.Sprintf("Fetch flock for %s failed.\n", e.dirPath)
}

// WhiteList contains infos about peers served by this shadow peer.
// It was stored by text file.
// A Flock would be applied to the directory contains whitelist.
type WhiteList struct {
	flock *util.Flock	// File lock for r/w whitelist file
	//clock *sync.Mutex	// cache lock for r/w cache
	dirPath string
	filePath string
	//cache map[string]interface{}
}


// NewWhiteList Create a new whitelist.
func NewWhiteListStore(dirPath string) (*WhiteList, error) {
	err := util.Mkdirp(dirPath)
	if err != nil {
		fmt.Printf("Error occur when make dir for white list\n")
		return nil, err
	}
	filePath := path.Join(dirPath, "whitelist")
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		// Create empty whitelist file
		f, err := os.Create(filePath)
		defer f.Close()
		if err != nil {
			fmt.Printf("Error occur when create empty white list file %s\n", filePath)
			return nil, err
		}
	} else if err != nil {
		fmt.Printf("Error occur when read white list file info\n")
		return nil, err
	}
	fileLock, err := util.NewFlock(dirPath)
	if err != nil {
		fmt.Printf("Error occur when create file lock for whitelist\n")
		return nil, err
	}
	return &WhiteList{
		flock:    fileLock,
		dirPath:  dirPath,
		filePath: filePath,
	}, nil
}

func (w *WhiteList) Add(peerId string) error {
	ok := w.flock.Lock()
	defer func() {
		err := w.flock.Unlock()
		if err != nil {
			fmt.Printf("Error occur when unlock whitelist dir:\n %s\n", err.Error())
		}
	}()

	if !ok {
		fmt.Printf("Fail to get flock of whitelist dir %s\n", w.dirPath)
		return &ErrFetchLockFail{dirPath:w.dirPath}
	}

	// Do Add
    err := w.writeWhiteList(peerId)
    if err != nil {
        return err
    }

    // Send inform

	return nil
}

func (w *WhiteList) Remove(peerId string) error {
	ok := w.flock.Lock()
	defer func() {
		err := w.flock.Unlock()
		if err != nil {
			fmt.Printf("Error occur when unlock whitelist dir:\n %s\n", err.Error())
		}
	}()

	if !ok {
		fmt.Printf("Fail to get flock of whitelist dir %s\n", w.dirPath)
		return &ErrFetchLockFail{dirPath:w.dirPath}
	}

	// Do Remove
    list, err := w.readWhiteList()
    if err != nil {
        return err
    }
    _, ok = list[peerId]
    if !ok {
        return nil
    }

    delete(list, peerId)
    f, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error occur when open whitelist file for write\n")
		return err
	}
	defer f.Close()

    for peerId, _ := range (list) {
		_, err = f.WriteString(peerId + "\n")
		if err != nil {
			fmt.Printf("Error occur when write peerId to whitelist file\n")
			return err
		}
    }
	return nil
}

// DO NOT ask for any lock within this func.
func (w *WhiteList) readWhiteList() (map[string]interface{}, error){
	data, err := ioutil.ReadFile(w.filePath)
	if err != nil {
		return nil, err
	}
	data_str := string(data[:])
	strs := strings.Split(data_str, "\n")
	res := make(map[string]interface{})
	for _, peerId := range strs {
		if peerId != "" {
			// Split may contains "" if data_str is ""
			res[peerId] = struct{}{}
		}
	}
	return res, nil
}

func (w *WhiteList) writeWhiteList(peerId string) error{
	// Open file writeonly
	f, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error occur when open whitelist file for write\n")
		return err
	} else {
		_, err = f.WriteString(peerId + "\n")
		if err != nil {
			fmt.Printf("Error occur when write peerId to whitelist file\n")
			return err
		}
	}
	defer f.Close()
	return nil
}
