package daemon

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type heartbeatPayload struct {
	ID string `json:"id"`
}

func (d *Daemon) initHeartbeat(addr string) error{
	id, err := utils.GetID()
	if err != nil{
		return err
	}

	hb := &heartbeatPayload{
		ID: id,
	}

	byteData, err := json.Marshal(hb)
	if err != nil{
		return err
	}

	resp, err := d.node.Dial(addr, "heartbeat", nil, byteData)
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