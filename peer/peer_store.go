package peer

import (
	"errors"
	"sync"

	"google.golang.org/grpc"
)

var (
	ErrInvalidPeerAddress     = errors.New("invalid peer address")
	ErrInvalidPeerName        = errors.New("invalid peer name")
	ErrInvalidPeerClusterName = errors.New("invalid peer cluster name")
	ErrPeerNotFouund          = errors.New("peer not found")
	ErrPeerAlreadyExists      = errors.New("peer already exists")
)



type PeerStore struct {
	lock  sync.RWMutex
	peers map[string]*PeerInfo
	conns map[string]*grpc.ClientConn
}

func NewPeerStore() *PeerStore {
	return &PeerStore{
		lock:  sync.RWMutex{},
		peers: make(map[string]*PeerInfo),
		conns: make(map[string]*grpc.ClientConn),
	}
}

func (ps *PeerStore) Exists(peer *Peer) bool {
	if peer.Addr() == "" {
		return false
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	_, ok := ps.peers[peer.Addr()]
	return ok
}

func (ps *PeerStore) AddPeer(peer *Peer) error {
	if ps.Exists(peer) {
		return ErrPeerAlreadyExists
	}
	if peer.Addr() == "" {
		return ErrInvalidPeerAddress
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	p := NewPeer(peer.Addr(), peer.Attributes())
	ps.peers[peer.Addr()] = p.PeerInfo

	return nil
}

func (ps *PeerStore) AddPeers(peers ...*Peer) error {
	for _, peer := range peers {
		if err := ps.AddPeer(peer); err != nil {
			return err
		}
	}
	return nil
}

func (ps *PeerStore) UpdatePeer(peer *Peer) error {
	if !ps.Exists(peer) {
		return ErrPeerNotFouund
	}
	if peer.Addr() == "" {
		return ErrInvalidPeerAddress
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.peers[peer.Addr()] = &PeerInfo{
		Attributes: peer.Attributes(),
	}
	return nil
}

func (ps *PeerStore) RemovePeer(peer *Peer) error {
	if peer.Addr() == "" {
		return ErrInvalidPeerAddress
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	delete(ps.peers, peer.Addr())
	return nil
}

func (ps *PeerStore) GetPeer(addr string) (*Peer, error) {
	if addr == "" {
		return nil, ErrInvalidPeerAddress
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	if _, ok := ps.peers[addr]; !ok {
		return nil, ErrPeerNotFouund
	}
	p := NewPeer(ps.peers[addr].Addr, ps.peers[addr].Attributes)
	p.SetConnection(ps.conns[addr])

	return p, nil
}

func (ps *PeerStore) GetAllPeers() []*Peer {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	var peers []*Peer
	for _, peer := range ps.peers {
		p := NewPeer(peer.Addr, peer.Attributes)
		peers = append(peers, p)
	}
	return peers
}

func (ps *PeerStore) GetPeerConnection(addr string) (*grpc.ClientConn, error) {
	if addr == "" {
		return nil, ErrInvalidPeerAddress
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	return ps.conns[addr], nil
}

func (ps *PeerStore) SetPeerConnection(addr string, conn *grpc.ClientConn) (*grpc.ClientConn, error) {
	if addr == "" {
		return nil, ErrInvalidPeerAddress
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.conns[addr] = conn
	return conn, nil
}
