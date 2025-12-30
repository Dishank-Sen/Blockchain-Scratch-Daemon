package ipc

import (
	"context"
	"net"
	"os"
)

type Server struct {
	listener   net.Listener
	socketPath string
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewServer(ctx context.Context, socketPath string) (*Server, error) {
	ipcCtx, ipcCancel := context.WithCancel(ctx)

	if err := os.RemoveAll(socketPath); err != nil {
		ipcCancel()
		return nil, err
	}
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		ipcCancel()
		return nil, err
	}

	server := &Server{
		listener:   listener,
		socketPath: socketPath,
		ctx:        ipcCtx,
		cancel:     ipcCancel,
	}

	return server, nil
}