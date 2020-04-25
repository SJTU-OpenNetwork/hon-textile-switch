package repo

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/util"
)

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
	flock *util.Flock
	dirPath string
	filePath string
	cache map[string]interface{}
}


// NewWhiteList Create a new whitelist.
// It does the following things:
//		- Create directory if not exists.
//		- Judge whether their already be a whitelist. Raise an error if so.
//		- Create Flock for this directory.
//		- Make cache to cache the whitelist
func NewWhiteList(dirPath string) (*WhiteList, error) {
	return nil, nil
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

	return nil
}

//func (w *WhiteList)
