package core

import "fmt"

func (t *Textile) WhitelistAddItem(peerId string) error {
	return t.whiteList.Add(peerId)
}

func (t *Textile) WhitelistRemoveItem(peerId string) error {
	return t.whiteList.Remove(peerId)
}

func (t* Textile) PrintWhiteList() {
	if t.whiteList == nil {
		fmt.Printf("Textile have no whitelist\n")
		return
	}
	t.whiteList.PrintOut()
}