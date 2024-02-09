package p2p

import "errors"

var (
	ErrServerGracefullyShutdownTimedout = errors.New("server gracefully shutdown timed out")
	ErrServerStartServerFailed          = errors.New("server start failed")

	ErrInvalidPeerAddress     = errors.New("invalid peer address")
	ErrInvalidPeerName        = errors.New("invalid peer name")
	ErrInvalidPeerClusterName = errors.New("invalid peer cluster name")
)
