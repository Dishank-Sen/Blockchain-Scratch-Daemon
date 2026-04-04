//go:build !windows

package ipc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	// "strings"
	"sync"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/constants"
	customerrors "github.com/Dishank-Sen/Blockchain-Scratch-Daemon/customErrors"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"golang.org/x/sync/errgroup"
)

type unixServer struct {
	listener   net.Listener
	daemonCtx  context.Context
	ctx        context.Context
	cancel     context.CancelFunc
	group     *errgroup.Group
	routes     map[routeKey]HandlerFunc
	mu         sync.RWMutex
}

func newUnixServer(ctx context.Context, g *errgroup.Group) (*unixServer, error) {
	ipcCtx, ipcCancel := context.WithCancel(ctx)

	if err := os.RemoveAll(constants.SocketPath); err != nil {
		ipcCancel()
		return nil, err
	}
	listener, err := net.Listen("unix", constants.SocketPath)
	if err != nil {
		ipcCancel()
		return nil, err
	}

	server := &unixServer{
		listener:   listener,
		daemonCtx:  ctx,
		ctx:        ipcCtx,
		cancel:     ipcCancel,
		group:      g,
		routes:     make(map[routeKey]HandlerFunc),
	}

	return server, nil
}

func (s *unixServer) Listen() error{
	s.group.Go(func() error{
		<-s.ctx.Done()
		logger.Warn("ipc server closing")
		return s.listener.Close()
	})

	for{
		conn, err := s.listener.Accept() // blocking
		if err != nil{
			if s.ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				logger.Info("listener shutting down")
				return customerrors.ErrServerShutdown
			}

			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				logger.Warn(fmt.Sprintf("temporary accept error: %v", err))
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return fmt.Errorf("listener accept failed: %v", err)
		}
		logger.Debug("new connection")

		s.group.Go(func() error{
			conn := newUnixConnection(s, conn)
			// if err := conn.Handle(); err != nil{
			// 	// logger.Warn(fmt.Sprintf("conn error: %v", err))
			// 	// if isTimeoutError(err){
			// 	// 	logger.Error("bootstrap connection timed out")
			// 	// 	go s.cancel()
			// 	// }
			// 	// go s.cancel()
			// 	return fmt.Errorf(fmt.Sprintf("ipc connection error: %v", err))
			// }
			return conn.Handle()
		})
		// go func (c net.Conn)  {
		// 	conn := newUnixConnection(s, c)
		// 	if err := conn.Handle(); err != nil{
		// 		logger.Warn(fmt.Sprintf("conn error: %v", err))
		// 		// if isTimeoutError(err){
		// 		// 	logger.Error("bootstrap connection timed out")
		// 		// 	go s.cancel()
		// 		// }
		// 		go s.cancel()
		// 	}
		// }(conn)
	}
}


func (s *unixServer) Get(endpoint string, h HandlerFunc){
	s.mu.Lock()
	s.routes[routeKey{method: "GET", path: endpoint}] = h
	s.mu.Unlock()
}

func (s *unixServer) Post(endpoint string, h HandlerFunc){
	s.mu.Lock()
	s.routes[routeKey{method: "POST", path: endpoint}] = h
	s.mu.Unlock()
}

func (s *unixServer) dispatch(req *types.Request) *types.Response {
	// logger.Debug("server.go - 96")
	// logger.Debug(req.Path)
	key := routeKey{
		method: req.Method,
		path:   req.Path,
	}

	s.mu.RLock()
	h, ok := s.routes[key]
	s.mu.RUnlock()

	if !ok {
		// logger.Debug("server.go - 135")
		return &types.Response{
			StatusCode: 404,
			Message:    "Not Found",
			Body:       []byte("route not found"),
		}
	}

	return h(req)
}

