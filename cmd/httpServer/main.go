package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rishavvajpayee/httpServerScratch/internal/request"
	"github.com/rishavvajpayee/httpServerScratch/internal/response"
	server "github.com/rishavvajpayee/httpServerScratch/internal/server"
)

const defaultport = 8000

func main() {
	portFlag := flag.Uint("p", defaultport, "an int")
	flag.Parse()
	port := uint16(*portFlag)
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		handlerErr := &server.HandlerError{}
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			handlerErr.StatusCode = response.StatusBadRequest
			handlerErr.Message = "Your problem is not my problem\n"
			return handlerErr
		case "/myproblem":
			handlerErr.StatusCode = response.StatusInternalServerError
			handlerErr.Message = "Woopsie, my bad\n"
			return handlerErr
		default:
			w.Write([]byte("All good, frfr\n"))
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
