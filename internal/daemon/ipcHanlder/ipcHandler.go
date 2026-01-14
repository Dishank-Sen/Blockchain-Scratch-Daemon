package ipchanlder

import "github.com/Dishank-Sen/quicnode/node"

type Handler struct{
	node *node.Node
}

func NewIpcHandler(n *node.Node) *Handler{
	return &Handler{
		node: n,
	}
}