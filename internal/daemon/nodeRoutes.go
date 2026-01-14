package daemon

import nodehandler "github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/daemon/nodeHandler"

func (d *Daemon) handleNodeRoutes(){
	n := d.node

	n.Handle("ping", nodehandler.Ping)
}