package peer

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Peer struct {
	Name        string
	ClusterName string
	Addr        string
	conn        *grpc.ClientConn
}

func NewPeer(addr, name, clusterName string) *Peer {
	return &Peer{
		Name:        name,
		ClusterName: clusterName,
		Addr:        addr,
	}
}

func (p *Peer) State() connectivity.State {
	if p.conn == nil {
		return connectivity.State(-1)
	}
	return p.conn.GetState()
}
