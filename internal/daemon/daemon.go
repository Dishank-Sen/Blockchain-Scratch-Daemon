package daemon

import (
	"context"
	"errors"
	"fmt"

	customerrors "github.com/Dishank-Sen/Blockchain-Scratch-Daemon/customErrors"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/ipc"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/node"
)

type Daemon struct{
	cfg  node.Config
	node *node.Node
	server ipc.Server
	ctx context.Context
	cancel context.CancelFunc
	addr string
}

func NewDaemon(ctx context.Context, addr string) (*Daemon, error) {
	daemonCtx, daemonCancel := context.WithCancel(ctx)

	cfg := getConfig(addr)
	n, err := node.NewNode(daemonCtx, cfg)
	if err != nil{
		logger.Debug("error in getting new node")
		daemonCancel()
		return nil, err
	}

	server, err := ipc.NewServer(daemonCtx)
	if err != nil{
		// logger.Debug("ipc error")
		daemonCancel()
		return nil, err
	}
	
	daemon := &Daemon{
		node: n,
		server: server,
		ctx: daemonCtx,
		cancel: daemonCancel,
		addr: addr,
	}

	return daemon, nil
}

func (d *Daemon) Run() error{
	go func() {
		<-d.ctx.Done()
		logger.Info("daemon shutting down..")
	}()

	// if err := quic.InitQuicService(d.ctx); err != nil{
	// 	logger.Error(fmt.Sprintf("error in initializing quic service: %v", err))
	// 	d.cancel()
	// 	return err
	// }

	// start node
	if err := d.node.Start(); err != nil{
		logger.Error(fmt.Sprintf("error while starting node: %v", err))
		d.node.Stop()
		d.cancel()
		return err
	}

	logger.Info("node started")

	go d.handleNodeRoutes()

	// listens for the socket connection requests
	server := d.server

	go d.handleIpcRoutes()

	// blocks here
	if err := server.Listen(); err != nil{
		logger.Debug("error in listening")
		if errors.Is(err, customerrors.ErrServerShutdown){
			logger.Info("server stopped listening")
		}
		d.cancel()
		return err
	}
	return nil
}