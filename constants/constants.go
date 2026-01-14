package constants

import "time"

const(
	SocketPath = "/tmp/blocd.sock"
	WindowsPipeName = `\\.\pipe\blockchain-scratch`
	PublicBootstrapUrl = "44.215.68.195:4242"
	LocalBootstrapUrl = "127.0.0.1:4242"
	DockerBootstrapUrl = "host.docker.internal:4242"
	QuicTimeout = 1*time.Hour
	QuicDialTimeout = 30*time.Second
	QuicStreamTimeout = 30*time.Second
)