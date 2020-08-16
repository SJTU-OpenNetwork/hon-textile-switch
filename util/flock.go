// +build !windows

package util

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

type ErrNotDir struct{
	path string
}
type ErrDirNotLock struct{
	dirPath string
}

func (e *ErrNotDir) Error() string {
	return fmt.Sprintf("%s is not a directory.", e.path)
}

func (e *ErrDirNotLock) Error() string {
	return fmt.Sprintf("Unlock a unlocked directory %s.", e.dirPath)
}

//func (e *ErrFetchLockFail)

// flock implements file lock
// Use this to avoid file r/w conflict between programs.
// Each file lock will lock a directory.
// Note:
//		Only unix based os supports flock.
type Flock struct{
	dirPath string
	lockFile *os.File
	lockPath string
	locked bool
}

// NewFlock create a new file lock
func NewFlock(dirPath string) (*Flock, error) {
	// Check whether directory exists.
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}
	if !dirInfo.IsDir() {
		return nil, &ErrNotDir{path:dirPath}
	}

	// Check whether lock file exists.
	// Create one if not exist
	lockPath := path.Join(dirPath, "flock")
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		fmt.Printf("Create lock file: %s\n", lockPath)
		_, e := os.Create(lockPath)
		if e!= nil {
			fmt.Printf("Create lock file failed\n")
			return nil, e
		}
	}
	return &Flock{
		dirPath:  dirPath,
		lockPath: lockPath,
		locked: false,
	}, nil
}

// lock the file
// Return whether this routine get the file lock successfully.
func (l *Flock) Lock() bool {
	if l.locked {
		fmt.Printf("Directory already locked by this thread.\n")
		return false
	} else {
		var err error
		l.lockFile, err = os.Open(l.lockPath)
		if err != nil {
			fmt.Printf("Open lock file %s failed\n", l.lockPath)
			return false
		}
		err = syscall.Flock(int(l.lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err != nil {
			fmt.Printf("Cannot flock directory %s - %s\n", l.dirPath, err)
			return false
		}
		fmt.Printf("Lock directory %s\n", l.dirPath)
		l.locked = true
		return true
	}
}

func (l *Flock) Unlock() error {
	defer func() {
		if l.locked {
			l.lockFile.Close()
			l.locked = false
		}
	}()
	if !l.locked {
		fmt.Printf("Error: Directory %s is not locked.\n", l.dirPath)
		return &ErrDirNotLock{dirPath:l.dirPath}
	}
	return syscall.Flock(int(l.lockFile.Fd()), syscall.LOCK_UN) //
}

