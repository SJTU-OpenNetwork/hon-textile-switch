package core

func (t *Textile) WhitelistAddItem(peerId string) error {
	return t.whiltelist.Add(peerId)
}

func (t *Textile) WhitelistRemoveItem(peerId string) error {
	return t.whiltelist.Remove(peerId)
}
