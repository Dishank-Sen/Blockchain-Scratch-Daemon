package daemon

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/constants"
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
	if ctx == nil{
		return nil, fmt.Errorf("context is nil")
	}

	if err := validateAddr(addr); err != nil{
		return nil, err
	}
	
	daemonCtx, daemonCancel := context.WithCancel(ctx)

	cfg, err := getConfig(addr)
	if err != nil{
		daemonCancel()
		return nil, err
	}
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

	go func(ctx context.Context) {
		<- ctx.Done()
		logger.Info("daemon context cancelled")
		daemon.stop()
	}(daemonCtx)
	return daemon, nil
}

func validateAddr(addr string) error{
	if addr == ""{
		return fmt.Errorf("addr is empty")
	}
	_, port, err := net.SplitHostPort(addr)
	if err != nil{
		return err
	}
	if port == ""{
		return fmt.Errorf("no port specified")
	}
	return nil
}

func (d *Daemon) Run() error{
	if d.node == nil{
		return fmt.Errorf("node is not defined, it is nil")
	}
	// start node
	if err := d.node.Start(); err != nil{
		logger.Error(fmt.Sprintf("error while starting node: %v", err))
		d.node.Stop()
		d.cancel()
		return err
	}

	logger.Info("node started")

	go d.handleNodeRoutes()
	go func (addr string){
		if err := d.initHeartbeat(addr); err != nil{
			logger.Error(fmt.Sprintf("error in heartbeat: %v", err))
		}
	}(constants.PublicBootstrapUrl)

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

func (d *Daemon) stop(){
	if err := d.node.Stop(); err != nil{
		logger.Error(err.Error())
	}

	d.cancel()
}