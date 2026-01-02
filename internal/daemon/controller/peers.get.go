package controller

import (
	"context"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/quic"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
)

func PeersController(ctx context.Context, req *types.Request) (*types.Response, error){
	return quic.Get(
		ctx,
		"/peers",
		req.Headers,
	)
}