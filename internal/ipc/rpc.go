package ipc

import (
	"context"
	"net"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
)

type Rpc struct{
	conn net.Conn
	request *types.Request
	ctx context.Context
	cancel context.CancelFunc
}

func NewRpc(ctx context.Context, conn net.Conn, req *types.Request) *Rpc{
	rpcCtx, rpcCancel := context.WithCancel(ctx)
	
	return &Rpc{
		conn: conn,
		request: req,
		ctx: rpcCtx,
		cancel: rpcCancel,
	}
}

func (r *Rpc) get() error{
	
	return nil
}

func (r *Rpc) post() error{
	return nil
}