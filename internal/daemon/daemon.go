package daemon

import (
	"context"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/daemon/controller"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/ipc"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type Daemon struct{
	server *ipc.Server
	ctx context.Context
	cancel context.CancelFunc
}

func NewDaemon(ctx context.Context) (*Daemon, error) {
	daemonCtx, daemonCancel := context.WithCancel(ctx)
	socketPath := "/tmp/blocd.sock"

	server, err := ipc.NewServer(daemonCtx, socketPath)
	if err != nil{
		daemonCancel()
		return nil, err
	}

	daemon := &Daemon{
		server: server,
		ctx: daemonCtx,
		cancel: daemonCancel,
	}

	return daemon, nil
}

func (d *Daemon) Run() error{
	go func() {
		<-d.ctx.Done()
		logger.Warn("daemon shutting down..")
	}()

	server := d.server

	server.Get("/ping", controller.PingController)
	server.Post("/register", controller.RegisterController)


	return server.Listen()  // blocks here
}