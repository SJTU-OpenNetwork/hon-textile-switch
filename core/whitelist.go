package core

func (t *Textile) WhitelistAddItem(peerId string) error {
	return t.whiteList.Add(peerId)
}

func (t *Textile) WhitelistRemoveItem(peerId string) error {
	return t.whiteList.Remove(peerId)
}

func (t* Textile) PrintWhiteList() {
	t.whiteList.PrintOut()
}