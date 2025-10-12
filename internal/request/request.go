package request

import (
	"bytes"
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	State       ParseState
}

var SEPARATOR = []byte("\r\n")
var ErrBadStartLine = fmt.Errorf("bad start line")
var ErrMalformedRequestLine = fmt.Errorf("malformed request line")
var ErrorRequestState = fmt.Errorf("error in request State")

type ParseState string

const (
	StateInit  ParseState = "init"
	StateDone  ParseState = "done"
	StateError ParseState = "error"
)

func NewRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case StateError:
			return 0, ErrorRequestState
		case StateInit:
			rl, n, err := ParseRequestline(data[read:])
			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.State = StateDone

		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
}

func ParseRequestline(s []byte) (*RequestLine, int, error) {
	idx := bytes.Index(s, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	startLine := s[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := NewRequest()

	// NOTE: buf could get overrun
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}
		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}
	return request, nil
}
