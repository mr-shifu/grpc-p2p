package p2p

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

func (ps *PeerService) AddPeer(peer *Peer) {
	ps.peerstore.AddPeer(peer)
}

func (ps *PeerService) RemovePeer(peer *Peer) {
	ps.peerstore.RemovePeer(peer)
}

func (ps *PeerService) GetPeer(peer *Peer) *Peer {
	return ps.peerstore.GetPeer(peer.Name)
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
