package nodehandler

import (
	"context"
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

func (h *Handler) AcceptRequest(ctx context.Context, req *types.Request) *types.Response{
	logger.Info("accept request handler :)")
	logger.Info(string(req.Body))
	id, err := utils.GetID()
	if err != nil{
		logger.Error(err.Error())
	}
	// go func (addr string){
	// 	logger.Info("initializing heartbeat after accept request")
	// 	if err := h.initHeartbeat(addr); err != nil{
	// 		logger.Error(fmt.Sprintf("error in heartbeat: %v", err))
	// 	}
	// }(req.SourceAddr.String())

	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: nil,
		Body: fmt.Appendf(nil, "accept request response from peer id : %s", id),
	}
}