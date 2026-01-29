package nodehandler

import "github.com/Dishank-Sen/quicnode/node"

type Handler struct{
	node *node.Node
}

func NewNodeHandler(node *node.Node) *Handler{
	return &Handler{
		node: node,
	}
}