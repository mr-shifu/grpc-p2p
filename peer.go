package p2p

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

func (p *Peer) State() connectivity.State {
	if p.conn == nil {
		return connectivity.State(-1)
	}
	return p.conn.GetState()
}
