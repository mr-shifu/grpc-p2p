package peer

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type PeerAttribute map[string]string

type Peer struct {
	Addr        string
	Attributes  map[string]string
	conn        *grpc.ClientConn
	Name        string
	ClusterName string
}

func NewPeer(addr, name, clusterName string, attrs map[string]string) *Peer {
	return &Peer{
		Addr:        addr,
		Attributes:  attrs,
		Name:        name,
		ClusterName: clusterName,
	}
}

func (p *Peer) State() connectivity.State {
	if p.conn == nil {
		return connectivity.Idle
	}
	return p.conn.GetState()
}
