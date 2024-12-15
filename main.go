package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

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

func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("New Connection from Clien to Orderbook Feed", ws.RemoteAddr())
	for {
		payload := fmt.Sprintf("Orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New Connection from Client", ws.RemoteAddr())

	s.conns[ws] = true
	s.readConnection(ws)
}

func (s *Server) readConnection(ws *websocket.Conn) {

	buf := make([]byte, 1024)
	for {
		// reads from the connection to the buffer
		// n is the number of bytes read
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection Closed")
				break
			}
			fmt.Println("Error reading from connection", err)
			continue
		}

		// msg is the message received from the client till n bytes
		msg := buf[:n]
		fmt.Println(string(msg))
		// ws.Write([]byte("Thankyou for your message"))
		s.broadcast(msg)
	}
}

func (s *Server) broadcast(msg []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(msg); err != nil {
				fmt.Println("Error broadcasting message", err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/orderbookFeed", websocket.Handler(server.handleWSOrderbook))
	http.ListenAndServe(":3000", nil)
}
