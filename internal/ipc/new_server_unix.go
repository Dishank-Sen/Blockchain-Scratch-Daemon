//go:build !windows

package ipc

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func NewServer(ctx context.Context, g *errgroup.Group) (Server, error) {
	return newUnixServer(ctx, g)
}
