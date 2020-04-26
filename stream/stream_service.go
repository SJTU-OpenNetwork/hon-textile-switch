// Service for sending/receving stream related data - add by Jerry 2020/02/25

package stream

import (
	"encoding/json"
	"fmt"

//    "bytes"
	"context"

	"github.com/golang/protobuf/ptypes"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/pb"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/service"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/util"
	"github.com/libp2p/go-libp2p-core/host"
)


// streamServiceProtocol is the current protocol tag
const streamServiceProtocol = protocol.ID("/textile/stream/1.0.0")
var ErrRedundantReq = fmt.Errorf("Request is redundant")
var ErrUnknowkStream = fmt.Errorf("Unknown stream")

type StreamService struct {
	service          *service.Service
	datastore        repo.Datastore
	online           bool
    subscribe        func(string) error 
	repoPath         string
	
    // for workers
    activeWorkers *workerStore
    
    // for providers
    providers *providerStore
	lostIndex chan *lostReport

	// Context for main routine
	ctx context.Context
}

// NewStreamService returns a new stream service
func NewStreamService(
	node func() host.Host,
	datastore repo.Datastore,
    repoPath  string,
    subscribe func(string) error,
	ctx context.Context,
) *StreamService {
	handler := &StreamService{
		datastore:        datastore,
        repoPath:         repoPath,
        subscribe:        subscribe,
		ctx:			  ctx,
		activeWorkers: newWorkerStore(),
        providers: newProviderStore(),
	}
	handler.service = service.NewService(handler, node)
	return handler
}

// Protocol returns the handler protocol
func (h *StreamService) Protocol() protocol.ID {
	return streamServiceProtocol
}

// Start begins online services
func (h *StreamService) Start() {
    h.online = true
	h.service.Start()
	h.service.Node().Network().Notify((*StreamNotifee)(h))
}


// Handle is called by the underlying service handler method
func (h *StreamService) Handle(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	fmt.Printf("core/stream_service.go Handler: New message receive from %s.\n", pid.Pretty())
	switch env.Message.Type {
	case pb.Message_STREAM_BLOCK_LIST:
		return h.handleStreamBlockList(env, pid)
	case pb.Message_STREAM_REQUEST:
		return h.handleStreamRequest(env, pid)
	case pb.Message_STREAM_UNSUBSCRIBE:
		return h.handleUnsubscribe(env, pid)
    default:
    	fmt.Printf("core/stream_service.go Handler: Unknown message type")
        return nil, nil
    }
}

func (h *StreamService) UnsubscribeStream(sid string) error{
    fmt.Printf("StreamService: Try to unsubscribe stream %s\n", sid)
    pid := h.providers.RemoveStream(sid)
    if pid != peer.ID("") {
        h.SendUnsubscribeRequest(pid.Pretty(), sid)
    }
    return nil
}

// ======================== FOR MESSAGE RECV/SEND ==================================
// handleStreamBlock receives a STREAM_BLOCK_LIST message
func (h *StreamService) handleStreamBlockList(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	//fmt.Printf("StreamService: New stream blk list receive from %s\n", pid.Pretty())
    streams := make(map[string]int)
    blks := new(pb.StreamBlockContentList)
    err := ptypes.UnmarshalAny(env.Message.Payload, blks)
    if err != nil {
        return nil, err
    }
    for _, blk := range blks.Blocks {
        size := 0
        if len(blk.Data) != 0 {
            //TODO: save blk in file system
            err := util.Store(h.repoPath, blk.Description, blk.Data)
            if err != nil {
                return nil, err
            }
        }
        model := &pb.StreamBlock {
            Id: blk.Description,
            Streamid: blk.StreamID,
            Index: blk.Index,
            Size: int32(size),
            IsRoot: blk.IsRoot,
            Description: string(blk.Description),
        }
        //fmt.Printf("StreamService: Received stream %s; index %d; cid %s\n", blk.StreamID, blk.Index, cid.String())
        //log.Debugf("[%s] Block %s, Stream %s, Index %d, From %s, Size %d", TAG_BLOCKRECEIVE, cid_str, blk.StreamID, blk.Index, pid.Pretty(), size)
        err = h.datastore.StreamBlocks().Add(model)
        if err != nil {
            return nil, err
        }
        //fmt.Printf("It is successfully stored in our database!\n")

        if blk.IsRoot {
            // we found a file !
            fmt.Print("It is a root node of a merkle-DAG!\n")
            err = h.handleRootBlk(pid, model)
            if err != nil {
                fmt.Printf("Handle root file failed\n")
                return nil, err
            }
        }
        streams[blk.StreamID] = 1
    }
    for id := range streams {
	    h.activeWorkers.newFileAdd(id)
    }
    return nil, nil
}

func (h *StreamService) handleRootBlk(pid peer.ID, blk *pb.StreamBlock) error {
    if blk.Id == "" {
        meta := h.datastore.StreamMetas().Get(blk.Streamid)
	    if meta == nil || meta.Nblocks > 0{
		    return nil
	    }
        err := h.datastore.StreamMetas().UpdateNblocks(blk.Streamid, blk.Index)
        if err != nil {
            //log.Error(err)
            return err
        }
    }
    return nil
}


// HandleStream is called by the underlying service handler method
func (h *StreamService) HandleStream(env *pb.Envelope, pid peer.ID) (chan *pb.Envelope, chan error, chan interface{}) {
	return make(chan *pb.Envelope), make(chan error), make(chan interface{})
}

// SendMessage sends a message to a peer
func (h *StreamService) SendMessage(ctx context.Context, peerId string, env *pb.Envelope) error {
	return h.service.SendMessage(ctx, peerId, env)
}

// HandleRequest
func (h *StreamService) handleStreamRequest(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	fmt.Printf("core/stream_service.go handleStreamRequest from %s\n", pid.Pretty())
	req := new(pb.StreamRequest)
	err := ptypes.UnmarshalAny(env.Message.Payload, req)
	if err != nil {
		return nil, err
	}

    // TODO: calculate capacity according to video rate
    if h.Workload() < 5 {
        err = h.responseRequest(pid, req)
        if err != nil {
            return nil, err
        }
	    return h.service.NewEnvelope(pb.Message_STREAM_REQUEST_HANDLE, &pb.StreamRequestHandle{
		    Value:1,
	    },nil, true)
    } else {
	    return h.service.NewEnvelope(pb.Message_STREAM_REQUEST_HANDLE, &pb.StreamRequestHandle{
		    Value:0,
	    },nil, true)
    }
}

func (h *StreamService) SendStreamRequest(peerId string, config *pb.StreamRequest) (*pb.Envelope, error) {
	//fmt.Printf("core/stream_service.go SendStreamRequest to %s\n", peerId)
	env, err := h.service.NewEnvelope(pb.Message_STREAM_REQUEST, config, nil, false)
	if err != nil {
		return nil,err
	}
	//log.Debugf("[%s] Stream %s, To %s", TAG_STREAMREQUEST, config.Id, peerId)
	return h.service.SendRequest(peerId, env)
}

func (h *StreamService) SendUnsubscribeRequest(peerId string, sid string) (*pb.Envelope, error) {
	//fmt.Printf("core/stream_service.go SendStreamRequest to %s\n", peerId)
	env, err := h.service.NewEnvelope(pb.Message_STREAM_UNSUBSCRIBE, &pb.StreamUnsubscribe{
        Id: sid,
    }, nil, false)
	if err != nil {
		return nil,err
	}
	//log.Debugf("[%s] Stream %s, To %s", TAG_STREAMREQUEST, sid, peerId)
	return h.service.SendRequest(peerId, env)
}

func (h *StreamService) handleUnsubscribe(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	fmt.Printf("core/stream_service.go handleUnsubscribe from %s\n", pid.Pretty())
	req := new(pb.StreamUnsubscribe)
	err := ptypes.UnmarshalAny(env.Message.Payload, req)
	if err != nil {
		return nil, err
	}
    
    //TODO: stop sending data to pid
    h.activeWorkers.endWorker(req.Id, pid.Pretty())

	return h.service.NewEnvelope(pb.Message_STREAM_UNSUBSCRIBE_RES, &pb.StreamUnsubscribeAck{
        Id:req.Id,
    }, nil, true)
}

// RequestAccepted is called when a stream request is accepted by some peer.
func (h *StreamService) RequestAccepted(peerId string, config *pb.StreamRequest) {
	acceptedSubstream := newProvidedSubstream(config.Id, config.StreamMap, 1, config.StartIndex, peerId, h.handleBlockLost)
	provider := h.providers.getOrCreate(peerId)
	// TODO: De-duplicated
	provider.add(acceptedSubstream)
}

// Handle lost block
// It would be called by providedSubstream.
func (h *StreamService) handleBlockLost(report *lostReport){
	fmt.Println("BlockLost")
	infos :=  report.Loggable()
	js, err := json.MarshalIndent(infos, "", "  ")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	} else {
		fmt.Printf("%s\n", string(js))
	}
}

// Call it when you decide to send blocks to requestor.
// Use "Response" to distinguish with "Handle".
func (h *StreamService) responseRequest(pid peer.ID, req *pb.StreamRequest) error {
	//log.Debugf("[%s] Stream %s, From %s", TAG_STREAMRESPONSE, req.Id, pid.Pretty())

	// Raise an error if obtain the same request (Same requestor with same streamid and substream requested by the same peer)
	if h.activeWorkers.isRedundant(pid, req) {
		return ErrRedundantReq
	}
	worker, err := h.createWorker(pid, req)
	if err != nil {
		return err
	}
	// add worker to activeworkers
	err = h.activeWorkers.add(worker)
	if err != nil {
		return err
	}
	// start worker
	return worker.start()
}

//SendStreamBlocks send a list of block to a peer
func (h *StreamService) SendStreamBlocks(peerId peer.ID, blks []*pb.StreamBlock) error{
	//fmt.Printf("StreamService: Send %d stream blks to %s\n", len(blks), peerId.Pretty())

	// Marshal blocks to pb
    blist := new(pb.StreamBlockContentList)
    for _, blk:= range blks {
        var data []byte
        if blk.Id != "" {
            //TODO: get block from database and file system
            var err error
            data, err = util.Get(h.repoPath, blk.Id)
		    if err != nil {
			    return err
		    }
        }
        content := &pb.StreamBlockContent{
            StreamID: blk.Streamid,
            Index: blk.Index,
            Data: data,
            IsRoot: blk.IsRoot,
            Description: blk.Id,
        }
        //log.Debugf("[%s] Block %s, Stream %s, Index %d, To %s, Size %d, description: %s", TAG_BLOCKSEND, blk.Id, blk.Streamid, blk.Index, peerId.Pretty(), blk.Size, blk.Description)
        blist.Blocks = append(blist.Blocks, content)
    }
	env, err := h.service.NewEnvelope(pb.Message_STREAM_BLOCK_LIST, blist, nil, false)
	if err != nil {
        //log.Error(err)
		return err
	}
	// Send envelope use StreamService.service.SendMessage
    err = h.service.SendMessage(nil, peerId.Pretty(), env)
    if err != nil {
        //log.Error(err)
    }
	return nil
}

// FetchBlocks fetches a list of blocks of a specific stream from database
func (h *StreamService) FetchBlocks(streamId string, startIndex uint64, maxNum int) ([]*pb.StreamBlock, error){
    // find blocks of the stream with id = streamId
    blks := h.datastore.StreamBlocks().ListByStream(streamId, int(startIndex),maxNum)
	if blks == nil{
		return nil,fmt.Errorf("stream blocks fetch failed")
	}
    return blks, nil
}

func (h *StreamService) AddPotential(pid string, config *pb.StreamRequest, hopcnt int) {
}


// =============== FOR WORKDERS ==================
func (h *StreamService) createWorker(pid peer.ID, req *pb.StreamRequest) (*streamWorker, error) {
	//fmt.Printf("stream/streamManager createWorker")
    stream := h.datastore.StreamMetas().Get(req.Id)
	if stream == nil {
		return nil, ErrUnknowkStream
	}
	return newStreamWorker(stream, pid, req, h.FetchBlocks, h.SendStreamBlocks), nil
}

func (h *StreamService) Workload() int {
    // return the required bitrate for all workers
    // for now, just return the number of workers
    return h.activeWorkers.Workload()
}

func (h *StreamService) WorkerStat() {
	//log.Debugf("StreamManager.WorkerStat()\n")
	h.activeWorkers.PrintOut()
}


// ============== FOR PEER MANAGEMENT ===================
func (h *StreamService)PeerDisconnected(pid peer.ID) {
	// Stop all the workers
	//log.Debugf("Peer %s disconnected", pid)
	h.activeWorkers.endPeer(pid.Pretty())
    provider := h.providers.remove(pid.Pretty())
    if provider != nil{
        // re-subscribe streams
        for _,stream := range(provider.subStreams) {
            err := h.subscribe(stream.streamId)
            if err != nil{
                //log.Error(err)
            }
        }
    }
}

func (h* StreamService) GetProvidedHopcnt(config *pb.StreamRequest) (int, bool){
	return 0, false
}

func (h* StreamService) GetProvider(sid string) peer.ID{
	return h.providers.GetProvider(sid)
}
// ===================== OTHERS =========================
func (h *StreamService) Loggable() map[string]interface{}{
	return h.activeWorkers.Loggable()
}

