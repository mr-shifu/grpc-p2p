package rpc

import (
	"context"
	"errors"
	"strings"

	"github.com/mr-shifu/grpc-p2p/peer"
	p2p_pb "github.com/mr-shifu/grpc-p2p/proto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
		return nil, errors.New("failed to validate peer")
	}
	defer r.ps.AddPeer(peer)
	
	peers := r.ps.GetPeers()
	pbPeers := peersToPbPeers(peers)
	return &p2p_pb.GetPeersResponse{
		Peers: pbPeers,
	}, nil
}

func getPeerFromContext(ctx context.Context) (*peer.Peer, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	addrs := md.Get("addr")
	if len(addrs) == 0 {
		return nil, errors.New("peer address not found")
	}

	attrs := make(map[string]string)
	for k, v := range md {
		if strings.HasPrefix(k, "attr-") {
			key := strings.TrimPrefix(k, "attr-")
			attrs[key] = v[0]
		}
	}

	p := peer.NewPeer(addrs[0], attrs)

	return p, nil
}

func peersToPbPeers(peers []*peer.Peer) []*p2p_pb.Peer {
	var pbPeers []*p2p_pb.Peer
	for _, peer := range peers {
		var attrs []*p2p_pb.Attribute
		for k, v := range peer.Attributes() {
			attrs = append(attrs, &p2p_pb.Attribute{
				Key:   k,
				Value: v,
			})
		}
		p := &p2p_pb.Peer{
			Address:    peer.Addr(),
			Attributes: attrs,
			State:      peer.GetState().String(),
		}
		pbPeers = append(pbPeers, p)
	}
	return pbPeers
}
