package ipchanlder

import (
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/node"
)

type Handler struct{
	node *node.Node
	ID string   // node id
}

func NewIpcHandler(n *node.Node) *Handler{
	id, err := utils.GetID()
	if err != nil{
		logger.Error(err.Error())
	}
	return &Handler{
		node: n,
		ID: id,
	}
}