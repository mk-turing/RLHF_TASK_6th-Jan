package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsHandler interface {
	OnOpen(ws *websocket.Conn)
	OnMessage(ws *websocket.Conn, message []byte)
	OnClose(ws *websocket.Conn)
}

type WsServer struct {
	handler  WsHandler
	upgrader websocket.Upgrader
}

func NewWsServer(handler WsHandler) *WsServer {
	return &WsServer{
		handler: handler,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for simplicity
			},
		},
	}
}

func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer ws.Close()

	s.handler.OnOpen(ws)
	defer s.handler.OnClose(ws)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		s.handler.OnMessage(ws, message)
	}
}

type SimpleWsHandler struct{}

func (h *SimpleWsHandler) OnOpen(ws *websocket.Conn) {
	log.Println("WebSocket connection opened")
}

func (h *SimpleWsHandler) OnMessage(ws *websocket.Conn, message []byte) {
	log.Printf("Received message: %s", string(message))
	err := ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Error writing message: %v", err)
	}
}

func (h *SimpleWsHandler) OnClose(ws *websocket.Conn) {
	log.Println("WebSocket connection closed")
}

func main() {
	handler := &SimpleWsHandler{}
	server := NewWsServer(handler)

	http.Handle("/ws", server)

	log.Println("Starting WebSocket server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
