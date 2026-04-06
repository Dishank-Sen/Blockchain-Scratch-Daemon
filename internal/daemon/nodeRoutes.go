package daemon

import nodehandler "github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/daemon/nodeHandler"

func (d *Daemon) handleNodeRoutes(){
	n := d.node
	handler := nodehandler.NewNodeHandler(n, d.peerStore)
	n.Handle("ping", handler.Ping)
	n.Handle("accept-peers", handler.AcceptPeers)  // from bootstrap
	n.Handle("connect-request", handler.ConnectRequest)
	n.Handle("accept-request", handler.AcceptRequest)
	n.Handle("message", handler.Message)
}