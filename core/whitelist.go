package core

import "fmt"

func (t *Textile) WhitelistAddItem(peerId string) error {
    err := t.whiteList.Add(peerId)
    if err != nil {
        return err
    }
    return t.shadow.WhitelistAddItem(peerId)
}

func (t *Textile) WhitelistRemoveItem(peerId string) error {
    err := t.whiteList.Remove(peerId)
    if err != nil {
        return err
    }
    t.shadow.WhitelistRemoveItem(peerId)
    return nil
}

func (t* Textile) PrintWhiteList() {
	if t.whiteList == nil {
		fmt.Printf("Textile have no whitelist\n")
		return
	}
	t.whiteList.PrintOut()
}
