package controller

import (
	"context"
	"encoding/json"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/quic"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type registerBody struct {
	ID string `json:"id"`
}

func RegisterController(ctx context.Context, req *types.Request) (*types.Response, error) {
	var body registerBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return nil, err
	}

	// resp, err := quic.Post(
	// 	ctx,
	// 	"/register",
	// 	req.Headers,
	// 	req.Body,
	// )
	// logger.Debug("register.post - 28")
	// logger.Debug(resp.Message)
	// if err != nil {
	// 	return resp, err
	// }

	// return resp, nil
	go func ()  {
		<-ctx.Done()
		logger.Debug("register.post.go - 37 - context cancelled")
	}()
	return quic.Post(
		"/register",
		req.Headers,
		req.Body,
	)
}
