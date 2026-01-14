package ipchanlder

import (
	"encoding/json"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/constants"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

type connectBody struct {
	ID string `json:"id"`
}

type peersList struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

func (h *Handler) ConnectController(req *types.Request) *types.Response {
	var body connectBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return nil
	}

	n := h.node

	// returns peer list in body
	resp, err := n.Dial(constants.PublicBootstrapUrl, "connect", req.Headers, req.Body)

	if err != nil || resp.StatusCode != 200{
		logger.Error(err.Error())
		return &types.Response{
			StatusCode: 500,
			Message: "Error",
			Headers: nil,
			Body: []byte("can't able to register node"),
		}
	}

	if peerExists(resp.Body){
		var p []peersList
		if err := json.Unmarshal(resp.Body, &p); err != nil{
			logger.Error(err.Error())
		}
		if h.connectToPeer(p){
			logger.Debug("connected to some peer")
			return &types.Response{
				StatusCode: 200,
				Message: "ok",
				Headers: nil,
				Body: []byte("connected"),
			}
		}
		return &types.Response{
			StatusCode: 400,
			Message: "Error",
			Headers: nil,
			Body: []byte("can't able to connect"),
		}
	}
	return &types.Response{
		StatusCode: 404,
		Message: "Error",
		Headers: nil,
		Body: []byte("no peer exists"),
	}
}

func peerExists(body []byte) bool{
	var p []peersList
	if err := json.Unmarshal(body, &p); err != nil{
		logger.Error(err.Error())
		return false
	}
	if len(p) == 0{
		logger.Debug("no peers")
		return false
	}
	return true
}

func (h *Handler) connectToPeer(p []peersList) bool{
	for _, peer := range p{
		logger.Debug("dialing...")
		resp, err := h.node.Dial(peer.Addr, "ping", nil, []byte("ping"))
		if err != nil{
			logger.Debug("some error while dialing")
			logger.Error(err.Error())
			continue
		}
		logger.Debug(string(resp.Body))
		return true
	}
	return false
}