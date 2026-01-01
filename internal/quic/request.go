package quic

import (
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/quic-go/quic-go"
)

func writeRequest(stream *quic.Stream, req *types.Request) error{
	// Request line
	if _, err := fmt.Fprintf(stream, "%s %s\r\n", req.Method, req.Path); err != nil{
		return err
	}

	// Body framing
	if len(req.Body) > 0 {
		if _, err := fmt.Fprintf(stream, "Content-Length: %d\r\n", len(req.Body)); err != nil {
			return err
		}
	}

	// headers
	for k, v := range req.Headers {
		if _, err := fmt.Fprintf(stream, "%s: %s\r\n", k, v); err != nil {
			return err
		}
	}

	// header/body delimiter
	if _, err := fmt.Fprint(stream, "\r\n"); err != nil {
		return err
	}

	// Body
	if len(req.Body) > 0 {
		if _, err := stream.Write(req.Body); err != nil {
			return err
		}
	}
	return nil
}