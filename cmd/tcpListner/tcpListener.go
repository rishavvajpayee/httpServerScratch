package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rishavvajpayee/httpServerScratch/internal/request"
)

func main() {

	fmt.Println("Rolling the Server")
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error @Listener : ", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error @conn", err)
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("Error @RequestFromReader : ", err)
		}
		fmt.Printf("Request Line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		r.Headers.ForEach(func(n, v string) {
			fmt.Printf("- %s: %s\n", n, v)
		})
		fmt.Printf("Body:\n")
		fmt.Printf("%s\n", r.Body)
	}
}
