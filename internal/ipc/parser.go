package ipc

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
)

type Parser struct{
	reader *bufio.Reader
}

func NewParser(conn net.Conn) *Parser{
	return &Parser{
		reader: bufio.NewReader(conn),
	}
}

func (p *Parser) ParseRequest() (*types.Request, error) {
	// Read headers (up to \r\n\r\n)
	rawHeaders, err := readUntilDelimiter(p.reader, []byte("\r\n\r\n"))
	if err != nil {
		return nil, err
	}

	// Parse header section
	lines := strings.Split(string(rawHeaders), "\r\n")

	// Request line: METHOD PATH VERSION
	parts := strings.Split(lines[0], " ")
	if len(parts) < 2 {
		return nil, errors.New("invalid request line")
	}

	req := &types.Request{
		Method:  parts[0],
		Path:    parts[1],
		Headers: make(map[string]string),
	}

	// Parse headers
	for _, line := range lines[1:] {
		if line == "" {
			break
		}
		kv := strings.SplitN(line, ":", 2)
		if len(kv) == 2 {
			req.Headers[strings.TrimSpace(kv[0])] =
				strings.TrimSpace(kv[1])
		}
	}

	// Read body (if Content-Length exists)
	if cl, ok := req.Headers["Content-Length"]; ok {
		n, err := strconv.Atoi(cl)
		if err != nil {
			return nil, err
		}

		req.Body = make([]byte, n)
		if _, err := io.ReadFull(p.reader, req.Body); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func readUntilDelimiter(r *bufio.Reader, delim []byte) ([]byte, error) {
	var buf []byte
	match := 0

	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		buf = append(buf, b)

		if b == delim[match] {
			match++
			if match == len(delim) {
				return buf, nil
			}
		} else {
			if b == delim[0] {
				match = 1
			} else {
				match = 0
			}
		}
	}
}