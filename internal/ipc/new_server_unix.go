//go:build !windows
package ipc

import "context"

func NewServer(ctx context.Context) (Server, error) {
	return newUnixServer(ctx)
}
