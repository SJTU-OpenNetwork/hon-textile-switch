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

type Textile struct {
	repoPath          string
	pinCode           string
	config            *config.Config
	ctx               context.Context
	stop              func() error
	node              *host.Host
	started           bool
	datastore         repo.Datastore
	online            chan struct{}
	done              chan struct{}
	shadow            *shadow.ShadowService //add shadowservice 2020.04.05
	lock              sync.Mutex
    stream            *stream.StreamService
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

func NewTextile(conf RunConfig) (*Textile, error) {
    if !fsrepo.IsInitialized(conf.RepoPath) {
		return nil, repo.ErrRepoDoesNotExist
	}

	// check if repo needs a major migration
	err := repo.Stat(conf.RepoPath)
	if err != nil {
		return nil, err
	}

	// force open the repo and datastore
	removeLocks(conf.RepoPath)

	node := &Textile{
		repoPath:          conf.RepoPath,
		pinCode:           conf.PinCode,
	}

	node.config, err = config.Read(node.repoPath)
	if err != nil {
		return nil, err
	}

	sqliteDb, err := db.Create(node.repoPath, node.pinCode)
	if err != nil {
		return nil, err
	}
	node.datastore = sqliteDb

	return node, nil
}

// Start creates an ipfs node and starts textile services
func (t *Textile) Start() error {
	t.lock.Lock()
	if t.started {
		t.lock.Unlock()
		return ErrStarted
	}

	t.online = make(chan struct{})
	t.done = make(chan struct{})

	// open db
	err = t.touchDatastore()
	if err != nil {
		return err
	}

	// create services
	t.stream = service.NewStreamService(
		t.account,
		t.Ipfs,
		t.datastore,
		t.sendNotification,
        t.SubscribeStream,
		context.Background())//Share the same ctx with textile. That is because we do not need to manually cancel it.
	t.shadow = shadow.NewShadowService(
        t.account,
        t.Ipfs,
        t.datastore,
        t.shadowMsgRecv,
        t.config.IsShadow,
        t.account.Address())
    t.cafe = NewCafeService(
		t.account,
		t.Ipfs,
		t.datastore,
		t.cafeInbox,
        t.stream,
        t.shadow)

	go func() {
		defer func() {
			close(t.online)
			t.lock.Unlock()
		}()

		t.stream.Start()
        t.shadow.Start()
    }
	t.started = true
    return nil
}


// touchDatastore ensures that we have a good db connection
func (t *Textile) touchDatastore() error {
	if err := t.datastore.Ping(); err != nil {
		log.Debug("re-opening datastore...")

		sqliteDB, err := db.Create(t.repoPath, t.pinCode)
		if err != nil {
			return err
		}
		t.datastore = sqliteDB
	}

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



