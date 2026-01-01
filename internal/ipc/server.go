package ipc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type HandlerFunc func(ctx context.Context, req *types.Request) (*types.Response, error)

type routeKey struct {
	method string
	path   string
}

type Server struct {
	listener   net.Listener
	socketPath string
	ctx        context.Context
	cancel     context.CancelFunc
	routes map[routeKey]HandlerFunc
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
		routes:     make(map[routeKey]HandlerFunc),
	}

	return server, nil
}

func (s *Server) Listen() error{
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
				return nil
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
			conn := NewConnection(s.ctx, s, c)
			if err := conn.Handle(); err != nil{
				logger.Warn(fmt.Sprintf("conn error: %v", err))
			}
		}(conn)
	}
}

func (s *Server) Get(endpoint string, h HandlerFunc){
	s.routes[routeKey{method: "GET", path: endpoint}] = h
}

func (s *Server) Post(endpoint string, h HandlerFunc){
	s.routes[routeKey{method: "POST", path: endpoint}] = h
}

func (s *Server) dispatch(ctx context.Context, req *types.Request) *types.Response {
	logger.Debug("server.go - 96")
	logger.Debug(req.Path)
	key := routeKey{
		method: req.Method,
		path:   req.Path,
	}

	h, ok := s.routes[key]
	if !ok {
		logger.Debug("server.go - 105")
		return &types.Response{
			StatusCode: 404,
			Message:    "Not Found",
			Body:       []byte("route not found"),
		}
	}

	resp, err := h(ctx, req)
	logger.Debug(resp.Message)

	if err != nil {
		logger.Debug("server.go - 112")
		logger.Error(err.Error())
		return &types.Response{
			StatusCode: 500,
			Message:    "Internal Error",
			Body:       []byte(err.Error()),
		}
	}

	return resp
}

