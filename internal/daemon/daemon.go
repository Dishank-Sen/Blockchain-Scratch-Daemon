package daemon

import (
	"context"
	"errors"
	"fmt"
	"net"

	customerrors "github.com/Dishank-Sen/Blockchain-Scratch-Daemon/customErrors"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/ipc"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/pkg/peerstore"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/node"
	"golang.org/x/sync/errgroup"
)

type Daemon struct{
	cfg  node.Config
	node *node.Node
	server ipc.Server
	peerStore *peerstore.PeerStore
	ctx context.Context
	cancel context.CancelFunc
	group *errgroup.Group
	gctx context.Context
	addr string
}

func NewDaemon(ctx context.Context, addr string) (*Daemon, error) {
	if ctx == nil{
		return nil, fmt.Errorf("context is nil")
	}

	if err := validateAddr(addr); err != nil{
		return nil, err
	}
	
	g, gctx := errgroup.WithContext(ctx)
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

	server, err := ipc.NewServer(daemonCtx, g)
	if err != nil{
		// logger.Debug("ipc error")
		daemonCancel()
		return nil, err
	}
	
	daemon := &Daemon{
		node: n,
		server: server,
		peerStore: peerstore.NewPeerStore(),
		ctx: daemonCtx,
		cancel: daemonCancel,
		group: g,
		gctx: gctx,
		addr: addr,
	}

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
	d.group.Go(func() error{
		<-d.ctx.Done()
		return d.node.Stop()
	})

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

	d.group.Go(func() error{
		return d.node.Wait()
	})
	d.group.Go(func() error{
		d.handleNodeRoutes()
		return nil
	})
	d.group.Go(func() error{
		d.handleIpcRoutes()
		return nil
	})

	// blocks here
	if err := d.server.Listen(); err != nil{
		logger.Debug("error in listening")
		if errors.Is(err, customerrors.ErrServerShutdown){
			logger.Info("server stopped listening")
		}
		d.cancel()
		return err
	}
	return nil
}