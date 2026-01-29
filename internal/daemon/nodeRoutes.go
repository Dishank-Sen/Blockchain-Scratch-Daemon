package daemon

import nodehandler "github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/daemon/nodeHandler"

func (d *Daemon) handleNodeRoutes(){
	n := d.node
	handler := nodehandler.NewNodeHandler(n)
	n.Handle("ping", handler.Ping)
	n.Handle("peer-info", handler.PeerInfo)
}