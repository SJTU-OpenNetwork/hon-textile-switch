package core

import (
	"context"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo/config"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo/db"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/shadow"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/stream"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"os"
	"path"
	ma "github.com/multiformats/go-multiaddr"

	//"strings"
	//"sync"
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
	//stop              func() error
	node              p2phost.Host
	//started           bool
	datastore         repo.Datastore
	//online            chan struct{}
	//done              chan struct{}
	shadow            *shadow.ShadowService //add shadowservice 2020.04.05
	//lock              sync.Mutex
    stream            *stream.StreamService
	whiteList         repo.WhitrListStore
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
	dbDir := path.Join(repoPath, "datastore")
	os.Mkdir(dbDir, os.ModePerm)
	sqliteDb, err := db.Create(repoPath, "")
	if err != nil {
		return err
	}
	err = sqliteDb.InitTables("")
	//err = sqliteDb.Config().Init("")
	if err != nil {
		return err
	}

	//return applyTextileConfigOptions(conf)
	return nil
}

// NewTextile create a textile instance.
// Note that the repo should be initialized before.
func NewTextile(conf RunConfig) (*Textile, error) {
	repoPath := conf.RepoPath
	if !repo.IsInitialized(repoPath){
		fmt.Printf("Repo has not been initialized.\n")
		return nil, fmt.Errorf("Repo %s has not been initialized.")
	}

	// force open the repo and datastore
	// removeLocks(conf.RepoPath)

	node := &Textile{
		repoPath:          conf.RepoPath,
		ctx:			   context.Background(),
	}
	var err error
	node.config, err = config.Read(node.repoPath)
	if err != nil {
		return nil, err
	}

	sqliteDb, err := db.Create(node.repoPath, node.pinCode)
	if err != nil {
		return nil, err
	}

	node.datastore = sqliteDb

	// load whiteList
	whPath := path.Join(repoPath, "whitelist")
	whiteList, err := repo.NewWhiteListStore(whPath)
	if err != nil {
		return nil, err
	}
	node.whiteList = whiteList

	// Create host
	node.node, err = host.NewHost(node.ctx, node.repoPath, node.config)
	if err != nil {
		fmt.Printf("Error occur when create host\n")
		return nil, err
	}
	return node, nil
}

// Start creates an ipfs node and starts textile services
func (t *Textile) Start() error {

	// open db
    err := t.touchDatastore()
	if err != nil {
		return err
	}

	// create services
	t.stream = stream.NewStreamService(
		t.Host,
		t.datastore,
		t.repoPath,
		t.SubscribeStream,
		t.ctx,
		)
	t.shadow = shadow.NewShadowService(
        t.Host,
        t.datastore,
        t.shadowMsgRecv,
        t.Host().ID().Pretty())
/*
	go func() {
		defer func() {
			close(t.online)
			t.lock.Unlock()
		}()

		t.stream.Start()
        t.shadow.Start()
    }()
	t.started = true
*/
 	t.stream.Start()
 	t.shadow.Start()

 	// Outprint peer info
 	fmt.Printf("Host start with:\n")
 	fmt.Printf("PeerId: %s\n", t.node.ID().Pretty())
 	fmt.Printf("MultiAddr:\n")
 	for _, addr := range t.node.Addrs(){
 		fmt.Printf("%s\n", addr.String())
	}

    return nil
}

func (t *Textile) Host() p2phost.Host{
	return t.node
}


// touchDatastore ensures that we have a good db connection
func (t *Textile) touchDatastore() error {
	if err := t.datastore.Ping(); err != nil {
		//log.Debug("re-opening datastore...")
		fmt.Printf("Error occur when ping datastore.\nTry to re-open datastore\n")
		sqliteDB, err := db.Create(t.repoPath, t.pinCode)
		if err != nil {
			return err
		}
		t.datastore = sqliteDB
	}

	return nil
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

func (t *Textile)Connect(peerId string, addr string) error {
	mulAddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		fmt.Printf("Error occur when build multi address from %s\n", addr)
		return err
	}
	// Note that peerId is encoded string
	// We need to decode it to get the real peerId
	decodedId, err := peer.IDB58Decode(peerId)
	if err != nil {
		fmt.Printf("Error occur when decode peerId\n%s\n", err.Error())
		return err
	}
	pi := peer.AddrInfo{
		ID:    decodedId,
		Addrs: []ma.Multiaddr{mulAddr},
	}
	fmt.Printf("Try to connect with peer info:\nPeerId: %s\naddress: %s\n", pi.ID.Pretty(), addr)
	err = t.node.Connect(t.ctx, pi)
	if err != nil {
		fmt.Printf("Error occur when connect %s:%s\n", peerId, addr)
		return err
	} else {
		fmt.Printf("Connect %s successful\n", peerId)
	}
	fmt.Printf("Connect command end\n")
	return nil
}

