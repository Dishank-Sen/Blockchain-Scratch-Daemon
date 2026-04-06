package ipchanlder

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type sendRequest struct {
	To      string `json:"to"`       // peer ID
	Message string `json:"message"`  // message content
}

type messagePayload struct {
	From    string    `json:"from"`
	To      string    `json:"to"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sent_at"`
}

func (h *Handler) SendController(req *types.Request) *types.Response {
	var sendReq sendRequest
	if err := json.Unmarshal(req.Body, &sendReq); err != nil {
		logger.Error(fmt.Sprintf("failed to unmarshal send request: %v", err))
		return &types.Response{
			StatusCode: 400,
			Message:    "Bad Request",
			Headers:    nil,
			Body:       []byte("invalid JSON payload"),
		}
	}

	// Validate inputs
	if sendReq.To == "" {
		return &types.Response{
			StatusCode: 400,
			Message:    "Bad Request",
			Headers:    nil,
			Body:       []byte("'to' field is required"),
		}
	}

	if sendReq.Message == "" {
		return &types.Response{
			StatusCode: 400,
			Message:    "Bad Request",
			Headers:    nil,
			Body:       []byte("'message' field is required"),
		}
	}

	// Get peer from store
	peer, exists := h.peerStore.Get(sendReq.To)
	if !exists {
		logger.Error(fmt.Sprintf("peer %s not found in store", sendReq.To))
		return &types.Response{
			StatusCode: 404,
			Message:    "Not Found",
			Headers:    nil,
			Body:       []byte(fmt.Sprintf("peer %s not connected", sendReq.To)),
		}
	}

	// Get own ID
	myID, err := utils.GetID()
	if err != nil {
		logger.Error(err.Error())
		return &types.Response{
			StatusCode: 500,
			Message:    "Error",
			Headers:    nil,
			Body:       []byte("failed to get node ID"),
		}
	}

	// Create message payload
	msg := messagePayload{
		From:    myID,
		To:      sendReq.To,
		Content: sendReq.Message,
		SentAt:  time.Now(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to marshal message: %v", err))
		return &types.Response{
			StatusCode: 500,
			Message:    "Error",
			Headers:    nil,
			Body:       []byte("failed to create message"),
		}
	}

	// Send message to peer
	logger.Debug(fmt.Sprintf("sending message to %s at %s", peer.ID, peer.Addr))
	resp, err := h.node.Dial(peer.Addr, "message", nil, msgBytes)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to send message: %v", err))
		return &types.Response{
			StatusCode: 500,
			Message:    "Error",
			Headers:    nil,
			Body:       []byte(fmt.Sprintf("failed to send message: %v", err)),
		}
	}

	if resp.StatusCode != 200 {
		logger.Error(fmt.Sprintf("peer returned error: %s", string(resp.Body)))
		return &types.Response{
			StatusCode: 500,
			Message:    "Error",
			Headers:    nil,
			Body:       []byte(fmt.Sprintf("peer error: %s", string(resp.Body))),
		}
	}

	logger.Info(fmt.Sprintf("✅ Message sent to %s", peer.ID))
	return &types.Response{
		StatusCode: 200,
		Message:    "ok",
		Headers:    nil,
		Body:       []byte(fmt.Sprintf("message sent to %s", peer.ID)),
	}
}
