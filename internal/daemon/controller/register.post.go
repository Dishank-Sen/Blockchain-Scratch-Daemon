package controller

import (
	"context"
	"encoding/json"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/quic"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
)

type registerBody struct {
	ID string `json:"id"`
}

func RegisterController(ctx context.Context, req *types.Request) (*types.Response, error) {
	var body registerBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return nil, err
	}

	resp, err := quic.Post(
		ctx,
		"/register",
		req.Headers,
		req.Body,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
