package host

import (
	"context"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo/config"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
    peer "github.com/libp2p/go-libp2p-core/peer"
	//repo "github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
)

// Create host for shadow peer

func NewHost(ctx context.Context, repoPath string, cfg *config.Config) (p2phost.Host, error){
	var privKey crypto.PrivKey
	var err error

	// load privKey
	opts := make([]libp2p.Option,0)
	if cfg.PrivKey!= nil {
		privKey, err = crypto.UnmarshalPrivateKey(cfg.PrivKey)
		if err != nil {
			fmt.Printf("Error occurs when ummarshal private key from config.\n%s\nCreate host with random key.\n", err.Error())
		} else {
			opts = append(opts, libp2p.Identity(privKey))
		}
	}

	// load swarm key
	protec := repo.GetProtector(repoPath)
	if protec != nil {
		opts = append(opts, libp2p.PrivateNetwork(protec))
	}
	return libp2p.New(ctx, opts...)
}

func Peers(host p2phost.Host) ([]peer.ID) {
	conns := host.Network().Conns()
    var result []peer.ID
	for _, c := range conns {
		pid := c.RemotePeer()
		//addr := c.RemoteMultiaddr()

		result = append(result, pid)
	}
	return result
}

func SwarmConnected(host p2phost.Host, pid peer.ID) (bool) {
    plist := Peers(host)
    for _, peer := range plist {
        if peer == pid {
            return true
        }
    }
    return false
}
