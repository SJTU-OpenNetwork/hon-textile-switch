package core

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

const discoveryConnTimeout = time.Second * 30

type discoveryNotifee struct {
	//PeerChan chan peer.AddrInfo
	node *Textile
}

//interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	//n.PeerChan <- pi
	ctx, cancel := context.WithTimeout(n.node.ctx, discoveryConnTimeout)
	defer cancel()
	fmt.Printf("Try to connect with peer info:\nPeerId: %s\naddress: %v\n", pi.ID.Pretty(), pi.Addrs)
	if err := n.node.Host().Connect(ctx, pi); err != nil {
		fmt.Printf("Connect mdns peer %s failed\n", pi.ID.Pretty())
		fmt.Printf("Error %v\n", err)
	} else {
		fmt.Printf("Connect mdns peer %s succeed\n", pi.ID.Pretty())
	}
}

//Initialize the MDNS service
func (t *Textile)initMDNS() error {
	// An hour might be a long long period in practical applications. But this is fine for us
	ser, err := discovery.NewMdnsService(t.ctx, t.Host(), time.Minute, discovery.ServiceTag)
	if err != nil {
		return err
	}

	//register with service so that we get notified about peer discovery
	n := &discoveryNotifee{
		//PeerChan: make(chan peer.AddrInfo, 10),
		node: t,
	}
	//n.PeerChan = make(chan peer.AddrInfo)

	ser.RegisterNotifee(n)
	/*
    go func(){
        //TODO: connect peer in the channel
        //TODO: return when according to ctx

    }()
	 */
	return nil
}
