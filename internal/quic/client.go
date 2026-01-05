package quic

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/constants"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/quic-go/quic-go"
)

type quicService struct {
	conn *quic.Conn
	ctx context.Context
	cancel context.CancelFunc
	mu   sync.Mutex
}

var (
	quicSvcMu sync.Mutex
	quicSvc   *quicService
)

func InitQuicService(ctx context.Context) error{
	quicSvcMu.Lock()
	defer quicSvcMu.Unlock()

	if quicSvc != nil {
		return nil // already initialized
	}

	clientCtx, clientCancel := context.WithCancel(ctx)

	// IMPORTANT: use context with timeout for dial
	dialCtx, dialCancel := context.WithTimeout(context.Background(), constants.QuicDialTimeout)
	defer dialCancel()
	conn, err := quic.DialAddr(
		dialCtx,
		constants.PublicBootstrapUrl,
		clientTLSConfig(),
		clientQuicConfig(),
	)
	if err != nil {
		clientCancel()
		return err
	}

	logger.Debug("quic connection established")
	
	quicSvc = &quicService{
		conn: conn,
		ctx: clientCtx,
		cancel: clientCancel,
	}
	quicSvc.watchLifecycle()
	return nil
}

func (q *quicService) watchLifecycle() {
    go func() {
        select {
        case <-q.ctx.Done():
            // daemon shutdown
            logger.Warn("daemon shutting down, closing quic")
            q.conn.CloseWithError(0, "daemon shutdown")

        case <-q.conn.Context().Done():
            // QUIC-owned close (idle timeout, peer close, reset)
            logger.Warn("quic connection closed: " + q.conn.Context().Err().Error())
        }

        quicSvcMu.Lock()
        quicSvc = nil
        quicSvcMu.Unlock()
    }()
}


func Get(path string, headers map[string]string) (*types.Response, error) {
	if quicSvc == nil {
		return errorResponse(fmt.Errorf("connection is closed")), fmt.Errorf("connection is closed")
	}
	return quicSvc.get(path, headers)
}

func (q *quicService) get(path string, headers map[string]string) (*types.Response, error) {
	streamCtx, streamCancel := context.WithTimeout(q.ctx, constants.QuicStreamTimeout)
	defer streamCancel()
	stream, err := q.conn.OpenStreamSync(streamCtx)
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
func Post(path string, headers map[string]string, body []byte) (*types.Response, error) {
	if quicSvc == nil {
		return errorResponse(fmt.Errorf("connection is closed")), fmt.Errorf("connection is closed")
	}
	return quicSvc.post(path, headers, body)
}

func (q *quicService) post(path string, headers map[string]string, body []byte) (*types.Response, error) {
	streamCtx, streamCancel := context.WithTimeout(q.ctx, constants.QuicStreamTimeout)
	defer streamCancel()
	stream, err := q.conn.OpenStreamSync(streamCtx)
	if err != nil {
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


func clientTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-example-v1"},
	}
}

func clientQuicConfig() *quic.Config{
	return &quic.Config{
		MaxIdleTimeout: constants.QuicTimeout,
	}
}
