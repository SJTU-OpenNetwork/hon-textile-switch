package host

import (
	"context"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo/config"
	"github.com/libp2p/go-libp2p"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	//repo "github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
)

// Create host for shadow peer

func NewHost(ctx context.Context, repoPath string, cfg *config.Config) (p2phost.Host, error){
	opts, err := option(repoPath, cfg)
	if err != nil {
		fmt.Printf("Error occurs when build libp2p options\n%s\n", err.Error())
		return nil, err
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
