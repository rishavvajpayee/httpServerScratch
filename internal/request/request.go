package request

import (
	"bytes"
	"fmt"
	"io"

	header "github.com/rishavvajpayee/httpServerScratch/internal/headers"
)

// file level Errors and initializers
var SEPARATOR = []byte("\r\n")
var ErrBadStartLine = fmt.Errorf("bad start line")
var ErrMalformedRequestLine = fmt.Errorf("malformed request line")
var ErrorRequestState = fmt.Errorf("error in request State")

type ParseState string

const (
	StateInit    ParseState = "init"
	StateDone    ParseState = "done"
	StateHeaders ParseState = "headers"
	StateBody    ParseState = "body"
	StateError   ParseState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     *header.Headers
	Body        []byte
	State       ParseState
}

func (r *Request) hasBody() bool {
	length := r.Headers.GetInt("content-length", 0)
	return length > 0
}

func NewRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: header.NewHeaders(),
	}
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break outer
		}
		switch r.State {

		case StateError:
			return 0, ErrorRequestState

		case StateInit:
			rl, n, err := ParseRequestline(currentData)
			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.State = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n

			// really i dont like this ..
			// in the real world we would not get an EOF after reading the data
			if done {
				if r.hasBody() {
					r.State = StateBody
				} else {
					r.State = StateDone
				}
			}

		case StateBody:
			length := r.Headers.GetInt("content-length", 0)
			if length == 0 {
				panic("Chunk not implemented")
			}
			// push the body to Request.Body
			remaining := min(length-len(r.Body), len(currentData))
			r.Body = append(r.Body, currentData[:remaining]...)
			read += remaining

			if len(r.Body) == length {
				r.State = StateDone
			}

		case StateDone:
			break outer

		default:
			panic("somehow we have programmed poorly")
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

	// this is a stateMachine and changes state as it goes parsing
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
