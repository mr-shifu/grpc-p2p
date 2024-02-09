package peer

import (
	"errors"
	"sync"
)

var (
	ErrInvalidPeerAddress     = errors.New("invalid peer address")
	ErrInvalidPeerName        = errors.New("invalid peer name")
	ErrInvalidPeerClusterName = errors.New("invalid peer cluster name")
)

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

func (ps *PeerStore) AddPeer(peer *Peer) error {
	if peer.Addr == "" {
		return ErrInvalidPeerAddress
	}
	if peer.Name == "" {
		return ErrInvalidPeerName
	}
	if peer.ClusterName == "" {
		return ErrInvalidPeerClusterName
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.peers[peer.Addr] = peer
	return nil
}

func (ps *PeerStore) RemovePeer(peer *Peer) error {
	if peer.Addr == "" {
		return ErrInvalidPeerAddress
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	delete(ps.peers, peer.Addr)
	return nil
}

func (ps *PeerStore) GetPeer(addr string) (*Peer, error) {
	if addr == "" {
		return nil, ErrInvalidPeerAddress
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	return ps.peers[addr], nil
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
