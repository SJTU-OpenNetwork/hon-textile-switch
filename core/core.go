package core

import (
	"fmt"
    "context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/service"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo/db"
	"strings"
	"sync"
)


// InitConfig is used to setup a textile node
type InitConfig struct {
	RepoPath        string
	//SwarmPorts      string
}


// RunConfig is used to define run options for a textile node
type RunConfig struct {
	RepoPath          string
}

type Variables struct {
	SwarmAddress    string
	FailedAddresses []string
    lock            sync.Mutex
}


// common errors
var ErrAccountRequired = fmt.Errorf("account required")
var ErrStarted = fmt.Errorf("node is started")
var ErrStopped = fmt.Errorf("node is stopped")
var ErrOffline = fmt.Errorf("node is offline")
var ErrMissingRepoConfig = fmt.Errorf("you must specify InitConfig.RepoPath or InitConfig.BaseRepoPath and InitConfig.Account")

// Repo returns the actual location of the configured repo
func (conf InitConfig) Repo() (string, error) {
	if len(conf.RepoPath) > 0 {
		return conf.RepoPath, nil
	} else {
		return "", ErrMissingRepoConfig
	}
}

// InitRepo initializes a new node repo.
// It does the following things:
//		- Create repo directory.
//		- Create datastore and save it to directory.
func InitRepo(conf InitConfig) error {
	repoPath, err := conf.Repo()
	if err != nil {
		return err
	}

	// init repo
	err = repo.Init(repoPath)

	if err != nil {
		return err
	}

	//sqliteDb, err := db.Create(repoPath, "")
	_, err = db.Create(repoPath, "")
	if err != nil {
		return err
	}
	//err = sqliteDb.Config().Init("")
	if err != nil {
		return err
	}


	//return applyTextileConfigOptions(conf)
	return nil
}

func Start(ctx context.Context, phost host.Host) error {
    service.NewShadowService(ctx, phost)
    return nil
}


type loggingWaitGroup struct {
	n  string
	wg sync.WaitGroup
}

func (lwg *loggingWaitGroup) Add(delta int, src string) {
	lwg.wg.Add(delta)
}

func (lwg *loggingWaitGroup) Done(src string) {
	lwg.wg.Done()
}

func (lwg *loggingWaitGroup) Wait(src string) {
	lwg.wg.Wait()
}

func getTarget(output string) string {
	str := strings.Split(output, " ")[1];
	return str;
}

func getStatus(output string) bool {
	str := strings.Split(output, " ")[2];
	if str=="success"{
		return true
	}else {
		return false
	}
}

func getId(output string) string {
	list := strings.Split(output, "/")
	return list[len(list)-1]
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}


func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}



