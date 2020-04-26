package repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo/config"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/util"
)

var ErrRepoExists = fmt.Errorf("repo not empty, reinitializing would overwrite your account")
var ErrRepoDoesNotExist = fmt.Errorf("repo does not exist, initialization is required")
var ErrMigrationRequired = fmt.Errorf("repo needs migration")
var ErrRepoCorrupted = fmt.Errorf("repo is corrupted")

const Repover = "19"


// Judge whether this repo is already initialized.
//		- For now we just check whether there is files in this repo
func IsInitialized(repoPath string) bool {
	files, err := ioutil.ReadDir(repoPath)
	if err != nil {
		fmt.Printf("Error occur when read repo %s\n")
		fmt.Printf("%s\n", err.Error())
		return false
	}
	if len(files) > 1 {
		return true
	}
	return false
}

func Init(repoPath string) error {
	// Make directory
	if !util.DirectoryExists(repoPath) {
		err := os.Mkdir(repoPath, os.ModePerm)
		if err != nil {return err}
	} else {
		isInit := IsInitialized(repoPath)
		if isInit {
			return ErrRepoExists
		}
	}
	err := checkWriteable(repoPath)
	if err != nil {
		return err
	}


	fmt.Printf("initializing repo at %s\n", repoPath)

	// write default textile config
	tconf, err := config.Init()
	if err != nil {
		return err
	}
	err = config.Write(repoPath, tconf)
	if err != nil {
		return err
	}

	return nil
}

func checkWriteable(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		// dir exists, make sure we can write to it
		testfile := path.Join(dir, "test")
		fi, err := os.Create(testfile)
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("%s is not writeable by the current user", dir)
			}
			return fmt.Errorf("unexpected error while checking writeablility of repo root: %s", err)
		}
		fi.Close()
		return os.Remove(testfile)
	}

	if os.IsNotExist(err) {
		// dir doesnt exist, check that we can create it
		return os.MkdirAll(dir, 0775)
	}

	if os.IsPermission(err) {
		return fmt.Errorf("cannot write to %s, incorrect permissions", err)
	}

	return err
}

func initializeBotFolder(repoPath string) error {
	botFolder := filepath.Join(repoPath, "bots")
	_, err := os.Stat(botFolder)
	if os.IsNotExist(err) {
		// dir doesnt exist, check that we can create it
		return os.MkdirAll(botFolder, 0775)
	}
	return err
}


func checkPermissions(path string) (bool, error) {
	_, err := os.Open(path)
	if os.IsNotExist(err) {
		// repo does not exist yet - don't load plugins, but also don't fail
		return false, nil
	}
	if os.IsPermission(err) {
		// repo is not accessible. error out.
		return false, fmt.Errorf("error opening repository at %s: permission denied", path)
	}

	return true, nil
}
