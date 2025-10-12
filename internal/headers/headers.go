package headers

import (
	"bytes"
	"fmt"
)

var ErrMalformedFieldLine = fmt.Errorf("malformed fieldLine")
var ErrMalformedFieldName = fmt.Errorf("malformed fieldName")
var ErrParseHeader = fmt.Errorf("error parsing headers")

type Headers map[string]string

var rn = []byte("\r\n")

func NewHeaders() Headers {
	return map[string]string{}
}

func ParseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ErrMalformedFieldLine
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	fmt.Printf("Name:%s\nValue:%s\n", name, value)

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", ErrMalformedFieldName
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rn)
		fmt.Println("IDX :", idx)
		if idx == -1 {
			break
		}

		// empty header
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}
		name, value, err := ParseHeader(data[read : read+idx])
		if err != nil {
			return read, done, ErrParseHeader
		}

		read += idx
		h[name] = value
	}
	return read, done, nil
}
