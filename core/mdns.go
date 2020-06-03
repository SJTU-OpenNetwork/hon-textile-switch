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
	if err := n.node.Host().Connect(ctx, pi); err != nil {
		fmt.Printf("Error %v\n", err)
	}
}

//Initialize the MDNS service
func (t *Textile)initMDNS() error {
	// An hour might be a long long period in practical applications. But this is fine for us
	ser, err := discovery.NewMdnsService(t.ctx, t.Host(), time.Minute, "")
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
