package nodehandler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

type payload struct{
	ID string `json:"id"`
	Addr string `json:"addr"`
}

type acceptRequestPayload struct {
	ID string `json:"id"`
}

func (h *Handler) AcceptPeers(ctx context.Context, req *types.Request) *types.Response{
	logger.Info("accept-peer handler :)")
	var p payload
	if err := json.Unmarshal(req.Body, &p); err != nil{
		return errorResponseWithLog("failed to parse accept-peers payload")
	}
	
	// Store the peer from bootstrap info
	h.peerStore.Add(p.ID, p.Addr)
	logger.Debug(fmt.Sprintf("peer %s stored at %s (from bootstrap)", p.ID, p.Addr))
	
	id, err := utils.GetID()
	if err != nil{
		logger.Error(err.Error())
		return errorResponseWithLog("failed to get node ID")
	}

	time.Sleep(1*time.Second)
	baseDelay := time.Second
	
	acceptPayload := acceptRequestPayload{
		ID: id,
	}
	
	payloadBytes, err := json.Marshal(acceptPayload)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to marshal payload: %v", err))
		return errorResponseWithLog("failed to create request payload")
	}

	for i := 0; i < 5; i++{
		logger.Debug(fmt.Sprintf("sending accept request to %s", p.ID))
		resp, err := h.node.Dial(p.Addr, "accept-request", nil, payloadBytes)
		if err != nil{
			logger.Error(err.Error())
			// exponential backoff
            delay := baseDelay * (1 << i) // 1s,2s,4s,8s,16s
            logger.Debug(fmt.Sprintf("retrying in %v", delay))
            time.Sleep(delay)
			continue
		}
		logger.Info(string(resp.Body))
		return &types.Response{
			StatusCode: 200,
			Message: "ok",
			Headers: nil,
			Body: []byte("accept request successfull"),
		}
	}
	return &types.Response{
		StatusCode: 500,
		Message: "Error",
		Headers: nil,
		Body: []byte("accept request not successfull"),
	}
}

func errorResponseWithLog(msg string) *types.Response{
	logger.Error(msg)
	errRes := &types.Response{
		StatusCode: 500,
		Message: "Error",
		Headers: map[string]string{},
		Body: []byte(msg),
	}
	return errRes
}