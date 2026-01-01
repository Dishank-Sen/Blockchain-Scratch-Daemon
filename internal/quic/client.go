package quic

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/quic-go/quic-go"
)

type Client struct {
	conn *quic.Conn
	mu   sync.Mutex
}

var (
	client     *Client
	clientOnce sync.Once
)

func getClient(ctx context.Context) (*Client, error) {
	var err error

	clientOnce.Do(func() {
		var conn *quic.Conn
		conn, err = quic.DialAddr(
			ctx,
			"127.0.0.1:4242",
			clientTLSConfig(),
			clientQuicConfig(),
		)
		if err != nil {
			return
		}
		client = &Client{conn: conn}
	})

	return client, err
}

func Get(ctx context.Context, path string, headers map[string]string) (*types.Response, error){
	c, err := getClient(ctx)
	if err != nil{
		return nil, err
	}

	return c.get(ctx, path, headers)
}

func (c *Client) get(ctx context.Context, path string, headers map[string]string) (*types.Response, error){
	c.mu.Lock()
	defer c.mu.Unlock()

	stream, err := c.conn.OpenStreamSync(ctx)
	if err != nil{
		return nil, err
	}

	defer stream.Close()

	req := &types.Request{
		Method: "GET",
		Path: path,
		Headers: headers,
		Body: []byte(""),
	}

	// Write request
	if err := writeRequest(stream, req); err != nil{
		return nil, err
	}

	// Read response
	return readResponse(stream) 
}

// main post function which every code calls
func Post(ctx context.Context, path string, headers map[string]string, body []byte) (*types.Response, error) {
	c, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	return c.post(ctx, path, headers, body)
}

func (c *Client) post(ctx context.Context, path string, headers map[string]string, body []byte) (*types.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	stream, err := c.conn.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	req := &types.Request{
		Method: "POST",
		Path: path,
		Headers: headers,
		Body: body,
	}

	// Write request
	if err := writeRequest(stream, req); err != nil{
		return nil, err
	}

	// Read response
	return readResponse(stream) 
}

func clientTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-example-v1"},
	}
}

func clientQuicConfig() *quic.Config{
	return &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	}
}

