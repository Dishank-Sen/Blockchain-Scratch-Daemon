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

type messagePayload struct {
	From    string    `json:"from"`
	To      string    `json:"to"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sent_at"`
}

func (h *Handler) Message(ctx context.Context, req *types.Request) *types.Response{
	var msg messagePayload
	if err := json.Unmarshal(req.Body, &msg); err != nil {
		logger.Error(fmt.Sprintf("failed to parse message: %v", err))
		return &types.Response{
			StatusCode: 400,
			Message: "Bad Request",
			Body: []byte("invalid message format"),
		}
	}
	
	myID, err := utils.GetID()
	if err != nil {
		logger.Error(err.Error())
		return &types.Response{
			StatusCode: 500,
			Message: "Error",
			Body: []byte("failed to get node ID"),
		}
	}
	
	logger.Info(fmt.Sprintf("📨 Message from %s: %s", msg.From, msg.Content))
	
	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Body: fmt.Appendf(nil, "message received by %s", myID),
	}
}
