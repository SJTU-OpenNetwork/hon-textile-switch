package core

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
    pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/pb"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/service"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/stream"
	"github.com/libp2p/go-libp2p-core/crypto"
)

// maxQueryWaitSeconds is used to limit a query request's max wait time
const maxQueryWaitSeconds = 30

// defaultQueryWaitSeconds is a query request's default wait time
const defaultQueryWaitSeconds = 5

// defaultQueryResultsLimit is a query request's default results limit
const defaultQueryResultsLimit = 5

// validation errors
const (
	errInvalidAddress = "invalid address"
	errUnauthorized   = "unauthorized"
	errForbidden      = "forbidden"
	errBadRequest     = "bad request"
)

// cafeServiceProtocol is the current protocol tag
const cafeServiceProtocol = protocol.ID("/textile/cafe/1.0.0")

// CafeService is a libp2p pinning and offline message service
type CafeService struct {
	service         *service.Service
	datastore       repo.Datastore
	inFlightQueries map[string]struct{}
    stream          *stream.StreamService
}

// NewCafeService returns a new threads service
func NewCafeService(
	node func() host.Host,
    ctx context.Context,
	datastore repo.Datastore,
    stream  *stream.StreamService,
	sk crypto.PrivKey) *CafeService {
	handler := &CafeService{
		datastore:       datastore,
		inFlightQueries: make(map[string]struct{}),
        stream:          stream,
	}
	handler.service = service.NewService(handler, node,sk)

    //subscribe topic
    ps, err := pubsub.NewGossipSub(ctx, handler.service.Node())
	if err != nil {
		fmt.Printf("error init pubsub")
	}
	sub, err = ps.Subscribe(string(cafeServiceProtocol))
	if err != nil {
		fmt.Printf("error subscribetopic")
	}
    go h.pubsubHandler(ctx, sub)
	return handler
}

func (h *CafeService) pubsubHandler(ctx context.Context, sub *pubsub.Subscription) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		mPeer := msg.From()
		if mPeer.Pretty() == srv.Node().Identity.Pretty() { // if the msg is from itself, just pass it
			continue
		}

		req := new(pb.Envelope)
		err := proto.Unmarshal(msg.Data(), req)
		if err != nil {
			log.Warningf("error unmarshaling pubsub message data from %s: %s", mPeer.Pretty(), err)
			continue
		}

		switch *req.Type {
		}
	}
}

// Protocol returns the handler protocol
func (h *CafeService) Protocol() protocol.ID {
	return cafeServiceProtocol
}

// Start begins online services
func (h *CafeService) Start() {
	h.service.Start()

}

// Handle is called by the underlying service handler method
func (h *CafeService) Handle(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	switch env.Message.Type {
	case pb.Message_CAFE_PUBSUB_QUERY:
		return h.handlePubSubQuery(env, pid)
//	case pb.Message_CAFE_PUBSUB_QUERY_RES:
//		return h.handlePubSubQueryResults(env, pid)
	default:
		return nil, nil
	}
}

// HandleStream is called by the underlying service handler method
func (h *CafeService) HandleStream(env *pb.Envelope, pid peer.ID) (chan *pb.Envelope, chan error, chan interface{}) {
	return make(chan *pb.Envelope), make(chan error), make(chan interface{})
}


func (h *CafeService) searchLocal(qtype pb.Query_Type, options *pb.QueryOptions, payload *any.Any, local bool, pid peer.ID) (*queryResultSet, error) {
	results := newQueryResultSet(options)

	switch qtype {
	case pb.Query_STREAM:
		fmt.Printf("cafe_service searchLocal: Search local stream\n")
        
        // have shadow node, return nothing
        // the shaodw node will connect other peers through other methods
        // [deprecated] Normal will return the peerId of shadow peer instead of return nothing.
        
		q := new(pb.StreamQuery)
		err := ptypes.UnmarshalAny(payload, q)
		if err != nil {
			return nil, err
		}

        meta := h.datastore.StreamMetas().Get(q.Id)
        if meta == nil{ 
            break
        }

		blocks := h.datastore.StreamBlocks().ListByStream(q.Id, int(q.Startindex),3)
        if q.Startindex != 0 && len(blocks) == 0 {
            break
        }

        peerId := h.service.Node().ID().Pretty()
		result := &pb.StreamQueryResultItem {
			Hopcnt: 0,
		}
		result.Pid = peerId
        
        value,_ := proto.Marshal(result)
		results.Add(&pb.QueryResult{
			Id:     q.Id,
			Local:  local,
			Value: &any.Any{
				TypeUrl: "/Stream",
                Value: value,
		    },
		})
	}

	return results, nil
}

// searchPubSub performs a network-wide search for the given query
//func (h *CafeService) searchPubSub(query *pb.Query, reply func(*pb.QueryResults) bool, cancelCh <-chan interface{}, fromCafe bool) error {
//	h.inFlightQueries[query.Id] = struct{}{}
//	defer func() {
//		delete(h.inFlightQueries, query.Id)
//	}()
//
//	err := h.publishQuery(&pb.PubSubQuery{
//		Id:           query.Id,
//		Type:         query.Type,
//		Payload:      query.Payload,
//		ResponseType: pb.PubSubQuery_P2P,
//		Exclude:      query.Options.Exclude,
//		Topic:        string(cafeServiceProtocol) + "/" + h.service.Node().Identity.Pretty(),
//		Timeout:      query.Options.Wait,
//	})
//	if err != nil {
//		return err
//	}
//
//	timer := time.NewTimer(time.Second * time.Duration(query.Options.Wait))
//	listener := h.queryResults.Listen()
//	doneCh := make(chan struct{})
//
//	done := func() {
//		listener.Close()
//		close(doneCh)
//	}
//
//	go func() {
//		<-timer.C
//		done()
//	}()
//
//	for {
//		select {
//		case <-cancelCh:
//			if timer.Stop() {
//				done()
//			}
//		case <-doneCh:
//			return nil
//		case value, ok := <-listener.Ch:
//			if !ok {
//				return nil
//			}
//			if r, ok := value.(*pb.PubSubQueryResults); ok && r.Id == query.Id && r.Results.Type == query.Type {
//				if reply(r.Results) {
//					if timer.Stop() {
//						done()
//					}
//				}
//			}
//		}
//	}
//}

// publishQuery publishes a search request to the network
//func (h *CafeService) publishQuery(req *pb.PubSubQuery) error {
//	env, err := h.service.NewEnvelope(pb.Message_CAFE_PUBSUB_QUERY, req, nil, false)
//	if err != nil {
//		return err
//	}
//	topic := string(cafeServiceProtocol)
//
//	log.Debugf("sending pubsub %s to %s", env.Message.Type.String(), topic)
//
//	payload, err := proto.Marshal(env)
//	if err != nil {
//		return err
//	}
//	//return h.pubsub.Publish(h.service.Node(), topic, payload)
//	
//    // TODO: send message to all neighbors
//    for () {
//    
//    }
//    return nil
//}

// handlePubSubQuery receives a query request over pubsub and responds with a direct message
func (h *CafeService) handlePubSubQuery(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	query := new(pb.PubSubQuery)
	err := ptypes.UnmarshalAny(env.Message.Payload, query)
	if err != nil {
		return nil, err
	}

	if _, ok := h.inFlightQueries[query.Id]; ok {
		return nil, nil
	}

	// return results, if any
	options := &pb.QueryOptions{
		Filter:  pb.QueryOptions_NO_FILTER,
		Exclude: query.Exclude,
	}
	results, err := h.searchLocal(query.Type, options, query.Payload, false, pid)
	if err != nil {
		return nil, err
	}
	if len(results.items) == 0 {
		return nil, nil
	}

	res := &pb.PubSubQueryResults{
		Id: query.Id,
		Results: &pb.QueryResults{
			Type:  query.Type,
			Items: results.List(),
		},
	}
	renv, err := h.service.NewEnvelope(pb.Message_CAFE_PUBSUB_QUERY_RES, res, nil, false)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), service.DefaultTimeout)
	defer cancel()

	err = h.service.SendMessage(ctx, pid.Pretty(), renv)
	if err != nil {
		fmt.Printf("error sending message response to %s: %s", pid, err)
	}
	return nil, nil
}

// handlePubSubQueryResults handles search results received from a pubsub query
//func (h *CafeService) handlePubSubQueryResults(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
//	res := new(pb.PubSubQueryResults)
//	err := ptypes.UnmarshalAny(env.Message.Payload, res)
//	if err != nil {
//		return nil, err
//	}
//
//	h.queryResults.Send(res)
//	return nil, nil
//}

// queryDefaults ensures the query is within the expected bounds
//func queryDefaults(query *pb.Query) *pb.Query {
//	if query.Options == nil {
//		query.Options = &pb.QueryOptions{
//			LocalOnly:  false,
//			RemoteOnly: false,
//			Limit:      defaultQueryResultsLimit,
//			Wait:       defaultQueryWaitSeconds,
//			Filter:     pb.QueryOptions_NO_FILTER,
//		}
//	}
//
//	if query.Options.Limit <= 0 {
//		query.Options.Limit = math.MaxInt32
//	}
//
//	if query.Options.Wait <= 0 {
//		query.Options.Wait = defaultQueryWaitSeconds
//	} else if query.Options.Wait > maxQueryWaitSeconds {
//		query.Options.Wait = maxQueryWaitSeconds
//	}
//
//	return query
//}
