package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
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
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())
	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF { // client disconnected
				break
			}
			fmt.Println("read error:", err) // other read error
			continue
		}
		msg := buf[:n]
		s.broadcast(msg)
	}
}

/* Broadcast message to all connected clients */
func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error:", err)
			}
		}(ws)
	}
}

func main() {
	signupEndpoint()
	structureEndpoint()
}

func signupEndpoint() {
	server := NewServer()
	http.Handle("/account", websocket.Handler(server.handleWS))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting signup server:", err)
	}
}

func structureEndpoint() {
	server := NewServer()
	http.Handle("/chats", websocket.Handler(server.handleWS))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting structure server:", err)
	}
}
