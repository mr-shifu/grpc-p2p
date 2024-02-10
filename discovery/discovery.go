package discovery

import (
	"context"
	"time"

	"github.com/mr-shifu/grpc-p2p/peer"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/connectivity"
)

type Discovery struct {
	ps     *peer.PeerService
	logger zerolog.Logger
}

func NewDiscovery(ps *peer.PeerService, logger zerolog.Logger) *Discovery {
	return &Discovery{
		ps:     ps,
		logger: logger,
	}
}

// scan starts discovery
// 1. Get adjacent peers from peers in the peerstore
// 2. Adds all adjacent peers to the peerstore
func (d *Discovery) scan(ctx context.Context) []peer.Peer {
	peers := d.ps.GetPeers()

	var allpeers []peer.Peer
	for _, peer := range peers {
		neighbors, err := d.ps.GetNeighbors(ctx, peer)
		if err != nil {
			continue
		}
		allpeers = append(allpeers, neighbors...)
	}

	allpeers = removeDuplicatePeers(allpeers)

	return allpeers
}

// refresh verifies peers in the peerstore and connects to the peers if not connected
func (d *Discovery) refresh(ctx context.Context, peers []peer.Peer) error {
	for _, peer := range peers {
		if peer.State() != connectivity.Ready {
			// connect to the peer and add to peerstore
			if _, err := d.ps.Connect(&peer); err != nil {
				continue
			}
		}
	}
	return nil
}

func (d *Discovery) Start(ctx context.Context) error {
	go func() {
		for {
			peers := d.scan(ctx)
			if err := d.refresh(ctx, peers); err != nil {
				d.logger.Error().Err(err).Msg("failed to refresh peers")
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// receives context done signal
	<-ctx.Done()

	return nil
}

func removeDuplicatePeers(peers []peer.Peer) []peer.Peer {
	encountered := map[string]bool{}
	result := []peer.Peer{}

	for _, p := range peers {
		if !encountered[p.Addr] {
			encountered[p.Addr] = true
			result = append(result, p)
		}
	}
	return result
}
