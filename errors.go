package p2p

import "errors"

var (
	ErrServerGracefullyShutdownTimedout = errors.New("server gracefully shutdown timed out")
	ErrServerStartServerFailed          = errors.New("server start failed")

	
)
