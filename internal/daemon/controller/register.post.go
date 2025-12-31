package controller

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/quic-go/quic-go"
)

type registerBody struct{
	ID string `json:"id"`
}

type streamMessage struct{
	Version uint16 `json:"version"`
	Header  map[string]string `json:"header"`
    Type    string `json:"type"`
    Length  uint32 `json:"length"`
    Payload []byte `json:"payload"`
}

func RegisterController(ctx context.Context, req *types.Request) (*types.Response, error){
	session, err := quic.DialAddr(ctx, "127.0.0.1:4242", clientTLSConfig(), clientQuicConfig())
	if err != nil{
		return nil, err
	}

	stream, err := session.OpenStreamSync(ctx)
	if err != nil{
		return nil, err
	}
	defer stream.Close()

	var b registerBody
	if err := json.Unmarshal(req.Body, &b); err != nil{
		return nil, err
	}

	logger.Debug(b.ID)

	msg := streamMessage{
		Version: 1,
		Header: req.Headers,
		Type: "register",
		Payload: req.Body,
	}

	msgBytes, _ := json.Marshal(msg)

	// ---- Write framed message ----
	if err := writeFramed(stream, msgBytes); err != nil {
		return nil, err
	}

	logger.Info("register sent")

	res := &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: map[string]string{},
		Body: []byte("peer registered"),
	}

	session.CloseWithError(0, "client shutdown")
	return res, nil
	// // ---- Optional: read response (if server sends any) ----
	// go readLoop(stream)

	// // 🔴 KEEP SESSION ALIVE
	// <-ctx.Done()

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

/* -------------------- HELPERS -------------------- */

func writeFramed(w io.Writer, data []byte) error {
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(data)))

	if _, err := w.Write(lenBuf[:]); err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func readLoop(r io.Reader) {
	for {
		msg, err := readFramed(r)
		if err != nil {
			logger.Error(fmt.Sprintf("read error: %v", err))
			return
		}
		logger.Info(fmt.Sprintf("server response: %s", string(msg)))
	}
}

func readFramed(r io.Reader) ([]byte, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lenBuf[:])
	if length == 0 || length > 1<<20 {
		return nil, fmt.Errorf("invalid length: %d", length)
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}