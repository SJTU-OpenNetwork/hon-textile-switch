package host

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo/config"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	secio "github.com/libp2p/go-libp2p-secio"
	tls "github.com/libp2p/go-libp2p-tls"
)

// option() build the libp2p options
// options include:
//		- private key, used to set the host identity
//		- protector, used to run host in private network
//		- security, used to set the security protocol used by connections
func option(repoPath string, cfg *config.Config) ([]libp2p.Option, error){
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

	// set security protocol
	opts = append(opts, libp2p.ChainOptions(libp2p.Security(secio.ID, secio.New), libp2p.Security(tls.ID, tls.New)))
	return opts, nil
}