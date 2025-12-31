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