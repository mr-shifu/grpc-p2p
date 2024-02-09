package peer

import (
	"context"

	"github.com/mr-shifu/grpc-p2p/config"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PeerService struct {
	self      *config.Peer
	bootstrap []config.Peer
	peerstore *PeerStore
	discovery *Discovery

	logger zerolog.Logger
}

func NewPeerService(cfg *config.Config, logger zerolog.Logger) *PeerService {
	peerstore := NewPeerStore()
	discovery := NewDiscovery(peerstore, logger)

	ps := &PeerService{
		self:      &cfg.Local,
		bootstrap: cfg.Bootstrap,
		peerstore: peerstore,
		discovery: discovery,
		logger:    logger,
	}

	// add bootstrap nodes into peerstore
	for _, peer := range cfg.Bootstrap {
		ps.peerstore.AddPeer(&Peer{
			Name:        peer.Name,
			ClusterName: peer.ClusterName,
			Addr:        peer.Addr,
		})
	}

	return ps
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

func (ps *PeerService) Connect(p *Peer) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(p.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (ps *PeerService) Disconnect(p *Peer) error {
	peer, err := ps.peerstore.GetPeer(p.Addr)
	if err != nil {
		return err
	}
	if peer == nil {
		return nil
	}
	if peer.conn != nil {
		return peer.conn.Close()
	}
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

func (ps *PeerService) StartDiscovery() error {
	return ps.discovery.Start(context.Background())
}
