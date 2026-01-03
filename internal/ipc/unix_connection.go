//go:build !windows
package ipc

import (
	"context"
	"io"
	"net"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type unixConnection struct{
	conn net.Conn
	ctx context.Context
	cancel context.CancelFunc
	server *unixServer
}

func newUnixConnection(ctx context.Context, server *unixServer, conn net.Conn) *unixConnection{
	connCtx, connCancel := context.WithCancel(ctx)

	return &unixConnection{
		conn: conn,
		ctx: connCtx,
		cancel: connCancel,
		server: server,
	}
}

func (c *unixConnection) Handle() error{
	go func ()  {
		<-c.ctx.Done()
		logger.Warn("closing connection")
		c.conn.Close()
	}()

	parser := NewParser(c.conn)
	req, err := parser.ParseRequest()
	if err != nil{
		if err == io.EOF {
			return nil // normal close
		}
		return err
	}

	resp, err := c.server.dispatch(c.ctx, req)
	if err != nil{
		logger.Debug("connection.go - 49")
		logger.Error(err.Error())
		logger.Debug("writing response - connection.go - 51")
		if rerr := writeResponse(c.conn, resp); rerr != nil{
			return rerr
		}
		return err
	}

	return writeResponse(c.conn, resp)
}