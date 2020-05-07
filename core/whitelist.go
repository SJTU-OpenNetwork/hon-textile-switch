package core

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/pb"
)

func (t *Textile) SubscribeStream(streamId string) error {
	fmt.Printf("Shadow peer should not call subscribe!!\nUse request directly\n")
	return nil
}

// RequestStream request a stream from a provider by sending stream request.
func (t* Textile) RequestStream(pid string, config *pb.StreamRequest) (*pb.Envelope, error){
	return t.stream.SendStreamRequest(pid, config)
}