package p2p

import (
	"context"
	"net"
	"time"

	p2p_pb "github.com/mr-shifu/grpc-p2p/proto"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Node struct {
	// The gRPC server to connect to and receive messages from the server
	server *grpc.Server

	// Node address to listen for incoming connections
	addr string

	//
	peerService *PeerService

	// Logger to use for debug logging
	logger zerolog.Logger
}

// NewNode creates a new node with the given address and logger
func NewNode(addr string, logger zerolog.Logger) *Node {
	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	self := NewPeer(addr, "", "")
	ps := NewPeerService(self)

	// register service
	p2p_pb.RegisterPeerServiceServer(server, ps)

	// enable rpc reflection
	reflection.Register(server)

	return &Node{
		addr:        addr,
		server:      server,
		peerService: ps,
		logger:      logger,
	}
}

// Start starts the node and listens for incoming connections
// It returns an error if the node fails to start
// It waits for receiving a signal from ctx to stop server grcefully
// It forces server to stop if gracefully shutdown failed ad returns an error
func (n *Node) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", n.addr)
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		n.logger.Debug().Msgf("Node Started with address %s", n.addr)
		if err := n.server.Serve(ln); err != nil {
			n.logger.Error().Err(err).Msg("Server failed to start")
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		return n.Stop()
	})
	return group.Wait()
}

// Stop first tries to gracefully shutdown the server and if it timed out then
// it forces the server to stop and returns an error
func (n *Node) Stop() error {
	n.logger.Debug().Msg("Gracefully Shutdown of Node")

	stopped := make(chan struct{})
	go func() {
		n.server.GracefulStop()
		close(stopped)
	}()

	timeout := time.NewTimer(5 * time.Second)

	select {
	case <-timeout.C:
		n.logger.Debug().Msg("Forcing Shutdown of Node")
		n.server.Stop()
		return ErrServerGracefullyShutdownTimedout
	case <-stopped:
		n.logger.Debug().Msg("Node Gracefully Shutdown Successfully")
		return nil
	}
}
