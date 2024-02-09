package rpc

import (
	"context"
	"errors"

	"github.com/mr-shifu/grpc-p2p/peer"
	p2p_pb "github.com/mr-shifu/grpc-p2p/proto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	gpeer "google.golang.org/grpc/peer"
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
	peer, err := getPeerFromContext(ctx)
	if err != nil {
		// return nil, errors.New("failed to validate peer")
	}
	r.ps.AddPeer(peer)

	peers := r.ps.GetPeers()
	pbPeers := peersToPbPeers(peers)
	return &p2p_pb.GetPeersResponse{
		Peers: pbPeers,
	}, nil
}

func getPeerFromContext(ctx context.Context) (*peer.Peer, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	cn := md.Get("cluster_name")
	pn := md.Get("peer_name")
	p, ok := gpeer.FromContext(ctx)
	if !ok {
		return nil, errors.New("peer not found")
	}
	return &peer.Peer{
		Name:        pn[0],
		ClusterName: cn[0],
		Addr:        p.LocalAddr.String(),
	}, nil
}

func peersToPbPeers(peers []*peer.Peer) []*p2p_pb.Peer {
	var pbPeers []*p2p_pb.Peer
	for _, peer := range peers {
		p := &p2p_pb.Peer{
			Name:        peer.Name,
			Address:     peer.Addr,
			CLusterName: peer.ClusterName,
			State:       peer.State().String(),
		}
		pbPeers = append(pbPeers, p)
	}
	return pbPeers
}
