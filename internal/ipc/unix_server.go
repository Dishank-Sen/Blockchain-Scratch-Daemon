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
)

type unixServer struct {
	listener   net.Listener
	ctx        context.Context
	cancel     context.CancelFunc
	routes     map[routeKey]HandlerFunc
	mu         sync.RWMutex
}

func newUnixServer(ctx context.Context) (*unixServer, error) {
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
		ctx:        ipcCtx,
		cancel:     ipcCancel,
		routes:     make(map[routeKey]HandlerFunc),
	}

	return server, nil
}

func (s *unixServer) Listen() error{
	go func ()  {
		<-s.ctx.Done()
		logger.Warn("server closing..")
		s.listener.Close()
	}()

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

		go func (c net.Conn)  {
			conn := newUnixConnection(s.ctx, s, c)
			if err := conn.Handle(); err != nil{
				logger.Warn(fmt.Sprintf("conn error: %v", err))
				// if isTimeoutError(err){
				// 	logger.Error("bootstrap connection timed out")
				// 	go s.cancel()
				// }
				go s.cancel()
			}
		}(conn)
	}
}

// func isTimeoutError(err error) bool {
// 	if err == nil {
// 		return false
// 	}

// 	// 1. Context deadline
// 	if errors.Is(err, context.DeadlineExceeded) {
// 		return true
// 	}

// 	// 2. net.Error timeout
// 	var netErr net.Error
// 	if errors.As(err, &netErr) && netErr.Timeout() {
// 		return true
// 	}

// 	// 3. quic-go idle timeout (fallback)
// 	if strings.Contains(err.Error(), "no recent network activity") {
// 		return true
// 	}

// 	return false
// }


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

func (s *unixServer) dispatch(ctx context.Context, req *types.Request) (*types.Response, error) {
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
		logger.Debug("server.go - 135")
		return &types.Response{
			StatusCode: 404,
			Message:    "Not Found",
			Body:       []byte("route not found"),
		}, nil
	}

	resp, err := h(ctx, req)  // IMPORTANT LINE
	logger.Debug("uni_server.go - 147")
	logger.Debug(resp.Message)
	logger.Debug(string(resp.Body))

	if err != nil {
		logger.Debug("server.go - 148")
		logger.Error(err.Error())
		return resp, err
	}

	return resp, nil
}

