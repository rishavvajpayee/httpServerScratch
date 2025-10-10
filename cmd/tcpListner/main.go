package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func ReadLinesFromReader(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(out)
		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				out <- str
				data = data[i+1:]
				str = ""
			}
			str += string(data)
		}
	}()
	return out
}

func main() {

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error @Listener : ", err)
	}
	for {
		fmt.Println("Rolling the Server")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error @conn : ", err)
		}

		lines := ReadLinesFromReader(conn)
		for line := range lines {
			fmt.Printf("read : %s\n", line)
		}
	}

}
