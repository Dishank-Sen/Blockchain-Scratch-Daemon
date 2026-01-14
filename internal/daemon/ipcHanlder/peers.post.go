package ipchanlder

import (
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/constants"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

func (h *Handler) PeersController(req *types.Request) *types.Response{
	n := h.node
	resp, err := n.Dial(constants.PublicBootstrapUrl, "peers", req.Headers, req.Body)
	logger.Debug(string(resp.Body))
	if err != nil{
		return &types.Response{
			StatusCode: 500,
			Message: "Error",
			Headers: nil,
			Body: []byte("Internal Server Error"),
		}
	}
	return &types.Response{
		StatusCode: resp.StatusCode,
		Message: resp.Message,
		Headers: resp.Headers,
		Body: resp.Body,
	}
}