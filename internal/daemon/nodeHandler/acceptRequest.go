package nodehandler

import (
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

func (h *Handler) AcceptRequest(req *types.Request) *types.Response{
	logger.Info(string(req.Body))
	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: nil,
		Body: []byte("got accept request"),
	}
}