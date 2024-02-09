package rpc

import (
	"context"

	"github.com/mr-shifu/grpc-p2p/peer"
	p2p_pb "github.com/mr-shifu/grpc-p2p/proto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type RpcService struct {
	ps     *peer.PeerService
	logger zerolog.Logger

	p2p_pb.UnimplementedPeerServiceServer
}

func NewRpcService(ps *peer.PeerService, logger zerolog.Logger) *RpcService {
	return &RpcService{
		ps:     ps,
		logger: logger,
	}
}

func (rs *RpcService) RegisterService(s grpc.ServiceRegistrar) {
	p2p_pb.RegisterPeerServiceServer(s, rs)
}

func (r *RpcService) GetPeers(ctx context.Context, req *p2p_pb.GetPeersRequest) (*p2p_pb.GetPeersResponse, error) {
	peers := r.ps.GetPeers()
	pbPeers := peersToPbPeers(peers)
	return &p2p_pb.GetPeersResponse{
		Peers: pbPeers,
	}, nil
}

func peersToPbPeers(peers []*peer.Peer) []*p2p_pb.Peer {
	var pbPeers []*p2p_pb.Peer
	for _, peer := range peers {
		p := &p2p_pb.Peer{
			Name:        peer.Name,
			Address:     peer.Addr,
			CLusterName: peer.ClusterName,
		}
		pbPeers = append(pbPeers, p)
	}
	return pbPeers
}
