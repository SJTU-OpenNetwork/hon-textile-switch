package core

import (
	"context"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/host"
	//"github.com/SJTU-OpenNetwork/hon-textile-switch/recorder"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo/config"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo/db"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/shadow"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/stream"
	"github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"os"
	"path"
	"strings"
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
	//record 			  *recorder.RecordService
	//lock              sync.Mutex
    stream            *stream.StreamService
    cafe              *CafeService
	whiteList         repo.WhiteListStore
	//pprofTask 		  *util.PprofTask
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
func InitRepo(conf InitConfig) error { // write the default conf files into usb dir
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

	pprofDir := path.Join(repoPath, "statistic")
	os.Mkdir(pprofDir, os.ModePerm)

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

	//node.pprofTask = util.NewPprofTask(path.Join(repoPath, "statistic"), path.Join(repoPath, "statistic"))
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
	sk, err := crypto.UnmarshalPrivateKey(t.config.PrivKey)
	if err != nil {
		fmt.Printf("Error occurs when unmarshal private key\n%s\n", err)
		return err
	}
	//streamCtx := context.WithValue(t.ctx, "pprof", t.pprofTask)
	t.stream = stream.NewStreamService(
		t.Host,
		t.datastore,
		t.repoPath,
		t.SubscribeStream,
		//streamCtx,
		t.ctx,
		t.getShadowIp,
		sk)

	var shadowIp1 string
	for _,x := range t.node.Addrs(){
		xx := strings.Split(x.String(),"/")
		if strings.Contains(xx[2],"192.168.3") {
			shadowIp1=xx[2]
			fmt.Println("get shadow IP:",shadowIp1)
		}
	}

	t.shadow = shadow.NewShadowService(
        t.Host,
        t.datastore,
        t.shadowMsgRecv,
        t.Host().ID().Pretty(),
        shadowIp1,
        sk,
        t.whiteList)
	//t.record = recorder.NewRecordService(t.Host, t.ctx, sk)
    t.cafe = NewCafeService(
        t.Host,
        t.ctx,
        t.datastore,
        t.stream,
        sk)
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
	//go t.pprofTask.StartMem(5*time.Second, true, context.Background())
	//go t.pprofTask.StartCpu(30*time.Second, true, context.Background())
 	t.stream.Start()
 	t.shadow.Start()
 	//t.record.Start()
 	t.cafe.Start()
    err = t.initMDNS()
    if err != nil {
        fmt.Printf("init mdns failed: %v\n", err)
        //return err
    }
	
 	// Outprint peer info
 	fmt.Printf("Host start with:\n")
 	fmt.Printf("PeerId: %s\n", t.node.ID().Pretty())
 	fmt.Printf("MultiAddr:\n")
 	for _, addr := range t.node.Addrs(){
 		fmt.Printf("%s\n", addr.String())
	}
	t.WritePeerInfo()
 	//t.tryExtractPublicKey()

 	//connect relay server
 	//"/ip4/202.120.38.100/tcp/4001/ipfs/QmZt8jsim548Y5UFN24GL9nX9x3eSS8QFMsbSRNMBAqKBb"
 	server1,_ := ma.NewMultiaddr("/ip4/202.120.38.100/tcp/4001/ipfs/QmZt8jsim548Y5UFN24GL9nX9x3eSS8QFMsbSRNMBAqKBb")
 	//server2,_ := ma.NewMultiaddr("/ip4/159.138.3.74/tcp/4001/ipfs/QmYovpcqB12c56AjGRaMUcwfoZg1DYinCFmEAzFHYvLb6R")
 	//server3,_ := ma.NewMultiaddr("/ip4/159.138.130.129/tcp/4001/ipfs/QmZX8WVgJ3cQCW3bNcodXhmK34rmNkvqk8Zg9u7f3JEFgN")
	addrInfos,_:=peer.AddrInfosFromP2pAddrs(server1)
	for _,addrInfo := range addrInfos {
		err=t.node.Connect(t.ctx,addrInfo)
		if err != nil {
			fmt.Println("connect relay server failed: ",addrInfo.ID)
		}else{
			fmt.Println("connect relay server successfully: ",addrInfo.ID)
		}
	}

    return nil
}

func (t *Textile) Host() p2phost.Host{
	return t.node
}
/*
type mdnsNotifee struct {
	h   p2phost.Host
	ctx context.Context
}

func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	err := m.h.Connect(m.ctx, pi)
	if err != nil {
		//fmt.Printf("error occur when connect mdns peer\nerr\n")
	} else {
		fmt.Printf("connect mdns peer %s\n", pi.ID.Pretty())
	}
}
*/
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
	fmt.Printf("Try to connect with peer info:\nPeerId: %s\naddress: %v\n", pi.ID.Pretty(), pi.Addrs)
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

func (t *Textile) PeerInfo() peer.AddrInfo {
	return peer.AddrInfo{
		ID: t.node.ID(),
		Addrs: t.node.Addrs(),
	}
}

func (t *Textile)tryExtractPublicKey() {
	fmt.Printf("Try extract public key from peer Id %s\n", t.node.ID().Pretty())
	_, err := t.node.ID().ExtractPublicKey()
	if err != nil {
		fmt.Printf("Error occur when extract public key from %s\n%s\n", t.node.ID().Pretty(), err)
	} else {
		fmt.Printf("Extract pubkey from peerId seccessfully\n")
	}

	pubk, err := crypto.UnmarshalPublicKey(t.config.Pubkey)
	if err != nil {
		fmt.Printf("Error occur when unmarshal public key\n", err)
	}

	fmt.Printf("Try extract peer Id frim public key\n")
	id, err := peer.IDFromPublicKey(pubk)
	if err != nil {
		fmt.Printf("Error occur when extract id from public key\n", err)
	} else {
		fmt.Printf("Extract peerId from public key seccessfully. Id: %s\n", id.Pretty())
	}

	privk, err := crypto.UnmarshalPrivateKey(t.config.PrivKey)
	if err != nil {
		fmt.Printf("Error occur when unmarshal private key\n", err)
	}

	fmt.Printf("Try extract peer Id frim private key \n")
	id2, err := peer.IDFromPrivateKey(privk)
	if err != nil {
		fmt.Printf("Error occur when extract id from private key\n", err)
	} else {
		fmt.Printf("Extract peerId from private key seccessfully. Id: %s\n", id2.Pretty())
	}

	//fmt.Printf("Try extract Id from public key %s\n", t.node.ID().Pretty())

}

func (t* Textile) WritePeerInfo() {
	infopath := path.Join(t.repoPath, "address")
	f, err := os.Create(infopath)
	if err != nil {
		fmt.Printf("Error occur when write info file\n%s\n", err)
		return
	}
	peerId := t.node.ID()
	for _, addr := range t.node.Addrs() {
		_, err = f.WriteString(addr.String()+"/ipfs/"+peerId.Pretty()+"\n")
		if err != nil {
			fmt.Printf("Error occur when write info file\n")
		}
	}

	defer f.Close()
}

func (t *Textile) getShadowIp() string {
	return t.shadow.ShadowIp()
}