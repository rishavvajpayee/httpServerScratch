package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/rishavvajpayee/httpServerScratch/internal/request"
	"github.com/rishavvajpayee/httpServerScratch/internal/response"
)

var RunServerErr = fmt.Errorf("Error in RunServer")

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct {
	closed  bool
	handler Handler
}

func handleConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)

	writer := bytes.NewBuffer([]byte{})
	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, headers)
		return
	}
	handlerErr := s.handler(writer, r)
	var body []byte = nil
	var status response.StatusCode = response.StatusOK
	if handlerErr != nil {
		status = handlerErr.StatusCode
		body = []byte(handlerErr.Message)
	} else {
		body = writer.Bytes()
	}
	headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
	response.WriteStatusLine(conn, status)
	response.WriteHeaders(conn, headers)
	conn.Write(body)
}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}
		go handleConnection(s, conn)
	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("%s : %s", RunServerErr, err)
	}
	server := &Server{
		closed:  false,
		handler: handler,
	}
	go runServer(server, listener)
	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
