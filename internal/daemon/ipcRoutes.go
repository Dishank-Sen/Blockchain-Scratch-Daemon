package daemon

import (
	ipchanlder "github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/daemon/ipcHanlder"
)

func (d *Daemon) handleIpcRoutes(){
	s := d.server
	h := ipchanlder.NewIpcHandler(d.node, d.peerStore)
	s.Post("/connect", h.ConnectController)
	s.Post("/peers", h.PeersController)
	s.Post("/send", h.SendController)
}