package daemon

import (
	"fmt"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

func (d *Daemon) initHeartbeat(addr string) error{
	resp, err := d.node.Dial(addr, "heartbeat", nil, nil)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200{
		return fmt.Errorf("not healthy for %s", addr)
	}
	logger.Info("healthy")
	time.Sleep(5*time.Second)
	d.initHeartbeat(addr)
	return nil
}