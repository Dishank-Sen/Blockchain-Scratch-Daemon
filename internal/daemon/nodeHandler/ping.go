package nodehandler

import "github.com/Dishank-Sen/quicnode/types"

func Ping(req *types.Request) *types.Response{
	return &types.Response{
		StatusCode: 200,
		Message: "ok",
		Headers: map[string]string{"test-header": "test1"},
		Body: []byte("pong"),
	}
}