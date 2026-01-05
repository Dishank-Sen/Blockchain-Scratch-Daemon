package quic

import (
	"bufio"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
)

func readUntilDelimiter(r *bufio.Reader, delim []byte) ([]byte, error) {
	var buf []byte
	match := 0

	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		buf = append(buf, b)

		switch b {
		case delim[match]:
			match++
			if match == len(delim) {
				return buf, nil
			}
		case delim[0]:
			match = 1
		default:
			match = 0
		}
	}
}

func errorResponse(err error) *types.Response {
	return &types.Response{
		StatusCode: 0,
		Message:    "error",
		Headers:    map[string]string{},
		Body:       []byte(err.Error()),
	}
}