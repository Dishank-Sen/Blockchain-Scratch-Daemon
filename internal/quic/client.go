package quic

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/constants"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/quic-go/quic-go"
)

type Client struct {
	conn *quic.Conn
	mu   sync.Mutex
}

var (
	clientMu sync.Mutex
	client   *Client
)

func getClient(ctx context.Context) (*Client, error) {
	clientMu.Lock()
	defer clientMu.Unlock()

	if client != nil {
		return client, nil
	}

	// IMPORTANT: use Background for dialing
	conn, err := quic.DialAddr(
		context.Background(),
		constants.PublicBootstrapUrl,
		clientTLSConfig(),
		clientQuicConfig(),
	)
	if err != nil {
		logger.Error("quic dial failed: " + err.Error())
		return nil, err
	}

	logger.Debug("quic connection established")

	client = &Client{conn: conn}

	go func() {
		<-ctx.Done()
		logger.Error(ctx.Err().Error())
		logger.Warn("closing quic connection")
		conn.CloseWithError(0, "client exiting")

		clientMu.Lock()
		client = nil // allow reconnect
		clientMu.Unlock()
	}()

	return client, nil
}


func Get(ctx context.Context, path string, headers map[string]string) (*types.Response, error) {
	c, err := getClient(ctx)
	if err != nil {
		return errorResponse(err), err
	}
	return c.get(ctx, path, headers)
}

func (c *Client) get(ctx context.Context, path string, headers map[string]string) (*types.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	stream, err := c.conn.OpenStreamSync(ctx)
	if err != nil {
		return errorResponse(err), err
	}
	defer stream.Close()

	req := &types.Request{
		Method:  "GET",
		Path:    path,
		Headers: headers,
	}

	if err := writeRequest(stream, req); err != nil {
		return errorResponse(err), err
	}

	resp, err := readResponse(stream)
	if err != nil {
		return errorResponse(err), err
	}

	return resp, nil
}


// main post function which every code calls
func Post(ctx context.Context, path string, headers map[string]string, body []byte) (*types.Response, error) {
	c, err := getClient(ctx)
	if err != nil {
		logger.Debug("some error - client.go - 103")
		return errorResponse(err), err
	}
	logger.Debug("client get - client.go - 105")
	return c.post(ctx, path, headers, body)
}

func (c *Client) post(ctx context.Context, path string, headers map[string]string, body []byte) (*types.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	stream, err := c.conn.OpenStreamSync(ctx)
	if err != nil {
		logger.Debug("some error - client - 116")
		logger.Error(err.Error())
		return errorResponse(err), err
	}
	defer stream.Close()

	req := &types.Request{
		Method:  "POST",
		Path:    path,
		Headers: headers,
		Body:    body,
	}

	if err := writeRequest(stream, req); err != nil {
		return errorResponse(err), err
	}

	resp, err := readResponse(stream)
	if err != nil {
		return errorResponse(err), err
	}

	return resp, nil
}


func errorResponse(err error) *types.Response {
	return &types.Response{
		StatusCode: 0,
		Message:    err.Error(),
		Headers:    map[string]string{},
		Body:       nil,
	}
}


func clientTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-example-v1"},
	}
}

func clientQuicConfig() *quic.Config{
	return &quic.Config{
		MaxIdleTimeout: 60 * time.Minute,
		KeepAlivePeriod: 15 * time.Second,
	}
}
