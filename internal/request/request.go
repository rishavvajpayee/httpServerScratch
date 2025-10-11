package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var SEPARATOR = "\r\n"
var ERROR_BAD_START_LINE = fmt.Errorf("bad start line")

type ParseState string

const (
	StateInit ParseState = "init"
	StateDone ParseState = "done"
)

func (r *Request) parse(data []byte) (int, error) {
	return 0, nil
}

func ParseRequestline(s string) (*RequestLine, string, error) {
	idx := strings.Index(s, SEPARATOR)
	if idx == -1 {
		return nil, s, nil
	}
	startLine := s[:idx]
	resOfMsg := s[idx+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, resOfMsg, ERROR_BAD_START_LINE
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, resOfMsg, ERROR_BAD_START_LINE
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return rl, resOfMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to io.ReadAll"), err)
	}
	strData := string(data)
	rl, _, err := ParseRequestline(strData)
	return &Request{
		RequestLine: *rl,
	}, err
}
