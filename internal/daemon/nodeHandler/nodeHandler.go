package nodehandler

import (
	"fmt"
	"time"

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

func (h *Handler) initHeartbeat(addr string) error{
	resp, err := h.node.Dial(addr, "heartbeat", nil, nil)
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