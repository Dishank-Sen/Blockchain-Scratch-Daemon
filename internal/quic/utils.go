package quic

import "bufio"

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