package discovery

import (
	"context"
	"sync"
	"time"

	"github.com/mr-shifu/grpc-p2p/peer"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/connectivity"
)

type Discovery struct {
	ps     *peer.PeerService
	logger zerolog.Logger
}

// NewDiscovery creates a new discovery service
func NewDiscovery(ps *peer.PeerService, logger zerolog.Logger) *Discovery {
	return &Discovery{
		ps:     ps,
		logger: logger,
	}
}

// scan starts discovery
// 1. Get adjacent peers from peers in the peerstore
// 2. Remove duplicate peers
func (d *Discovery) scan(ctx context.Context) []*peer.Peer {
	peers := d.ps.GetPeers()

	var allpeers []*peer.Peer
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

func (d *Discovery) addPeers(peers []*peer.Peer) {
	var list []*peer.Peer
	for _, p := range peers {
		if p.Addr() != d.ps.Self().Addr() {
			list = append(list, p)
		}
	}
	added, _ := d.ps.AddPeers(list)
	for _, p := range added {
		d.logger.Info().Str("peer", p.Addr()).Msg("added peer")
	}
}

// refresh verifies peers in the peerstore and connects to the peers if not connected
func (d *Discovery) refresh(ctx context.Context) error {
	peers := d.ps.GetPeers()

	var wg sync.WaitGroup

	for _, p := range peers {
		if p.Addr() == d.ps.Self().Addr() {
			continue
		}
		state, err := d.ps.GetState(p.Addr())
		if err != nil {
			continue
		}
		if state != peer.Ready {
			wg.Add(1)
			go func(p *peer.Peer) {
				defer wg.Done()

				// connect to the peer and add to peerstore
				con, err := d.ps.Connect(p.Addr())

				if err != nil {
					d.logger.Error().Err(err).Str("peer", p.Addr()).Msg("connection failed")
					return
				}

				timedout := time.After(5 * time.Second)
				for con.GetState() != connectivity.Ready {
					select {
					case <-timedout:
						d.logger.Error().Str("peer", p.Addr()).Msg("connection timed out")
						return
					default:
						continue
					}
				}

				d.logger.Info().Str("peer", p.Addr()).Msg("connected")
			}(p)
		}
	}
	wg.Wait()
	return nil
}

// Start starts peer discovery
// 1. Scans all peers in the peerstore to get adjacent peers
// 2. Adds all adjacent peers to the peerstore
// 3. Refreshes peers' connections
func (d *Discovery) Start(ctx context.Context) error {
	go func() {
		for {
			// scan all peers in the peerstore to get adjacent peers
			peers := d.scan(ctx)

			// update peersetore with adjacent peers
			d.addPeers(peers)

			// refresh peers' connections
			if err := d.refresh(ctx); err != nil {
				d.logger.Error().Err(err).Msg("failed to refresh peers")
			}
			// ToDo - make this configurable
			// sleep for 1 second
			time.Sleep(1 * time.Second)
		}
	}()

	// receives context done signal
	<-ctx.Done()

	return nil
}

func removeDuplicatePeers(peers []*peer.Peer) []*peer.Peer {
	encountered := map[string]bool{}
	result := []*peer.Peer{}

	for _, p := range peers {
		if !encountered[p.Addr()] {
			encountered[p.Addr()] = true
			result = append(result, p)
		}
	}
	return result
}
