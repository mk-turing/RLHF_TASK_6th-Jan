package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Define an interface for handling different message types
type WsMessageHandler interface {
	OnOpen(ws *websocket.Conn)
	OnTextMessage(ws *websocket.Conn, message []byte)
	OnBinaryMessage(ws *websocket.Conn, message []byte)
	OnPing(ws *websocket.Conn, message []byte)
	OnPong(ws *websocket.Conn, message []byte)
	OnClose(ws *websocket.Conn)
}

type WsServer struct {
	handler  WsMessageHandler
	upgrader websocket.Upgrader
}

func NewWsServer(handler WsMessageHandler) *WsServer {
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
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		switch messageType {
		case websocket.TextMessage:
			s.handler.OnTextMessage(ws, message)
		case websocket.BinaryMessage:
			s.handler.OnBinaryMessage(ws, message)
		case websocket.PingMessage:
			s.handler.OnPing(ws, message)
		case websocket.PongMessage:
			s.handler.OnPong(ws, message)
		default:
			log.Printf("Unsupported message type: %d", messageType)
		}
	}
}

// SimpleWsHandler implements the WsMessageHandler interface
type SimpleWsHandler struct{}

func (h *SimpleWsHandler) OnOpen(ws *websocket.Conn) {
	log.Println("WebSocket connection opened")
}

func (h *SimpleWsHandler) OnTextMessage(ws *websocket.Conn, message []byte) {
	log.Printf("Received text message: %s", string(message))
	err := ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Error writing text message: %v", err)
	}
}

func (h *SimpleWsHandler) OnBinaryMessage(ws *websocket.Conn, message []byte) {
	log.Printf("Received binary message: %x", message)
	// Handle binary message here
}

func (h *SimpleWsHandler) OnPing(ws *websocket.Conn, message []byte) {
	log.Println("Received ping message")
	err := ws.WriteMessage(websocket.PongMessage, nil)
	if err != nil {
		log.Printf("Error responding with pong message: %v", err)
	}
}

func (h *SimpleWsHandler) OnPong(ws *websocket.Conn, message []byte) {
	log.Println("Received pong message")
	// Handle pong message here
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
