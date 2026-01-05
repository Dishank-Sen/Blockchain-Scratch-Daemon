package ipc

import (
	"context"
	"fmt"
	"net"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
)

type HandlerFunc func(ctx context.Context, req *types.Request) (*types.Response, error)

type routeKey struct {
	method string
	path   string
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