// Service for sending/receving messages to the shadow node - add by Jerry 2020.04.05

package shadow

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/pb"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	"sync"
	"time"

	"github.com/SJTU-OpenNetwork/hon-textile-switch/service"
	"github.com/libp2p/go-libp2p-core/crypto"
)


// streamServiceProtocol is the current protocol tag
const shadowServiceProtocol = protocol.ID("/textile/shadow/1.0.0")
const informRetry = 5
const informInternal = time.Second * 3
var ErrWrongRole = fmt.Errorf("Wrong role.")	//shadow function called at normal peer or vice versa.

type ShadowService struct {
	service          *service.Service
	datastore        repo.Datastore
	online           bool
	msgRecv          func(*pb.Envelope, peer.ID) error
	address			 string  // public key. textile.account.Address()
    whiteList        repo.WhiteListStore
	lock             sync.Mutex
}

type shadowInfo struct {
	peerId		peer.ID
	multiAddress ma.Multiaddr
}

func NewShadowService(
	node func() host.Host,
	datastore repo.Datastore,
	msgRecv func(*pb.Envelope, peer.ID) error,
	address string,
	key crypto.PrivKey,
	whiteList repo.WhiteListStore,
) *ShadowService {
	handler := &ShadowService{
		datastore:        datastore,
		msgRecv:          msgRecv,
		address:		  address,
		whiteList:		  whiteList,
	}
	handler.service = service.NewService(handler, node, key)
	return handler
}

// Protocol returns the handler protocol
func (h *ShadowService) Protocol() protocol.ID {
	return shadowServiceProtocol
}

// Start begins online services
func (h *ShadowService) Start() {
	h.online = true
	h.service.Start()
	h.service.Node().Network().Notify((*ShadowNotifee)(h))
}

// Handle is called by the underlying service handler method
func (h *ShadowService) Handle(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	fmt.Printf("core/shadow_service.go Handler: New message receive from %s.\n", pid.Pretty())
	switch env.Message.Type {
	case pb.Message_SHADOW_STREAM_META:
		return h.handleStreamMeta(env, pid)
	case pb.Message_SHADOW_INFORM:
		return h.handleInform(env, pid)
    default:
        return nil, nil
    }
}

func (h *ShadowService) handleInform(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	fmt.Printf("Inform message from %s\n", pid.Pretty())
	return nil, nil
}


func (h *ShadowService) PeerDisconnected(pid peer.ID) error{
    h.lock.Lock()
    defer h.lock.Unlock()
    //fmt.Println("")
	return nil
}


// TODO:
// 		Avoid to call it multi time for the same peer!!
func (h *ShadowService) PeerConnected(pid peer.ID, multiaddr ma.Multiaddr) {
	// check whitelist
	exist := h.whiteList.Check(pid.Pretty())
	if exist {
		fmt.Printf("A peer within whitelist connected. Try to send inform to %s\n", pid.Pretty())
		err := h.inform(pid); if err != nil {
			fmt.Printf("Error occurs when inform peer %s\n%s\n", pid.Pretty(), err.Error())
		}
	}
    //return nil
}

// TODO: inform pid about my information (e.g., public key), could use ``contact'' directly
func (h *ShadowService) inform(pid peer.ID) error {
	fmt.Printf("Shadow: Send inform to %s\n", pid.Pretty())

	inform := &pb.ShadowInform{}
	inform.PublicKey = h.address
	env, err := h.service.NewEnvelope(pb.Message_SHADOW_INFORM, inform, nil, true); if err != nil {return err}
	informTime := 0;
	for {
		informTime += 1
		err = h.service.SendMessage(nil, pid.Pretty(), env)
		if err != nil {
			fmt.Printf("Send inform failed: %v\n", err)
			if informTime < informRetry {
				fmt.Printf("Retry after %v \n", informInternal)
				time.Sleep(informInternal)
				continue
			} else {
				return err
			}
			//return err
		} else {
			return nil
		}
	}
    return nil
}

func (h *ShadowService) handleStreamMeta(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
    err := h.msgRecv(env, pid)
    if err != nil {
    	fmt.Printf("Error occur when receive stream meta %v\n", err)
	}
	return nil, err
}

func (h *ShadowService) ShadowStat() *pb.ShadowStat {
	res := &pb.ShadowStat{}
	res.Role = "shadow"
	//ulist:= make([]string, 0, len(h.users))

	return res
}

// Retuen the stat of ShadowService
func (h *ShadowService) Loggable() map[string]interface{} {
	res := make(map[string] interface{})
	res["role"] = "shadow"
	return res
}

// HandleStream is called by the underlying service handler method
func (h *ShadowService) HandleStream(env *pb.Envelope, pid peer.ID) (chan *pb.Envelope, chan error, chan interface{}) {
	return make(chan *pb.Envelope), make(chan error), make(chan interface{})
}
