package types

type Request struct{
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

type Response struct {
	StatusCode int
	Message    string
	Headers    map[string]string
	Body       []byte
}

type StreamMessage struct{
	Version uint16 `json:"version"`
	Header  map[string]string `json:"header"`
    Endpoint    string `json:"endpoint"`
    Length  uint32 `json:"length"`
    Body []byte `json:"body"`
}