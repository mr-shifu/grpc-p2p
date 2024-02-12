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
	addr, err := validatePeerAddr(addr)
	if err != nil {
		return false, ErrInvalidPeerAddress
	}

	return ps.exists(addr), nil
}

func (ps *PeerStore) AddPeer(p *Peer) error {
	addr, err := validatePeerAddr(p.Addr())
	if err != nil {
		return ErrInvalidPeerAddress
	}

	if exists := ps.exists(addr); exists {
		return ErrPeerAlreadyExists
	}

	np := NewPeer(addr, p.Attributes())
	ps.addPeer(np)

	return nil
}

func (ps *PeerStore) AddPeers(peers []*Peer, skip_errors bool) (int, error) {
	count := 0
	for _, peer := range peers {
		if err := ps.AddPeer(peer); err != nil {
			if skip_errors || err == ErrPeerAlreadyExists {
				continue
			} else {
				return count, err
			}
		}
		count++
	}
	return count, nil
}

func (ps *PeerStore) UpdatePeer(p *Peer) error {
	addr, err := validatePeerAddr(p.Addr())
	if err != nil {
		return ErrInvalidPeerAddress
	}

	if exists := ps.exists(addr); !exists {
		return ErrPeerNotFouund
	}

	np := NewPeer(addr, p.Attributes())
	ps.addPeer(np)

	return nil
}

func (ps *PeerStore) RemovePeer(p *Peer) error {
	addr, err := validatePeerAddr(p.Addr())
	if err != nil {
		return ErrInvalidPeerAddress
	}

	if exists := ps.exists(addr); !exists {
		return ErrPeerNotFouund
	}

	ps.removePeer(addr)

	return nil
}

func (ps *PeerStore) GetPeer(addr string) (*Peer, error) {
	addr, err := validatePeerAddr(addr)
	if err != nil {
		return nil, ErrInvalidPeerAddress
	}

	peerInfo, err := ps.getPeer(addr)
	if err != nil {
		return nil, err
	}

	p := NewPeer(peerInfo.Addr, peerInfo.Attributes)
	p.SetConnection(ps.conns[addr])

	return p, nil
}

func (ps *PeerStore) GetAllPeers() []*Peer {
	return ps.getAllPeers()
}

func (ps *PeerStore) GetPeerConnection(addr string) (*grpc.ClientConn, error) {
	addr, err := validatePeerAddr(addr)
	if err != nil {
		return nil, ErrInvalidPeerAddress
	}
	if exists := ps.exists(addr); !exists {
		return nil, ErrPeerNotFouund
	}

	return ps.getPeerConnection(addr)
}

func (ps *PeerStore) SetPeerConnection(addr string, conn *grpc.ClientConn) (*grpc.ClientConn, error) {
	addr, err := validatePeerAddr(addr)
	if err != nil {
		return nil, ErrInvalidPeerAddress
	}
	if exists := ps.exists(addr); !exists {
		return nil, ErrPeerNotFouund
	}

	return ps.setPeerConnection(addr, conn)
}

func (ps *PeerStore) exists(addr string) bool {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	_, ok := ps.peers[addr]
	return ok
}

func (ps *PeerStore) getPeer(addr string) (*PeerInfo, error) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if _, ok := ps.peers[addr]; !ok {
		return nil, ErrPeerNotFouund
	}
	return ps.peers[addr], nil
}

func (ps *PeerStore) getAllPeers() []*Peer {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	var peers []*Peer
	for _, peer := range ps.peers {
		p := NewPeer(peer.Addr, peer.Attributes)
		p.SetConnection(ps.conns[peer.Addr])
		peers = append(peers, p)
	}
	return peers
}

func (ps *PeerStore) addPeer(p *Peer) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.peers[p.Addr()] = p.PeerInfo
}

func (ps *PeerStore) removePeer(addr string) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	delete(ps.peers, addr)
}

func (ps *PeerStore) getPeerConnection(addr string) (*grpc.ClientConn, error) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	return ps.conns[addr], nil
}

func (ps *PeerStore) setPeerConnection(addr string, conn *grpc.ClientConn) (*grpc.ClientConn, error) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.conns[addr] = conn
	return conn, nil
}

func validatePeerAddr(addr string) (string, error) {
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
