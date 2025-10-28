package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var ErrMalformedFieldLine = fmt.Errorf("malformed fieldLine")
var ErrMalformedFieldName = fmt.Errorf("malformed fieldName")
var ErrParseHeader = fmt.Errorf("error parsing headers")

func isValidFieldName(chars []byte) bool {
	/*
		In other words, a field-name must contain only:
		Uppercase letters: A-Z
		Lowercase letters: a-z
		Digits: 0-9
		Special characters: !, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~
	*/

	for _, ch := range chars {
		found := false
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}

		if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' {
			found = true
		}

		if !found {
			return false
		}
	}
	return true
}

var rn = []byte("\r\n")

func ParseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ErrMalformedFieldLine
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", ErrMalformedFieldName
	}

	return string(name), string(value), nil
}

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) (string, bool) {
	value, ok := h.headers[strings.ToLower(name)]
	return value, ok
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s,%s", v, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rn)
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

		if !isValidFieldName([]byte(name)) {
			return read, done, ErrMalformedFieldName
		}

		read += idx + len(rn)
		h.Set(name, value)
	}
	return read, done, nil
}
