package nodehandler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

func (h *Handler) AcceptRequest(ctx context.Context, req *types.Request) *types.Response{
	logger.Info("accept request handler :)")
	
	var payload acceptRequestPayload
	if err := json.Unmarshal(req.Body, &payload); err != nil {
		logger.Error(fmt.Sprintf("failed to parse accept request: %v", err))
		return errorRes()
	}
	
	logger.Info(fmt.Sprintf("accept request from %s", payload.ID))
	
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
		Body: fmt.Appendf(nil, "accept request response from peer id : %s", id),
	}
}

func errorRes() *types.Response{
	return &types.Response{
		StatusCode: 500,
		Message: "Error",
		Headers: nil,
		Body: []byte("Internal Server Error"),
	}
}