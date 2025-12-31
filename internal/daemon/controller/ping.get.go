package controller

import (
	"context"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
)

func PingController(ctx context.Context, req *types.Request) (*types.Response, error){
	header := map[string]string{
		"test": "test-header",
	}
	res := &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: header,
		Body: []byte("pong"),
	}
	return res, nil
}