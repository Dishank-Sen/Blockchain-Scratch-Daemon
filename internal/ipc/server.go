package ipc

import (
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
)

type Server interface{
	Listen() error
	Get(endpoint string, h HandlerFunc)
	Post(endpoint string, h HandlerFunc)
	dispatch(req *types.Request) *types.Response
}