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
	store := NewPeerStore()

	ps := &PeerService{
		self:      &cfg.Local,
		bootstrap: cfg.Bootstrap,
		peerstore: store,
		logger:    logger,
	}

	// add bootstrap nodes into peerstore
	for _, peer := range cfg.Bootstrap {
		if peer.Addr != ps.self.Addr {
			p := NewPeer(peer.Addr, peer.Attributes)
			ps.peerstore.AddPeer(p)
		}
	}

	return ps
}

func (ps *PeerService) AddPeer(p *Peer) error {
	return ps.peerstore.AddPeer(p)
}

func (ps *PeerService) AddPeers(peers []*Peer) error {
	return ps.peerstore.AddPeers(peers...)
}

func (ps *PeerService) GetPeers() []*Peer {
	return ps.peerstore.GetAllPeers()
}

func (ps *PeerService) GetNeighbors(ctx context.Context, p *Peer) ([]*Peer, error) {
	conn, err := ps.Connect(p.Addr())
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, errors.New("connection failed")
	}
	if conn.GetState() != connectivity.Ready {
		return nil, errors.New("connection not ready")
	}

	var opts []string
	opts = append(opts, "addr", ps.self.Addr)
	for k, v := range ps.self.Attributes {
		key := "attr-" + k
		opts = append(opts, key, v)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, opts...)

	client := p2p_pb.NewPeerServiceClient(conn)
	neighbors, err := client.GetPeers(ctx, &p2p_pb.GetPeersRequest{})
	if err != nil {
		return nil, err
	}

	np := peersFromPbPeers(neighbors.Peers)

	return np, nil
}

// Connect connects to a peer and returns a client connection and updates peer connection at peerstore
// throws error if connection fails
func (ps *PeerService) Connect(addr string) (*grpc.ClientConn, error) {
	if addr == ps.self.Addr {
		return nil, errors.New("cannot connect to self")
	}

	p, err := ps.peerstore.GetPeer(addr)
	if err != nil {
		return nil, err
	}
	if p.GetState() == Ready {
		return p.conn, nil
	}

	// connect to peer
	conn, err := grpc.Dial(p.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	ps.peerstore.SetPeerConnection(p.Addr(), conn)

	return conn, err
}

// Disconnect disconnects from a peer and updates peer connection at peerstore
// throws error if disconnection fails
func (ps *PeerService) Disconnect(addr string) error {
	p, err := ps.peerstore.GetPeer(addr)
	if err != nil {
		return err
	}

	// update peer connection at peerstore
	if p.GetState() == Ready {
		p.conn.Close()
		ps.peerstore.SetPeerConnection(p.Addr(), nil)
	}
	return nil
}

func (ps *PeerService) DisconnectAll() error {
	peers := ps.peerstore.GetAllPeers()
	for _, peer := range peers {
		if err := ps.Disconnect(peer.Addr()); err != nil {
			return err
		}
	}
	return nil
}

func peersFromPbPeers(pbPeers []*p2p_pb.Peer) []*Peer {
	var peers []*Peer
	for _, peer := range pbPeers {
		attrs := make(map[string]string)
		for _, attr := range peer.Attributes {
			attrs[attr.Key] = attr.Value
		}
		p := NewPeer(peer.Address, attrs)
		peers = append(peers, p)
	}
	return peers
}
