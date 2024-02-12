package peer

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

var (
	ErrInvalidPeerAddress = errors.New("peerstore: invalid peer address")
	ErrPeerNotFouund      = errors.New("peerstore: peer not found")
	ErrPeerAlreadyExists  = errors.New("peerstore: failed to add peer. peer already exists")
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

func (ps *PeerStore) Exists(addr string) (bool, error) {
	addr, err := ps.validatePeerAddr(addr)
	if err != nil {
		return false, err
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	_, ok := ps.peers[addr]
	return ok, nil
}

func (ps *PeerStore) AddPeer(peer *Peer) error {
	exists, err := ps.Exists(peer.Addr())
	if err != nil {
		return err
	}
	if exists {
		return ErrPeerAlreadyExists
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
	exists, err := ps.Exists(peer.Addr())
	if err != nil {
		return err
	}
	if exists {
		return ErrPeerNotFouund
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.peers[peer.Addr()] = &PeerInfo{
		Attributes: peer.Attributes(),
	}
	return nil
}

func (ps *PeerStore) RemovePeer(peer *Peer) error {
	exists, err := ps.Exists(peer.Addr())
	if err != nil {
		return err
	}
	if exists {
		return ErrPeerNotFouund
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	delete(ps.peers, peer.Addr())
	return nil
}

func (ps *PeerStore) GetPeer(addr string) (*Peer, error) {
	exists, err := ps.Exists(addr)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPeerNotFouund
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

func (ps *PeerStore) validatePeerAddr(addr string) (string, error) {
	host, err := parseIP(addr)
	if err == nil && host != "" {
		return host, nil
	}

	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	host = strings.TrimPrefix(u.Host, "www.")

	return host, nil
}

func parseIP(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		if ip := net.ParseIP(addr); ip != nil {
			return ip.String(), nil
		}
		return "", ErrInvalidPeerAddress
	}
	if ip := net.ParseIP(host); ip == nil {
		return "", ErrInvalidPeerAddress
	}
	if p, err := net.LookupPort("tcp", port); err != nil || p < 0 {
		return "", ErrInvalidPeerAddress
	}
	return fmt.Sprintf("%s:%s", host, port), nil
}
