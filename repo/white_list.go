package repo

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/util"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

const writeRetry = 3	// Time to retry when write whitelist file failed.

type ErrFetchLockFail struct {
	dirPath string
}

func (e *ErrFetchLockFail) Error() string {
	return fmt.Sprintf("Fetch flock for %s failed.\n", e.dirPath)
}

// WhiteList contains infos about peers served by this shadow peer.
// It was stored by text file and cached by a string map.
// A Flock would be applied to the directory contains whitelist.
// How whitelist work:
//		- Load from text file to cache when start peer
//		- Read from cache when judge whether a peer is in white list
//		- Write to cache and further write back to text file when write white list
// NOTE:
//		- Flock has been removed from WhiteList.
//			That is because only unix based os supports flock.
//			And we need to run it on windows sometimes.
type WhiteList struct {
	//flock *util.Flock	// File lock for r/w whitelist file
	clock sync.Mutex	// cache lock for r/w cache
	dirPath string
	filePath string
	cache map[string]interface{}
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
	var res = &WhiteList{
		dirPath:	dirPath,
		filePath:	filePath,
		cache:		make(map[string]interface{}),
	}
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
	} else {
		// White list already exists
		fmt.Printf("White list %s exists\nLoad white list to cache\n", filePath)
		reloadList, err := res.readWhiteList()
		if err != nil {
			fmt.Printf("Error occur when reload white list from %s\n", filePath)
		} else {
			res.cache = reloadList
		}
	}
	fileLock, err := util.NewFlock(dirPath)
	if err != nil {
		fmt.Printf("Error occur when create file lock for whitelist\n")
		return nil, err
	}
	//res.flock = fileLock
	return res, nil
}

func (w *WhiteList) Check(peerId string) bool {
	w.clock.Lock()
	defer w.clock.Unlock()
	_, ok := w.cache[peerId]
	return ok
}

func (w *WhiteList) Add(peerId string) error {
	//ok := w.flock.Lock()
	w.clock.Lock()
	defer func() {
		w.clock.Unlock()
		//err := w.flock.Unlock()
		//if err != nil {
		//	fmt.Printf("Error occur when unlock whitelist dir:\n %s\n", err.Error())
		//}
	}()

	//if !ok {
	//	fmt.Printf("Fail to get flock of whitelist dir %s\n", w.dirPath)
	//	return &ErrFetchLockFail{dirPath:w.dirPath}
	//}

	// Do Add

	_, exists := w.cache[peerId]
	if exists {
		fmt.Printf("Add a already existing peer to whitelist: %s\n", peerId)
		return nil
	}
	w.cache[peerId] = struct{}{}
	err := w.writebackWhiteList(peerId)
	if err != nil {
		fmt.Printf("Error occurs when write peer %s back to whitelist text file.\n%s\n", err)
		// Do not raise this error out of this function
	}
	return nil
}

func (w *WhiteList) Remove(peerId string) error {
	//ok := w.flock.Lock()
	w.clock.Lock()
	defer func() {
		//err := w.flock.Unlock()
		//if err != nil {
		//	fmt.Printf("Error occur when unlock whitelist dir:\n %s\n", err.Error())
		//}
		w.clock.Unlock()
	}()

	//if !ok {
	//	fmt.Printf("Fail to get flock of whitelist dir %s\n", w.dirPath)
	//	return &ErrFetchLockFail{dirPath:w.dirPath}
	//}

	// Do Remove
	delete(w.cache, peerId)
	err := w.rewriteWhiteList()
	if err != nil {
		fmt.Printf("Error occurs when rewrite white list file\n%s\n", err)
		return nil
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

func (w *WhiteList) writebackWhiteList(peerId string) error{
	// Open file writeonly
	f, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_APPEND, 0644)
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error occur when close whitelist file %v\n", err)
		}
	}()
	if err != nil {
		fmt.Printf("Error occur when open whitelist file for write\n")
		return err
	} else {
		_, err = f.WriteString(peerId+"\n")
		if err != nil {
			fmt.Printf("Error occur when write peerId to whitelist file\n")
			return err
		}
	}
	return nil
}

func (w *WhiteList) rewriteWhiteList() error {
	f, err := os.Create(w.filePath)
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error occur when close whitelist file %v\n", err)
		}
	}()
	if err != nil {
		fmt.Printf("Error occur when flush whitelist file %s\n", err)
		return err
	}
	for peerId, _ := range w.cache {
		_, err = f.WriteString(peerId + "\n")
		if err != nil {
			fmt.Printf("Error occur when rewrite %s to whitelist file\n", peerId)
			return err
		}
	}
	return nil
}

func (w *WhiteList) PrintOut() {
	w.clock.Lock()
	defer w.clock.Unlock()
	fmt.Printf("Peers in whitelist:\n")
	if w.cache == nil {
		fmt.Printf("Nil white list cache\n")
		return
	}
	for peerId, _ := range w.cache{
		fmt.Printf("\t%s\n", peerId)
	}
}
