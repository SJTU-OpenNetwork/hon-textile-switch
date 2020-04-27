package core

import "fmt"

func (t *Textile) SubscribeStream(streamId string) error {
	fmt.Printf("Shadow peer should not call subscribe!!\nUse request directly\n")
	return nil
}
