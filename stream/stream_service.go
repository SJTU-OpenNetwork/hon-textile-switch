// Service for sending/receving stream related data - add by Jerry 2019/02/25

package stream

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"io/ioutil"
	"net"
	"time"

	//    "bytes"
	"context"

	"github.com/SJTU-OpenNetwork/hon-textile-switch/pb"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/recorder"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/service"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/util"
	"github.com/golang/protobuf/ptypes"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
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

	shadowIp func() string

	//pprofTask *util.PprofTask
}

// NewStreamService returns a new stream service
func NewStreamService(
	node func() host.Host,
	datastore repo.Datastore,
    repoPath  string,
    subscribe func(string) error,
	ctx context.Context,
	shadowIp func() string,
	key crypto.PrivKey,
) *StreamService {
	handler := &StreamService{
		datastore:        datastore,
        repoPath:         repoPath,
        subscribe:        subscribe,
		ctx:			  ctx,
		activeWorkers: newWorkerStore(),
        providers: newProviderStore(),
        shadowIp: shadowIp,
        //pprofTask: ctx.Value("pprof").(*util.PprofTask),
	}
	handler.service = service.NewService(handler, node, key)
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
	//h.service.Node().Network().Notify((*StreamNotifee)(h))

	//TODO: start TCP server socket, listen and accept stream block push from peer
	go func(){
		listener,_:=net.Listen("tcp",":40121")
		fmt.Println("start listen TCP block push at:",h.shadowIp())
		for{
			conn,_:=listener.Accept()
			go h.handleTCPBlockList(conn)
		}
	}()
}

func (h *StreamService) handleTCPBlockList(conn2 net.Conn) error{
	//tmp := make([]byte,1024)
	//var buff bytes.Buffer
	//for {
	//	c,_:=conn2.Read(tmp)
	//	if string(tmp[0:c]) == "tcpend" {
	//		fmt.Println("get an tcp end")
	//		break
	//	}
	//	buff.Write(tmp[0:c])
	//}
	//buf:=buff.Bytes()

	buf,_:=ioutil.ReadAll(conn2)

	streams := make(map[string]int)
	blks := new(pb.StreamBlockContentList)
	proto.Unmarshal(buf,blks)

	fmt.Println("get an connection with",len(blks.Blocks),"blocks")

	for _,blk := range blks.Blocks{
		size := 0
		var cid string
		if len(blk.Data) != 0 {
			m := make(map[string]string)
			err := json.Unmarshal(blk.Description, &m)
			if err != nil{
				return err
			}
			cid = m["CID"]

			err = util.WriteFileByPath(h.repoPath+"/blocks/"+cid, blk.Data)
			if err != nil {
				fmt.Printf("error occur when store file")
				return err
			}
		} else {
			cid = ""
		}
		model := &pb.StreamBlock {
			Id: cid,  //get id from description
			Streamid: blk.StreamID,
			Index: blk.Index,
			Size: int32(size),
			IsRoot: blk.IsRoot,
			Description: string(blk.Description),
		}
		//fmt.Printf("[BLKRECV] Block %s, Stream %s, Index %d, Size %d\n", model.Id, blk.StreamID, blk.Index, size)
		h.datastore.StreamBlocks().Add(model)

		if blk.IsRoot {
			//fmt.Print("It is a root node of a merkle-DAG!\n")
			h.handleRootBlk("", model)
		}
		streams[blk.StreamID] = 1
	}
	for id := range streams {
		h.activeWorkers.newFileAdd(id)
	}


	return nil
}


// Handle is called by the underlying service handler method
func (h *StreamService) Handle(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	fmt.Printf("core/stream_service.go Handler: New message receive from %s.\n", pid.Pretty())
	switch env.Message.Type {
	case pb.Message_STREAM_BLOCK_LIST:
		//return h.handleStreamBlockList(env, pid)
		return nil, nil
	case pb.Message_STREAM_REQUEST:
		return h.handleStreamRequest(env, pid)
	case pb.Message_STREAM_UNSUBSCRIBE:
		return h.handleUnsubscribe(env, pid)
    default:
    	fmt.Printf("core/stream_service.go Handler: Unknown message type")
        return nil, nil
    }
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
        fmt.Printf("size: %d\n", len(blk.Data))
        var cid string
        if len(blk.Data) != 0 {
            m := make(map[string]string)
            err := json.Unmarshal(blk.Description, &m)
            if err != nil{
                return nil, err
            }
            cid = m["CID"]

            err = util.WriteFileByPath(h.repoPath+"/blocks/"+cid, blk.Data)
            if err != nil {
                fmt.Printf("error occur when store file")
                return nil, err
            }
        } else {
            cid = ""
        }
        model := &pb.StreamBlock {
            Id: cid,  //get id from description
            Streamid: blk.StreamID,
            Index: blk.Index,
            Size: int32(size),
            IsRoot: blk.IsRoot,
            Description: string(blk.Description),
        }
        fmt.Printf("[BLKRECV] Block %s, Stream %s, Index %d, From %s, Size %d", model.Id, blk.StreamID, blk.Index, pid.Pretty(), size)
        err = h.datastore.StreamBlocks().Add(model)
        if err != nil {
            return nil, err
        }

        if blk.IsRoot {
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
	//h.pprofTask.NoticeMem()
	//h.pprofTask.NoticeCpu()
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

        // Send notification back
		record2 := &pb.Notification{
			Block: blk.Streamid,
			Date:  ptypes.TimestampNow(),
			//Actor:                t.node().Identity.Pretty(),	// Whether this is id of this peer ?
			Subject: recorder.Event_DoneSHADOW,
			Target:  pid.Pretty(),
			Read:    false, // Do not send to notification channel directly
		}
		recorder.RecordCh <- record2
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
        fmt.Print(err)
		return nil, err
	}

    err = h.responseRequest(pid, req)
    if err != nil {
        return nil, err
    }
    return h.service.NewEnvelope(pb.Message_STREAM_REQUEST_HANDLE, &pb.StreamRequestHandle{
	    Value:1,
    },nil, true)

    // TODO: calculate capacity according to video rate
//    if h.Workload() < 5 {
//        err = h.responseRequest(pid, req)
//        if err != nil {
//            return nil, err
//        }
//	    return h.service.NewEnvelope(pb.Message_STREAM_REQUEST_HANDLE, &pb.StreamRequestHandle{
//		    Value:1,
//	    },nil, true)
//    } else {
//	    return h.service.NewEnvelope(pb.Message_STREAM_REQUEST_HANDLE, &pb.StreamRequestHandle{
//		    Value:0,
//	    },nil, true)
//    }
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
	tStart:=time.Now()
    for _, blk:= range blks {
        var data []byte
        if blk.Id != "" {
            //TODO: get block from database and file system
            var err error

            data, err = util.ReadFileByPath(h.repoPath + "/blocks/", blk.Id)
		    if err != nil {
			    return err
		    }
        }
        content := &pb.StreamBlockContent{
            StreamID: blk.Streamid,
            Index: blk.Index,
            Data: data,
            IsRoot: blk.IsRoot,
            Description: []byte(blk.Description),
        }
        fmt.Println("Send Index ", blk.Index)
        //log.Debugf("[%s] Block %s, Stream %s, Index %d, To %s, Size %d, description: %s", TAG_BLOCKSEND, blk.Id, blk.Streamid, blk.Index, peerId.Pretty(), blk.Size, blk.Description)
        blist.Blocks = append(blist.Blocks, content)
    }
	fmt.Println("\n---after read: ",time.Since(tStart))

	env, err := h.service.NewEnvelope(pb.Message_STREAM_BLOCK_LIST, blist, nil, false)
	if err != nil {
        //log.Error(err)
		return err
	}
	// Send envelope use StreamService.service.SendMessage
	tStart=time.Now()
    err = h.service.SendMessage(nil, peerId.Pretty(), env)

    // TODO: remove this after test
    time.Sleep(100*time.Millisecond)

    if err != nil {
        //log.Error(err)
    	fmt.Println("Error, send message failed: ", err)
    }
    fmt.Println("\n---after send: ",time.Since(tStart))
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
	return newStreamWorker(stream, pid, req, h.FetchBlocks, h.SendStreamBlocks, h.ctx), nil
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
//        for _,stream := range(provider.subStreams) {
//            err := h.subscribe(stream.streamId)
//            if err != nil{
//                //log.Error(err)
//            }
//        }
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

