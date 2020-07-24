package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	ggio "github.com/gogo/protobuf/io"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	inet "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-msgio"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/pb"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/util"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/crypto"
	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
)

// service represents a libp2p service
type Service struct {
	Node    func() host.Host
	handler Handler
	strmap map[peer.ID]*messageSender
	smlk   sync.Mutex
	privateKey libp2pcrypto.PrivKey
}

// DefaultTimeout is the context timeout for sending / requesting messages
const DefaultTimeout = time.Second * 10 //change from 30 to 10 2019.11.26

// PeerStatus is the possible results from pinging another peer
type PeerStatus string

const (
	PeerOnline  PeerStatus = "online"
	PeerOffline PeerStatus = "offline"
)

// Handler is used to handle messages for a specific protocol
type Handler interface {
	Protocol() protocol.ID
	Start()
	Handle(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error)
	HandleStream(env *pb.Envelope, pid peer.ID) (chan *pb.Envelope, chan error, chan interface{})
}

// NewService returns a service for the given config
func NewService(handler Handler, node func() host.Host, privKey libp2pcrypto.PrivKey) *Service {
	return &Service{
		Node:    node,
		handler: handler,
		strmap:  make(map[peer.ID]*messageSender),
		privateKey:privKey,
	}
}

// Start sets the peer host stream handler
func (srv *Service) Start() {
	srv.Node().SetStreamHandler(srv.handler.Protocol(), srv.handleNewStream)
	//go srv.listen("")
	//go srv.listen(srv.Node().Identity.Pretty())
}

// SendRequest sends out a request
func (srv *Service) SendRequest(p string, pmes *pb.Envelope) (*pb.Envelope, error) {

	pid, err := peer.IDB58Decode(p)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	ms, err := srv.messageSenderForPeer(ctx, pid)
	if err != nil {
		return nil, err
	}

	rpmes, err := ms.SendRequest(ctx, pmes)
	if err != nil {
		return nil, err
	}

	if rpmes == nil {
		err = fmt.Errorf("no response from %s\n", p)
		return nil, err
	}
	fmt.Printf("service/service.go SendRequest: Received request response from %s\n", p)
	err = srv.handleError(rpmes)
	if err != nil {
		return nil, err
	}

	return rpmes, nil
}

// SendHTTPStreamRequest sends a request over HTTP
func (srv *Service) SendHTTPStreamRequest(addr string, pmes *pb.Envelope, access string) (chan *pb.Envelope, chan error, *func()) {
	envCh := make(chan *pb.Envelope)
	errCh := make(chan error)

	var cancel func()
	go func() {
		defer close(envCh)

		payload, err := proto.Marshal(pmes)
		if err != nil {
			errCh <- err
			return
		}

		req, err := http.NewRequest("POST", addr, bytes.NewReader(payload))
		if err != nil {
			errCh <- err
			return
		}
		req.Header.Set("Authorization", "Basic "+access)

		tr := &http.Transport{}
		client := &http.Client{Transport: tr}

		cancel = func() {
			tr.CancelRequest(req)
		}

		res, err := client.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		defer res.Body.Close()

		if res.StatusCode >= 400 {
			res, err := util.UnmarshalString(res.Body)
			if err != nil {
				errCh <- err
			} else {
				errCh <- fmt.Errorf(res)
			}
			return
		}

		reader := bufio.NewReader(res.Body)
		for {
			size := make([]byte, 2)
			_, err := io.ReadFull(reader, size)
			if err == io.EOF {
				return
			} else if err != nil {
				errCh <- err
				return
			}

			mes := make([]byte, binary.LittleEndian.Uint16(size))
			_, err = io.ReadFull(reader, mes)
			if err == io.EOF {
				return
			} else if err != nil {
				errCh <- err
				return
			}

			env := new(pb.Envelope)
			err = proto.Unmarshal(mes, env)
			if err != nil {
				errCh <- err
				return
			}
			if env.Message == nil {
				errCh <- fmt.Errorf("message is nil")
				return
			}

			err = srv.handleError(env)
			if err != nil {
				errCh <- err
				return
			}
			envCh <- env
		}
	}()

	return envCh, errCh, &cancel
}

// SendMessage sends out a message
func (srv *Service) SendMessage(ctx context.Context, p string, pmes *pb.Envelope) error {
	pid, err := peer.IDB58Decode(p)
	if err != nil {
		return err
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()
	}

	ms, err := srv.messageSenderForPeer(ctx, pid)
	if err != nil {
		return err
	}

	return ms.SendMessage(ctx, pmes)
}

// NewEnvelope returns a signed pb message for transport
func (srv *Service) NewEnvelope(mtype pb.Message_Type, msg proto.Message, id *int32, response bool) (*pb.Envelope, error) {
	var payload *any.Any
	if msg != nil {
		var err error
		payload, err = ptypes.MarshalAny(msg)
		if err != nil {
			return nil, err
		}
	}
	message := &pb.Message{Type: mtype, Payload: payload}

	if id != nil {
		message.Request = *id
	}
	if response {
		message.Response = true
	}

	ser, err := proto.Marshal(message)
	//_, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	sig, err := srv.privateKey.Sign(ser)
	if err != nil {
		return nil, err
	}

	return &pb.Envelope{Message: message, Sig: sig}, nil
}

// NewError returns a signed pb error message
func (srv *Service) NewError(code int, msg string, id int32) (*pb.Envelope, error) {
	return srv.NewEnvelope(pb.Message_ERROR, &pb.Error{
		Code:    uint32(code),
		Message: msg,
	}, &id, true)
}

// VerifyEnvelope verifies the authenticity of an envelope
func (srv *Service) VerifyEnvelope(env *pb.Envelope, pid peer.ID) error {
	ser, err := proto.Marshal(env.Message)
	if err != nil {
		return err
	}

	pk, err := pid.ExtractPublicKey()
	if err != nil {
		return err
	}

	return crypto.Verify(pk, ser, env.Sig)
}

// handleError receives an error response
func (srv *Service) handleError(env *pb.Envelope) error {
	if env.Message.Payload == nil && env.Message.Type != pb.Message_PONG {
		return fmt.Errorf("message payload with type %s is nil", env.Message.Type.String())
	}

	if env.Message.Type != pb.Message_ERROR {
		return nil
	} else {
		errMsg := new(pb.Error)
		err := ptypes.UnmarshalAny(env.Message.Payload, errMsg)
		if err != nil {
			return err
		}
		return fmt.Errorf(errMsg.Message)
	}
}

// handleCore provides service level handlers for common message types
func (srv *Service) handleCore(mtype pb.Message_Type) func(*pb.Envelope, peer.ID) (*pb.Envelope, error) {
	switch mtype {
	case pb.Message_PING:
		return srv.handlePing
	default:
		return nil
	}
}

// handlePing receives a PING message
func (srv *Service) handlePing(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	return srv.NewEnvelope(pb.Message_PONG, nil, &env.Message.Request, true)
}

var dhtReadMessageTimeout = time.Minute
var dhtStreamIdleTimeout = 10 * time.Minute
var ErrReadTimeout = fmt.Errorf("timed out reading response")

// The Protobuf writer performs multiple small writes when writing a message.
// We need to buffer those writes, to make sure that we're not sending a new
// packet for every single write.
type bufferedDelimitedWriter struct {
	*bufio.Writer
	ggio.WriteCloser
}

var writerPool = sync.Pool{
	New: func() interface{} {
		w := bufio.NewWriter(nil)
		return &bufferedDelimitedWriter{
			Writer:      w,
			WriteCloser: ggio.NewDelimitedWriter(w),
		}
	},
}

func writeMsg(w io.Writer, mes *pb.Envelope) error {
	bw := writerPool.Get().(*bufferedDelimitedWriter)
	bw.Reset(w)
	tStart:=time.Now()
	err := bw.WriteMsg(mes)
	if err == nil {
		err = bw.Flush()
		fmt.Println("=== bw flush: ",time.Since(tStart))
	}
	bw.Reset(nil)
	writerPool.Put(bw)
	return err
}

func (w *bufferedDelimitedWriter) Flush() error {
	return w.Writer.Flush()
}

// handleNewStream implements the inet.StreamHandler
func (srv *Service) handleNewStream(s inet.Stream) {
	defer s.Reset()
	if srv.handleNewMessage(s) {
		// Gracefully close the stream for writes.
		_ = s.Close()
	}
}

func (srv *Service) handleNewMessage(s inet.Stream) bool {
	//ctx := srv.Node().Context()

	r := msgio.NewVarintReaderSize(s, inet.MessageSizeMax)

	mPeer := s.Conn().RemotePeer()

	timer := time.AfterFunc(dhtStreamIdleTimeout, func() { s.Reset() })
	defer timer.Stop()

	for {
		var req pb.Envelope
		msgbytes, err := r.ReadMsg()
		if err != nil {
			defer r.ReleaseMsg(msgbytes)
			if err == io.EOF {
				return true
			}
			// This string test is necessary because there isn't a single stream reset error
			// instance	in use.
			if err.Error() != "stream reset" {
				//log.Debugf("error reading message: %#v", err)
			}
			return false
		}
		err = proto.Unmarshal(msgbytes, &req)
		r.ReleaseMsg(msgbytes)
		if err != nil {
			//log.Debugf("error unmarshalling message: %#v", err)
			return false
		}

		timer.Reset(dhtStreamIdleTimeout)

		if err := srv.VerifyEnvelope(&req, mPeer); err != nil {
			//log.Warningf("error verifying message: %s", err)
			fmt.Printf("Error when verify message: %s\n", err)
			continue
		}

        handler := srv.handler.Handle

		//log.Debugf("received %s from %s", req.Message.Type.String(), mPeer.Pretty())
		rpmes, err := handler(&req, mPeer)
		if err != nil {
			fmt.Printf("error handling message %s: %s\n", req.Message.Type.String(), err)
			return false
		}

		if rpmes == nil {
			continue
		}

		// send out response msg
		//log.Debugf("responding with %s to %s", rpmes.Message.Type.String(), mPeer.Pretty())

		// send out response msg
		err = writeMsg(s, rpmes)
		if err != nil {
			//log.Debugf("error writing response: %s", err)
			return false
		}
	}
}
