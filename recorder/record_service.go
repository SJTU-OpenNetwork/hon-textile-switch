package recorder

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/pb"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/service"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/opentracing/opentracing-go/log"
)

// recrod service is used to collect statistics.
// For now it is used to analyze the time used for some distributing tasks.
// record service use pb.Notification directly to transform info.
// Fields in notification is used as following:
//		- id string: A unique id used to distinguish different notification.
//			Note that each notification should have its own id.
//		- block string: A unique id used to distinguish different event.
//			It would the cid of file when collecting file distributing information.
//		- actor string: Self peer id.
//			Record may be collected by some other peers. Use this to distinguish the peer generate this record.
//			This may be empty when received from RecordCh. In that case, fill it with self id before sending it.
//		- subject string: Event type of record.
//			A special "final" event is used to represent the final event for an unique id.
//		- date timestamp: Time when this event happens.
//		- target string:  A remote peer id.
//			Record would be sent to that peer if it needs to be collected by other peer.
//			It would be empty if this notification should not be sent to other peer.
//		- read bool:  Whether to send this record to notification channel.
// The workflow of record service is:
//		- Useful information will be sent to RecordCh when some events happen.
//		- Check "read == true ?". Send notification to notification channel if it is.
//		- Check target. Send notification to collector.
//		- Collector receives messages containing notifications from other peers.
//			Send notifications to notification channel if final event received.

// streamServiceProtocol is the current protocol tag
const recordServiceProtocol = protocol.ID("/textile/record/1.0.0")

var RecordCh = make(chan *pb.Notification, 10)
var Online = false

type RecordService struct {
	service          *service.Service
	online           bool
	peerId 			 string // self peer id.
	// Context for main routine
	ctx context.Context
}

// NewStreamService returns a new stream service
func NewRecordService(
	node func() host.Host,
	ctx context.Context,
	key crypto.PrivKey,
) *RecordService {
	handler := &RecordService{
		//sendNotification: sendNotification,
		ctx:			  ctx,
		//peerId : node().ID().Pretty(),
		//peerId:node().Identity.Pretty(),	//Should call it after the ipfs node start
	}
	handler.service = service.NewService(handler, node, key)
	return handler
}

// Protocol returns the handler protocol
func (h *RecordService) Protocol() protocol.ID {
	return recordServiceProtocol
}

// Start begins online services
func (h *RecordService) Start() {
	h.online = true
	h.service.Start()
	h.peerId = h.service.Node().ID().Pretty()
	//h.service.Node().Identity.ShortString()
	go h.ListenRecordCh()
	Online = true
}

// Handle is called by the underlying service handler method
func (h *RecordService) Handle(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	fmt.Printf("Error: Receive recorder from %s. recorder of shadow peer should now receive anything from other peer.\n", pid.Pretty())
	return nil, nil
}


// HandleStream is called by the underlying service handler method
func (h *RecordService) HandleStream(env *pb.Envelope, pid peer.ID) (chan *pb.Envelope, chan error, chan interface{}) {
	return make(chan *pb.Envelope), make(chan error), make(chan interface{})
}

// SendMessage sends a message to a peer.
func (h *RecordService) sendMessage(ctx context.Context, peerId string, env *pb.Envelope) error {
	return h.service.SendMessage(ctx, peerId, env)
}

func (h *RecordService) sendNotificationToPeer(notification *pb.Notification, peerId string) error {
	fmt.Printf("Send record notification to %s\n", peerId)
	env, err := h.service.NewEnvelope(pb.Message_RECORD_NOTIFICATION, notification, nil, false)
	if err != nil {
		fmt.Printf("Error occurs when send notification %v\n", err)
		return err
	}
	return h.service.SendMessage(nil, peerId, env)
}

func (h *RecordService) ListenRecordCh() {
	for {
		select {
		case n := <- RecordCh:
			err := h.handleRecordChannel(n)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (h *RecordService) handleRecordChannel(notification *pb.Notification) error {
	// fill self peer
	notification.Actor = h.peerId

	// fill notification type
	notification.Type = pb.Notification_RECORD_REPORT


	// check whether need to send it to other peer
	if notification.Target != "" {
		err := h.sendNotificationToPeer(notification, notification.Target)
		if err != nil {
			fmt.Println("Error occur when send notification from record channel to notification channel.")
			fmt.Printf("%v\n", err)
			return err
		}
	} else {
		// Do nothing ??
	}
	return nil
}
