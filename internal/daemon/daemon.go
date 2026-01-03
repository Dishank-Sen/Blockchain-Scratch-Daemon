package daemon

import (
	"context"
	"errors"

	customerrors "github.com/Dishank-Sen/Blockchain-Scratch-Daemon/customErrors"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/daemon/controller"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/ipc"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type Daemon struct{
	server ipc.Server
	ctx context.Context
	cancel context.CancelFunc
}

func NewDaemon(ctx context.Context) (*Daemon, error) {
	daemonCtx, daemonCancel := context.WithCancel(ctx)

	server, err := ipc.NewServer(daemonCtx)
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
		logger.Info("daemon shutting down..")
	}()

	// listens for the socket connection requests
	server := d.server

	server.Get("/ping", controller.PingController)
	server.Post("/register", controller.RegisterController)
	server.Get("/peers", controller.PeersController)

	// blocks here
	if err := server.Listen(); err != nil{
		if errors.Is(err, customerrors.ErrServerShutdown){
			logger.Info("server stopped listening")
		}
		d.cancel()
		return err
	}
	return nil
}