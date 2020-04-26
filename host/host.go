package host

import (
	"context"
	"github.com/libp2p/go-libp2p"
	p2phost "github.com/libp2p/go-libp2p-core/host"
    peer "github.com/libp2p/go-libp2p-core/peer"
)

// Create host for shadow peer

func NewHost(ctx context.Context) (p2phost.Host, error){
	return libp2p.New(ctx)
}

func NewHost2(ctx context.Context)(p2phost.Host, error){
	return nil, nil
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
