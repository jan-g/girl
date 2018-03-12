package server

import (
	"context"
	"io"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"errors"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/jan-g/girl/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server interface {
	AddPeer(address string) error
	Serve()
	Start(context.Context)
}

func NewServer(network string, address string, spi model.LimiterSPI) (Server, error) {
	lst, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer()
	server := &server{
		name:    address,
		lst:     lst,
		srv:     grpcServer,
		limiter: spi,
	}
	server.peers.Store([]string{})
	model.RegisterGossipProtocolServer(grpcServer, server)
	model.RegisterLimiterProtocolServer(grpcServer, server)

	return server, nil
}

type server struct {
	name string

	lst net.Listener
	srv *grpc.Server

	limiter model.LimiterSPI

	peers atomic.Value
}

func (s *server) AddPeer(address string) error {
	peers := s.peers.Load().([]string)
	peers = append(peers, address)
	s.peers.Store(peers)
	return nil
}

func (s *server) Serve() {
	s.srv.Serve(s.lst)
}

func (s *server) Start(ctx context.Context) {
	go s.gossip(ctx)
}

// Handle incoming gossip.
func (s *server) Gossip(srv model.GossipProtocol_GossipServer) error {
	// Three-way handshake
	v, err := srv.Recv()
	if err != nil {
		return err
	}
	theyHave, ok := v.Hs.(*model.ConnectorHandshake_IHave)
	if !ok {
		return errors.New("Expected IHave")
	}
	incomingEpoch := theyHave.IHave.Epoch
	logrus.WithField("epoch", incomingEpoch).Debug("Epoch received")
	response := s.limiter.GossipIn(theyHave.IHave)

	err = srv.Send(response)
	if err != nil {
		return err
	}

	v, err = srv.Recv()
	if err != nil {
		return err
	}
	theyPush, ok := v.Hs.(*model.ConnectorHandshake_Push)
	if !ok {
		return errors.New("Expected Push")
	}
	incomingPush := theyPush.Push
	s.limiter.ReceivePush(incomingEpoch, incomingPush.Traffic)

	return nil
}

func (s *server) Sync(context.Context, *google_protobuf.Empty) (*google_protobuf.Empty, error) {
	return &google_protobuf.Empty{}, nil
}

// Handle periodic ticks to manage the current epoch
// Handle periodic ticks to manage outgoing gossip
func (s *server) gossip(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			epoch := time.Now().Unix()
			if epoch != s.limiter.Epoch() {
				// The protocol relies on this being fairly accurate at the moment
				s.limiter.Tick(epoch)
			}
			peers := s.peers.Load().([]string)
			addr := peers[r.Intn(len(peers))]
			// c, can := context.WithTimeout(ctx, 90*time.Millisecond)
			s.gossipWith(addr, ctx)
			// can()
		}
	}
}

func (s *server) gossipWith(address string, ctx context.Context) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := model.NewGossipProtocolClient(conn)
	goss, err := client.Gossip(ctx)
	if err != nil {
		return err
	}

	iHave := s.limiter.GossipOut()

	logrus.WithField("peer", address).Debug("Sending initial handshake")
	err = goss.Send(&model.ConnectorHandshake{
		Hs: &model.ConnectorHandshake_IHave{IHave: iHave},
	})
	if err != nil {
		return err
	}
	logrus.WithField("peer", address).Debug("Receiving response")
	ret, err := goss.Recv()
	if err != nil {
		return err
	}
	s.limiter.ReceivePush(iHave.Epoch, ret.Push.Traffic)

	push := &model.Push{
		Traffic: s.limiter.OriginatePush(iHave.Epoch, ret.IWant),
	}
	logrus.WithField("peer", address).Debug("Sending final handshake")
	err = goss.Send(&model.ConnectorHandshake{
		Hs: &model.ConnectorHandshake_Push{Push: push},
	})
	if err != nil {
		return err
	}
	err = goss.RecvMsg(nil)
	if err != io.EOF {
		panic(err)
	}
	goss.CloseSend()
	logrus.WithField("peer", address).Debug("Gossip complete")

	return nil
}

func (s *server) Use(ctx context.Context, rq *model.UseRequest) (*model.UseResponse, error) {
	f := rq.Facet
	q := s.limiter.(model.Limiter).Ask(f, rq.Quantity)
	r := s.limiter.(model.Limiter).Level(f)
	return &model.UseResponse{Facet: f, Quantity: q, Remaining: r}, nil
}
