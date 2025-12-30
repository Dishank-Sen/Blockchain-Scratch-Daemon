package daemon

import (
	"context"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/ipc"
)

type Daemon struct{
	socket *ipc.Server
	ctx context.Context
	cancel context.CancelFunc
}

func NewDaemon(ctx context.Context) (*Daemon, error) {
	daemonCtx, daemonCancel := context.WithCancel(ctx)
	socketPath := "/tmp/blocd.sock"

	socket, err := ipc.NewServer(daemonCtx, socketPath)
	if err != nil{
		daemonCancel()
		return nil, err
	}

	daemon := &Daemon{
		socket: socket,
		ctx: daemonCtx,
		cancel: daemonCancel,
	}

	return daemon, nil
}

func (d *Daemon) Run() error{
	
}