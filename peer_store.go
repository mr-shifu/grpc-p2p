package p2p

import "sync"

type PeerStore struct {
	lock  sync.RWMutex
	peers map[string]*Peer
}

func NewPeerStore() *PeerStore {
	return &PeerStore{
		lock:  sync.RWMutex{},
		peers: make(map[string]*Peer),
	}
}

func (ps *PeerStore) AddPeer(peer *Peer) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.peers[peer.Name] = peer
}

func (ps *PeerStore) RemovePeer(peer *Peer) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	delete(ps.peers, peer.Name)
}

func (ps *PeerStore) GetPeer(name string) *Peer {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	return ps.peers[name]
}

func (ps *PeerStore) GetAllPeers() []*Peer {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	var peers []*Peer
	for _, peer := range ps.peers {
		peers = append(peers, peer)
	}
	return peers
}
