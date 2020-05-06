package repo

import (
	"bytes"
	"fmt"
	pnet "github.com/libp2p/go-libp2p-pnet"
	ipnet "github.com/libp2p/go-libp2p-core/pnet"
	"io/ioutil"
	"os"
	"path"
)

// pnet.go handle private net.
// provide methods to get libpp2p protector

const swarmKeyFile = "swarm.key"
const defaultSwarmKey = "/key/swarm/psk/1.0.0/\nbase16/\n7894ae706f3b54675785afd43a5a554463744e89594ab6e274fb817ccd9a58d4\n"

// swarmKey read the swarmKeyFile
func swarmKey(repoPath string) ([]byte, error) {
	//repoPath := filepath.Clean(r.path)
	spath := path.Join(repoPath, swarmKeyFile)

	f, err := os.Open(spath)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return nil, err
	}
	defer f.Close()
	fmt.Printf("Read swarm key from %s\n", spath)
	return ioutil.ReadAll(f)
}

// GetProtector fetch swarmkey in swarmKeyFile and return the libp2p protector.
func GetProtector(repoPath string) ipnet.Protector {
	swarmkey, err := swarmKey(repoPath)
	if err != nil {
		fmt.Printf("Error occurs when read swarm key.\n%s\nHost will run on public network\n", err.Error())
		return nil
	} else if swarmkey == nil {
		fmt.Printf("No swarm key\nHost will run on public network\n")
		return nil
	}

	protec, err := pnet.NewProtector(bytes.NewReader(swarmkey))
	if err != nil {
		fmt.Printf("Error occurs when build protector from swarmkey.\nHost will run on public network\n")
		return nil
	}
	return protec
}