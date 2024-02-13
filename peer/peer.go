package peer

import (
	"google.golang.org/grpc"
)

type PeerAttribute map[string]string

type PeerState int

// this piece of code is copied from google.golang.org/grpc/connectivity.go
// NoConnection state is added in case the connection is not created (connection is nil)
const (
	// Idle indicates the ClientConn is idle.
	Idle PeerState = iota
	// Connecting indicates the ClientConn is connecting.
	Connecting
	// Ready indicates the ClientConn is ready for work.
	Ready
	// TransientFailure indicates the ClientConn has seen a failure but expects to recover.
	TransientFailure
	// Shutdown indicates the ClientConn has started shutting down.
	Shutdown
	// NoConnection indicates the ClientConn has not created.
	NoConnection
)

// String returns the string representation of the PeerState
func (s PeerState) String() string {
	switch s {
	case Idle:
		return "IDLE"
	case Connecting:
		return "CONNECTING"
	case Ready:
		return "READY"
	case TransientFailure:
		return "TRANSIENT_FAILURE"
	case Shutdown:
		return "SHUTDOWN"
	case NoConnection:
		return "NO_CONNECTION"
	default:
		return "INVALID_STATE"
	}
}

// PeerStateFromString returns the PeerState from the given string
func PeerStateFromString(s string) PeerState {
	switch s {
	case "IDLE":
		return Idle
	case "CONNECTING":
		return Connecting
	case "READY":
		return Ready
	case "TRANSIENT_FAILURE":
		return TransientFailure
	case "SHUTDOWN":
		return Shutdown
	case "NO_CONNECTION":
		return NoConnection
	default:
		return NoConnection
	}
}

// PeerInfo contains the address and attributes of a peer
type PeerInfo struct {
	Addr       string
	Attributes map[string]string
}

// Peer contains the peer information and the connection
type Peer struct {
	*PeerInfo
	conn *grpc.ClientConn
}

// NewPeer creates a new peer with the given address and attributes
func NewPeer(addr string, attrs map[string]string) *Peer {
	return &Peer{
		PeerInfo: &PeerInfo{
			Addr:       addr,
			Attributes: attrs,
		},
		conn: nil,
	}
}

// Addr returns the address of the peer
func (p *Peer) Addr() string {
	return p.PeerInfo.Addr
}

// Attributes returns the attributes of the peer
func (p *Peer) Attributes() map[string]string {
	return p.PeerInfo.Attributes
}

func (p *Peer) HasAttributes(opts PeerAttribute) bool {
	attrs := p.Attributes()
	for k, v := range opts {
		if attrs[k] != v {
			return false
		}
	}
	return true
}

// SetConnection sets the connection of the peer
func (p *Peer) SetConnection(conn *grpc.ClientConn) {
	p.conn = conn
}

// GetConnection returns the connection of the peer
func (p *Peer) GetState() PeerState {
	if p.conn == nil {
		return NoConnection
	}
	return PeerStateFromString(p.conn.GetState().String())
}
