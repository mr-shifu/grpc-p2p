package p2p

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PeerService struct {
	self      *Peer
	peerstore *PeerStore
}

func NewPeerService(self *Peer) *PeerService {
	return &PeerService{
		self:      self,
		peerstore: NewPeerStore(),
	}
}

func (ps *PeerService) AddPeer(peer *Peer) error {
	return ps.peerstore.AddPeer(peer)
}

func (ps *PeerService) RemovePeer(peer *Peer) error {
	return ps.peerstore.RemovePeer(peer)
}

func (ps *PeerService) GetPeer(name string) (*Peer, error) {
	return ps.peerstore.GetPeer(name)
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
	peer, err := ps.GetPeer(p.Addr)
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
	peers := ps.GetPeers()
	for _, peer := range peers {
		if err := ps.Disconnect(peer); err != nil {
			return err
		}
	}
	return nil
}
