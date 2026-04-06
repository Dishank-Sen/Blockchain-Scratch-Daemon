package ipchanlder

import (
	"encoding/json"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type peerInfo struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

func (h *Handler) PeersController(req *types.Request) *types.Response{
	// Get all connected peers from peer store
	peers := h.peerStore.GetAll()
	
	// Convert to response format
	peerList := make([]peerInfo, 0, len(peers))
	for _, p := range peers {
		peerList = append(peerList, peerInfo{
			ID:   p.ID,
			Addr: p.Addr,
		})
	}
	
	responseData, err := json.Marshal(peerList)
	if err != nil {
		logger.Error("failed to marshal peer list")
		return &types.Response{
			StatusCode: 500,
			Message: "Error",
			Headers: nil,
			Body: []byte("Internal Server Error"),
		}
	}
	
	logger.Debug(string(responseData))
	
	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: nil,
		Body: responseData,
	}
}