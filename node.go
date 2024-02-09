package p2p

import (
	"context"
	"net"
	"time"

	"github.com/mr-shifu/grpc-p2p/config"
	"github.com/mr-shifu/grpc-p2p/peer"
	p2p_pb "github.com/mr-shifu/grpc-p2p/proto"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Node struct {
	// The gRPC server to connect to and receive messages from the server
	server *grpc.Server

	// Node local config including name, cluster name, address to listen for incoming connections
	local *config.Peer

	//
	peerService *peer.PeerService

	// Logger to use for debug logging
	logger zerolog.Logger
}

// NewNode creates a new node with the given address and logger
func NewNode(cfgpath string, logger zerolog.Logger) *Node {
	cfg, err := config.FromFile(cfgpath)
	if err != nil {
		return nil
	}

	// create new grpc server
	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	// instantiate a new peer service
	ps := peer.NewPeerService(cfg)

	// register service
	p2p_pb.RegisterPeerServiceServer(server, ps)

	// enable rpc reflection
	reflection.Register(server)

	// start discovery
	go ps.StartDiscovery()

	return &Node{
		local:       &cfg.Local,
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
	ln, err := net.Listen("tcp", n.Self().Addr)
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		n.logger.Debug().Msgf("Node Started with address %s", n.Self().Addr)
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

func (n *Node) Self() *config.Peer {
	return n.local
}