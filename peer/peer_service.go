package peer

import (
	"context"

	"github.com/mr-shifu/grpc-p2p/config"
	p2p_pb "github.com/mr-shifu/grpc-p2p/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PeerService struct {
	self      *config.Peer
	bootstrap []config.Peer
	peerstore *PeerStore
	discovery *Discovery

	p2p_pb.UnimplementedPeerServiceServer
}

func NewPeerService(cfg *config.Config) *PeerService {
	peerstore := NewPeerStore()
	discovery := NewDiscovery(peerstore)

	ps := &PeerService{
		self:      &cfg.Local,
		bootstrap: cfg.Bootstrap,
		peerstore: peerstore,
		discovery: discovery,
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

func (ps *PeerService) GetPeers(context.Context, *p2p_pb.GetPeersRequest) (*p2p_pb.GetPeersResponse, error) {
	peers := ps.peerstore.GetAllPeers()
	pbPeers := peersToPbPeers(peers)
	return &p2p_pb.GetPeersResponse{
		Peers: pbPeers,
	}, nil
}

func peersToPbPeers(peers []*Peer) []*p2p_pb.Peer {
	var pbPeers []*p2p_pb.Peer
	for _, peer := range peers {
		p := &p2p_pb.Peer{
			Name:        peer.Name,
			Address:     peer.Addr,
			CLusterName: peer.ClusterName,
		}
		pbPeers = append(pbPeers, p)
	}
	return pbPeers
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
