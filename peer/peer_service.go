package peer

import (
	"context"
	"errors"

	"github.com/mr-shifu/grpc-p2p/config"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type PeerService struct {
	self      *Peer
	bootstrap []config.Peer
	peerstore *PeerStore
	client    *Client

	logger zerolog.Logger
}

func NewPeerService(cfg *config.Config, logger zerolog.Logger) *PeerService {
	store := NewPeerStore()

	self := NewPeer(cfg.Local.Addr, cfg.Local.Attributes)

	ps := &PeerService{
		self:      self,
		bootstrap: cfg.Bootstrap,
		peerstore: store,
		client:    NewClient(),
		logger:    logger,
	}

	// add bootstrap nodes into peerstore
	for _, peer := range cfg.Bootstrap {
		if peer.Addr != ps.self.Addr() {
			p := NewPeer(peer.Addr, peer.Attributes)
			ps.peerstore.AddPeer(p)
		}
	}

	return ps
}

func (ps *PeerService) Self() *Peer {
	return ps.self
}

func (ps *PeerService) AddPeer(p *Peer) error {
	return ps.peerstore.AddPeer(p)
}

func (ps *PeerService) AddPeers(peers []*Peer) ([]*Peer, error) {
	return ps.peerstore.AddPeers(peers, true)
}

func (ps *PeerService) GetPeers() []*Peer {
	return ps.peerstore.GetPeers()
}

func (ps *PeerService) GetPeersWithAttributes(attrs map[string]string) []*Peer {
	return ps.peerstore.GetPeersWithAttributes(attrs)
}

func (ps *PeerService) GetState(addr string) (PeerState, error) {
	p, err := ps.peerstore.GetPeer(addr)
	if err != nil {
		return 0, err
	}
	return p.GetState(), nil
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

	neihgbors, err := ps.client.GetPeers(ctx, conn, ps.self.PeerInfo)
	if err != nil {
		return nil, err
	}

	return neihgbors, nil
}

// Connect connects to a peer and returns a client connection and updates peer connection at peerstore
// throws error if connection fails
func (ps *PeerService) Connect(addr string) (*grpc.ClientConn, error) {
	if addr == ps.self.Addr() {
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
	peers := ps.peerstore.GetPeers()
	for _, peer := range peers {
		if err := ps.Disconnect(peer.Addr()); err != nil {
			return err
		}
	}
	return nil
}
