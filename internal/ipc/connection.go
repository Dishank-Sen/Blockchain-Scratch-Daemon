package ipc

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type Connection struct{
	conn net.Conn
	ctx context.Context
	cancel context.CancelFunc
	server *Server
}

func NewConnection(ctx context.Context, server *Server, conn net.Conn) *Connection{
	connCtx, connCancel := context.WithCancel(ctx)

	return &Connection{
		conn: conn,
		ctx: connCtx,
		cancel: connCancel,
		server: server,
	}
}

func (c *Connection) Handle() error{
	parser := NewParser(c.conn)
	req, err := parser.ParseRequest()
	if err != nil{
		if err == io.EOF {
			return nil // normal close
		}
		return err
	}

	logger.Debug(req.Method)
	logger.Debug(req.Path)
	logger.Debug("headers")
	for k, p := range req.Headers{
		logger.Debug(fmt.Sprintf("%s: %s", k, p))
	}
	logger.Debug(string(req.Body))

	resp := c.server.dispatch(c.ctx, req)
	return writeResponse(c.conn, resp)
}

func writeResponse(conn net.Conn, resp *types.Response) error {
	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}

	// Content-Length is mandatory
	resp.Headers["Content-Length"] = fmt.Sprintf("%d", len(resp.Body))

	// 1. Status line
	if _, err := fmt.Fprintf(
		conn,
		"%d %s\r\n",
		resp.StatusCode,
		resp.Message,
	); err != nil {
		return err
	}

	// 2. Headers
	for k, v := range resp.Headers {
		if _, err := fmt.Fprintf(conn, "%s: %s\r\n", k, v); err != nil {
			return err
		}
	}

	// 3. Header/body delimiter
	if _, err := fmt.Fprint(conn, "\r\n"); err != nil {
		return err
	}

	// 4. Body
	_, err := conn.Write(resp.Body)
	return err
}