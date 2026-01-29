package nodehandler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/types"
)

type payload struct{
	ID string `json:"id"`
	Addr string `json:"addr"`
}

func (h *Handler) PeerInfo(req *types.Request) *types.Response{
	var pl payload
	if err := json.Unmarshal(req.Body, &pl); err != nil{
		return errorRes()
	}
	logger.Info(pl.Addr)
	logger.Info(pl.ID)
	time.Sleep(3*time.Second)

	for range 5{
		logger.Debug("dialing...")
		resp, err := h.node.Dial(pl.Addr, "peer-side", nil, fmt.Appendf(nil, "dialed %s", pl.ID))
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
		break
	}
	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: nil,
		Body: []byte("got response from another peer"),
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