package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client: ", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)
		// EOF means that the connection from the client is closed, so need to check on that
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("read error: ", err)
			continue
		}

		// since n has its own size, we should assigned buffer only up to n's size
		// to prevent from sending the whole buffer
		msg := buf[:n]
		fmt.Println(string(msg))
		ws.Write([]byte("thank you for the msg!!"))
	}
}

func main() {
	// initialise server
	s := NewServer()
	http.Handle("/ws", websocket.Handler(s.handleWS))
	http.ListenAndServe(":3000", nil)
}