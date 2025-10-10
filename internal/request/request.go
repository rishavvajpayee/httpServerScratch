package request

import "io"

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	// Headers     map[string]string
	// Body        []byte
}

var ERROR_BAD_START_LINE error

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestLine := RequestLine{
		HttpVersion:   "1.1",
		RequestTarget: "/",
		Method:        "GET",
	}
	r := Request{
		RequestLine: requestLine,
	}
	return &r, nil
}
