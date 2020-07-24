package stream

import (
	"context"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/pb"
	"github.com/libp2p/go-libp2p-core/peer"
	"time"
)

const maxBlockFetchNum = 1

// StreamWorker is used do blocksending task.
// Each streamrequest will create a independent worker.
// Note:
//		If a peer request two different substream of the same stream with two seperate request,
//		StreamManager will create two independent worker.
//		However, if a peer request two substream with one single request, there will be only one worker created.
type streamWorker struct {
	stream *pb.StreamMeta	// Contains stream info
	req *pb.StreamRequest 	// Contains core information such as substream and index
	pid peer.ID				// Contains information about destination
	currentIndex uint64		// The index of block sending now
    end bool
	workSignal chan interface{}
	cancelSignal chan interface{}
	blockFetcher func(streamId string, startIndex uint64, maxNum int) ([] *pb.StreamBlock, error)
	blockSender func (destination peer.ID, streamBlk [] *pb.StreamBlock) error
	//stopSignal chan interface{}
	ctx context.Context
	//pprofTask *util.PprofTask
}

func newStreamWorker(
	stream *pb.StreamMeta,
	pid peer.ID,
	req *pb.StreamRequest,
	blockFetcher func(streamId string, startIndex uint64, maxNum int) ([] *pb.StreamBlock, error),
	blockSender func (destination peer.ID, streamBlk [] *pb.StreamBlock) error,
	ctx context.Context) *streamWorker{

		//var pprofTask *util.PprofTask
		//pprofTask, ok := ctx.Value("pprof").(*util.PprofTask)
		//if !ok {
		//	pprofTask = nil
		//	fmt.Println("Error: no pprof in worker context!")
		//}

		return &streamWorker{
			stream: stream,
			req: req,
			pid: pid,
			currentIndex: req.StartIndex,
			workSignal: make(chan interface{}, 1),
			cancelSignal: make(chan interface{}, 1),
			blockFetcher: blockFetcher,
			blockSender: blockSender,
            end: false,
            ctx: ctx,
            //pprofTask: pprofTask,
		}
}

func (sw *streamWorker) notice() {
	// workSignal has buffer size 1.
	// notice() would not block if the worker has already been noticed.
	// Note that default has lower priority than case in golang.
	select{
		case sw.workSignal <- struct{}{}:
		default:
	}
}

// cancel worker
// Note:
// 		cancel have nothing to do with workerstore.
func (sw *streamWorker) cancel(){
	select {
		case sw.cancelSignal <- struct{}{}:
		default:
	}
}

func (sw *streamWorker) start() error {
	//log.Debugf("[%s] Stream %s, To %s", TAG_WORKERSTART, sw.stream.Id, sw.pid.Pretty())
	//fmt.Printf("stream/streamWorker.go start(): Worker for stream %s to %s start\n", sw.stream.Id, sw.pid.Pretty())
	// Start the block sending routine
	sw.currentIndex = sw.req.StartIndex
	sw.notice() //notice once at begining
	go func(){
		//defer fmt.Printf("stream/streamWorker.go start(): worker for stream %s to %s end\n", sw.stream.Id, sw.pid.Pretty())
		//defer log.Debugf("[%s] Stream %s, To %s", TAG_WORKEREND, sw.stream.Id, sw.pid.Pretty())
		tmpTime := time.Now()
		for {
			var retry int
			//tmpTime := time.Now()
			select {
				case <-sw.workSignal:
					// Do sending
					// Block if there is no signal
					//if sw.pprofTask != nil {
						//sw.pprofTask.NoticeCpu()
						//sw.pprofTask.NoticeMem()
					//}
					//tmpTime := time.Now()
					blks, _ := sw.blockFetcher(sw.req.Id, sw.currentIndex, maxBlockFetchNum)
					//duration := time.Since(tmpTime)
					//fmt.Println("----- Fetch ", len(blks), " blocks with ", time.Since(tmpTime) , "ms")
					fmt.Println("=== 2 circle and fetch: ", time.Since(tmpTime))
					tmpTime=time.Now()
					if blks != nil && len(blks) > 0 {
						fmt.Printf("stream/streamWorker.go start(): send %d blks for stream %s to %s start\n", len(blks), sw.stream.Id, sw.pid.Pretty())
						//fblks := sw.filterBlocks(blks)

						err := sw.blockSender(sw.pid, blks)
						if err != nil {
							retry += 1
							fmt.Println("Error when sending blocks: ", err)
							//log.Errorf("%s\nError occur when sending blocks.", err.Error())
							if retry <= 3 {
                            	time.Sleep(time.Duration(100)*time.Millisecond) //something wrong, maybe the connection breaks, if that happens, the worker will be canceled
								sw.notice() // Resend block if the connection is still there
							}
                            break
						}
						retry = 0
						sw.currentIndex = sw.currentIndex + uint64(len(blks))
						// Notice the worker again if there maybe more blocks can be fetched.
						if len(blks) >= maxBlockFetchNum {
							sw.notice()
						}
                        if blks[len(blks)-1].Id == "" {
                            sw.cancel()
                        }
					}
					fmt.Println("worker send: ",time.Since(tmpTime))
					tmpTime=time.Now()

					//duration = time.Since(tmpTime)
					//fmt.Println("===== Fetch and send ", len(blks), " blocks with ", duration.Milliseconds(), "ms")

				case <- sw.cancelSignal:
					// Note that break will break select only.
                    sw.end = true
					return
			}
		}
	}()

	return nil
}


// filter blocks to find the blocks belongs to certain substream.
func (sw *streamWorker) filterBlocks(blks []*pb.StreamBlock) []*pb.StreamBlock {
	streamMap := sw.req.StreamMap
	res := make([]*pb.StreamBlock, 0)
	for _, blk := range blks {
		subIndex := blk.Index % uint64(sw.stream.Nsubstreams)
		subMap := uint64(1) << subIndex
		if subMap & streamMap != 0 {
			res = append(res, blk)
		}
	}
	return res
}

func (sw *streamWorker) isSame(pid peer.ID, req *pb.StreamRequest) bool {
	return false
}

// Convert basic info of worker to loggable map
//func (sw *streamWorker) Loggab<F11>le() map[string]interface{} {
//
//}


