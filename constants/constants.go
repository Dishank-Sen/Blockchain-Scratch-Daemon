package constants

import "time"

const(
	SocketPath = "/tmp/blocd.sock"
	WindowsPipeName = `\\.\pipe\blockchain-scratch`
	PublicBootstrapUrl = "44.204.118.99:4242"
	LocalBootstrapUrl = "127.0.0.1:4242"
	DockerBootstrapUrl = "host.docker.internal:4242"
	MaxIdleTimeout = 30*time.Second
	KeepAlivePeriod = 10*time.Second
	QuicDialTimeout = 30*time.Second
	QuicStreamTimeout = 30*time.Second
)