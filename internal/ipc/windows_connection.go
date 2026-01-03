//go:build windows

package ipc

import (
	"context"
	"io"
	"net"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type windowsConnection struct {
	conn   net.Conn
	ctx    context.Context
	cancel context.CancelFunc
	server *windowsServer
}

func newWindowsConnection(
	ctx context.Context,
	server *windowsServer,
	conn net.Conn,
) *windowsConnection {
	connCtx, connCancel := context.WithCancel(ctx)

	return &windowsConnection{
		conn:   conn,
		ctx:    connCtx,
		cancel: connCancel,
		server: server,
	}
}

func (c *windowsConnection) Handle() error {
	defer c.cancel()

	go func() {
		<-c.ctx.Done()
		logger.Warn("closing windows pipe connection")
		c.conn.Close()
	}()

	parser := NewParser(c.conn)

	req, err := parser.ParseRequest()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	resp, err := c.server.dispatch(c.ctx, req)
	if err != nil {
		if rerr := writeResponse(c.conn, resp); rerr != nil {
			return rerr
		}
		return err
	}

	return writeResponse(c.conn, resp)
}
