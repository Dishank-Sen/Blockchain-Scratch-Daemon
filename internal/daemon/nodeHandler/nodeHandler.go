package nodehandler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/node"
)

type Handler struct{
	node *node.Node
}

func NewNodeHandler(node *node.Node) *Handler{
	return &Handler{
		node: node,
	}
}

type heartbeatPayload struct {
	ID string `json:"id"`
}

func (h *Handler) initHeartbeat(addr string) error{
	id, err := utils.GetID()
	if err != nil{
		return err
	}

	hb := &heartbeatPayload{
		ID: id,
	}

	byteData, err := json.Marshal(hb)
	if err != nil{
		return err
	}

	resp, err := h.node.Dial(addr, "heartbeat", nil, byteData)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200{
		return fmt.Errorf("not healthy for %s", addr)
	}
	logger.Info("healthy")
	time.Sleep(5*time.Second)
	h.initHeartbeat(addr)
	return nil
}