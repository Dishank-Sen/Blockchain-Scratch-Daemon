package nodehandler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

type connectRequestPayload struct {
	ID string `json:"id"`
}

func (h *Handler) ConnectRequest(ctx context.Context, req *types.Request) *types.Response{
	logger.Info("connect-request handler :)")
	
	var payload connectRequestPayload
	if err := json.Unmarshal(req.Body, &payload); err != nil {
		logger.Error(fmt.Sprintf("failed to parse connect request: %v", err))
		return &types.Response{
			StatusCode: 400,
			Message: "Bad Request",
			Body: []byte("invalid payload"),
		}
	}
	
	logger.Info(fmt.Sprintf("connect request from %s", payload.ID))
	
	// Store the peer
	h.peerStore.Add(payload.ID, req.SourceAddr.String())
	logger.Debug(fmt.Sprintf("peer %s stored at %s", payload.ID, req.SourceAddr.String()))
	
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