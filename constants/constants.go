package constants

import "time"

const(
	SocketPath = "/tmp/blocd.sock"
	WindowsPipeName = `\\.\pipe\blockchain-scratch`
	PublicBootstrapUrl = "100.48.90.87:4242"
	LocalBootstrapUrl = "127.0.0.1:4242"
	QuicTimeout = 1*time.Hour
	QuicDialTimeout = 5*time.Second
	QuicStreamTimeout = 5*time.Second
)