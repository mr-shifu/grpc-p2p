package peer

import (
	"context"
	"time"

	p2p_pb "github.com/mr-shifu/grpc-p2p/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type Discovery struct {
	store *PeerStore
}

func NewDiscovery(store *PeerStore) *Discovery {
	return &Discovery{
		store: store,
	}
}

// getPeers gets adjacent peers from a peer calling GetPeers RPC
func (d *Discovery) getPeers(ctx context.Context, peer *Peer) ([]*Peer, error) {
	if peer.State() != connectivity.Ready {
		conn, err := grpc.Dial(peer.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		peer.conn = conn
	}

	client := p2p_pb.NewPeerServiceClient(peer.conn)
	neighbors, err := client.GetPeers(ctx, &p2p_pb.GetPeersRequest{})
	if err != nil {
		return nil, err
	}

	np := peersFromPbPeers(neighbors.Peers)

	return np, nil
}

// scan starts discovery
// 1. Get adjacent peers from peers in the peerstore
// 2. Adds all adjacent peers to the peerstore
func (d *Discovery) scan(ctx context.Context) error {
	peers := d.store.GetAllPeers()

	var allpeers []*Peer
	for _, peer := range peers {
		neighbors, err := d.getPeers(ctx, peer)
		if err != nil {
			continue
		}
		allpeers = append(allpeers, neighbors...)
	}

	if err := d.store.AddPeers(allpeers...); err != nil {
		return err
	}

	return nil
}

// refresh verifies peers in the peerstore and connects to the peers if not connected
func (d *Discovery) refresh(ctx context.Context) error {
	peers := d.store.GetAllPeers()

	for _, peer := range peers {
		if peer.State() != connectivity.Ready {
			// connect to the peer
			conn, err := grpc.Dial(peer.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				continue
			}
			peer.conn = conn
		}
		// add peer to peerstore if not exists or update peer connetion
		if err := d.store.AddPeer(peer); err != nil {
			continue
		}
	}
	return nil
}

func (d *Discovery) Start(ctx context.Context) error {
	go func() {
		for {
			d.scan(ctx)
			d.refresh(ctx)
			time.Sleep(1 * time.Second)
		}
	}()

	// receives context done signal
	<-ctx.Done()

	return nil
}
