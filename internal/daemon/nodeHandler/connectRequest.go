package nodehandler

import (
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

func (h *Handler) ConnectRequest(req *types.Request) *types.Response{
	logger.Info("connect-request handler :)")
	logger.Info(string(req.Body))
	id, err := utils.GetID()
	if err != nil{
		logger.Error(err.Error())
	}
	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: nil,
		Body: fmt.Appendf(nil, "connect request accepted and here is my id : %s", id),
	}
}