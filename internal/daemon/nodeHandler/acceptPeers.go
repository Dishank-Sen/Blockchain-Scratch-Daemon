package nodehandler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

type payload struct{
	ID string `json:"id"`
	Addr string `json:"addr"`
}

func (h *Handler) AcceptPeers(req *types.Request) *types.Response{
	var p payload
	if err := json.Unmarshal(req.Body, &p); err != nil{
		return errorRes()
	}
	id, err := utils.GetID()
	if err != nil{
		logger.Error(err.Error())
	}

	time.Sleep(3*time.Second)

	for range 5{
		logger.Debug(fmt.Sprintf("sending accept request to %s", p.ID))
		resp, err := h.node.Dial(p.Addr, "accept-request", nil, fmt.Appendf(nil, "accept request from %s", id))
		if err != nil{
			logger.Error(err.Error())
			continue
			// return &types.Response{
			// 	StatusCode: 500,
			// 	Message: "Error",
			// 	Headers: nil,
			// 	Body: []byte("not able to connect to another peer"),
			// }
		}
		logger.Info(string(resp.Body))
		return &types.Response{
			StatusCode: 200,
			Message: "ok",
			Headers: nil,
			Body: []byte("accept request successfull"),
		}
	}
	return &types.Response{
		StatusCode: 500,
		Message: "Error",
		Headers: nil,
		Body: []byte("accept request not successfull"),
	}
}

func errorRes() *types.Response{
	errRes := &types.Response{
		StatusCode: 500,
		Message: "Error",
		Headers: map[string]string{},
		Body: []byte("Internal Server Error"),
	}
	return errRes
}