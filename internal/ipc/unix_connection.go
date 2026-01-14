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
	daemonCtx context.Context
	ctx context.Context
	cancel context.CancelFunc
	server *unixServer
}

func newUnixConnection(server *unixServer, conn net.Conn) *unixConnection{
	connCtx, connCancel := context.WithCancel(server.ctx)

	return &unixConnection{
		conn: conn,
		daemonCtx: server.daemonCtx,
		ctx: connCtx,
		cancel: connCancel,
		server: server,
	}
}

func (c *unixConnection) Handle() error{
	defer c.conn.Close()
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

	resp := c.server.dispatch(req)

	return writeResponse(c.conn, resp)
}