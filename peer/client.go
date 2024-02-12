package peer

import (
	"context"

	p2p_pb "github.com/mr-shifu/grpc-p2p/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetPeers(ctx context.Context, cc *grpc.ClientConn, self *PeerInfo) ([]*Peer, error) {
	var opts []string
	opts = append(opts, "addr", self.Addr)
	for k, v := range self.Attributes {
		key := "attr-" + k
		opts = append(opts, key, v)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, opts...)

	client := p2p_pb.NewPeerServiceClient(cc)
	neighbors, err := client.GetPeers(ctx, &p2p_pb.GetPeersRequest{})
	if err != nil {
		return nil, err
	}

	return peersFromPbPeers(neighbors.Peers), nil
}

func peersFromPbPeers(pbPeers []*p2p_pb.Peer) []*Peer {
	var peers []*Peer
	for _, p := range pbPeers {
		attrs := make(map[string]string)
		for _, attr := range p.Attributes {
			attrs[attr.Key] = attr.Value
		}
		p := NewPeer(p.Address, attrs)
		peers = append(peers, p)
	}
	return peers
}
