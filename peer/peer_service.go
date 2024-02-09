package peer

import (
	"context"
	"errors"

	"github.com/mr-shifu/grpc-p2p/config"
	p2p_pb "github.com/mr-shifu/grpc-p2p/proto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type PeerService struct {
	self      *config.Peer
	bootstrap []config.Peer
	peerstore *PeerStore

	logger zerolog.Logger
}

func NewPeerService(cfg *config.Config, logger zerolog.Logger) *PeerService {
	peerstore := NewPeerStore()

	ps := &PeerService{
		self:      &cfg.Local,
		bootstrap: cfg.Bootstrap,
		peerstore: peerstore,
		logger:    logger,
	}

	// add bootstrap nodes into peerstore
	for _, peer := range cfg.Bootstrap {
		if peer.Addr != ps.self.Addr {
			ps.peerstore.AddPeer(&Peer{
				Name:        peer.Name,
				ClusterName: peer.ClusterName,
				Addr:        peer.Addr,
			})
		}
	}

	return ps
}

func (ps *PeerService) AddPeer(p *Peer) error {
	return ps.peerstore.AddPeer(p)
}

func (ps *PeerService) GetPeers() []*Peer {
	return ps.peerstore.GetAllPeers()
}

func (ps *PeerService) GetClusterPeers() []*Peer {
	peers := ps.peerstore.GetAllPeers()
	var clusterPeers []*Peer
	for _, peer := range peers {
		if peer.ClusterName == ps.self.ClusterName {
			clusterPeers = append(clusterPeers, peer)
		}
	}
	return clusterPeers
}

func (ps *PeerService) GetNeighbors(ctx context.Context, p *Peer) ([]*Peer, error) {
	if p.State() != connectivity.Ready {
		_, err := ps.Connect(p)
		if err != nil {
			return nil, err
		}
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "cluster_name", ps.self.ClusterName, "peer_name", ps.self.Name)

	client := p2p_pb.NewPeerServiceClient(p.conn)
	neighbors, err := client.GetPeers(ctx, &p2p_pb.GetPeersRequest{})
	if err != nil {
		return nil, err
	}

	np := peersFromPbPeers(neighbors.Peers)

	return np, nil
}

// Connect connects to a peer and returns a client connection and updates peer connection at peerstore
// throws error if connection fails
func (ps *PeerService) Connect(p *Peer) (*grpc.ClientConn, error) {
	// connect to peer
	conn, err := grpc.Dial(p.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// update peer connection at peerstore
	p.conn = conn
	ps.peerstore.AddPeer(p)

	return conn, err
}

// Disconnect disconnects from a peer and updates peer connection at peerstore
// throws error if disconnection fails
func (ps *PeerService) Disconnect(p *Peer) error {
	if p.conn == nil {
		return errors.New("peer not connected")
	}
	if err := p.conn.Close(); err != nil {
		return err
	}

	// update peer connection at peerstore
	p.conn = nil
	ps.peerstore.AddPeer(p)

	return nil
}

func (ps *PeerService) DisconnectAll() error {
	peers := ps.peerstore.GetAllPeers()
	for _, peer := range peers {
		if err := ps.Disconnect(peer); err != nil {
			return err
		}
	}
	return nil
}

func peersFromPbPeers(pbPeers []*p2p_pb.Peer) []*Peer {
	var peers []*Peer
	for _, peer := range pbPeers {
		p := &Peer{
			Name:        peer.Name,
			Addr:        peer.Address,
			ClusterName: peer.CLusterName,
		}
		peers = append(peers, p)
	}
	return peers
}
