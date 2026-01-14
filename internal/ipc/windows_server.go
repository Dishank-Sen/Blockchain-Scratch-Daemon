//go:build windows

package ipc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/Microsoft/go-winio"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/constants"
	customerrors "github.com/Dishank-Sen/Blockchain-Scratch-Daemon/customErrors"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type windowsServer struct {
	listener net.Listener
	daemonCtx context.Context
	ctx      context.Context
	cancel   context.CancelFunc
	routes   map[routeKey]HandlerFunc
	mu       sync.RWMutex
}

func newWindowsServer(ctx context.Context) (*windowsServer, error) {
	ipcCtx, ipcCancel := context.WithCancel(ctx)

	ln, err := winio.ListenPipe(
		constants.WindowsPipeName,
		&winio.PipeConfig{
			InputBufferSize:  64 * 1024,
			OutputBufferSize: 64 * 1024,
			MessageMode:      false, // stream semantics like unix socket
		},
	)
	if err != nil {
		ipcCancel()
		return nil, err
	}

	return &windowsServer{
		listener: ln,
		daemonCtx: ctx,
		ctx:      ipcCtx,
		cancel:   ipcCancel,
		routes:   make(map[routeKey]HandlerFunc),
	}, nil
}

func (s *windowsServer) Listen() error {
	go func() {
		<-s.ctx.Done()
		logger.Warn("windows server closing")
		s.listener.Close()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				logger.Info("windows listener shutting down")
				return customerrors.ErrServerShutdown
			}
			return fmt.Errorf("pipe accept failed: %w", err)
		}

		logger.Debug("new windows pipe connection")

		go func(c net.Conn) {
			conn := newWindowsConnection(s, c)
			if err := conn.Handle(); err != nil {
				logger.Warn(fmt.Sprintf("pipe conn error: %v", err))
			}
		}(conn)
	}
}

func (s *windowsServer) Get(endpoint string, h HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes[routeKey{method: "GET", path: endpoint}] = h
}

func (s *windowsServer) Post(endpoint string, h HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes[routeKey{method: "POST", path: endpoint}] = h
}

func (s *windowsServer) dispatch(req *types.Request) *types.Response {
	s.mu.RLock()
	h, ok := s.routes[routeKey{method: req.Method, path: req.Path}]
	s.mu.RUnlock()

	if !ok {
		return &types.Response{
			StatusCode: 404,
			Message:    "Not Found",
			Body:       []byte("route not found"),
		}
	}

	return h(req)
}
