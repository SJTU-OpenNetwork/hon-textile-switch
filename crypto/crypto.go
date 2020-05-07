package crypto

import (
	"fmt"
	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
)


func Verify(pk libp2pc.PubKey, data []byte, sig []byte) error {
	good, err := pk.Verify(data, sig)
	if err != nil || !good {
		return fmt.Errorf("bad signature")
	}
	return nil
}


