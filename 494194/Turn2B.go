package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// MessageType defines the different types of WebSocket messages.
type MessageType int

const (
	TextMessage   MessageType = iota
	BinaryMessage
	PingMessage
	PongMessage
	CloseMessage
)

// Message represents a WebSocket message.
type Message struct {
	Type    MessageType
	Payload []byte
}

// WsMessageHandler defines the interface for handling different message types.
type WsMessageHandler interface {
	HandleTextMessage(ws *websocket.Conn, message []byte)
	HandleBinaryMessage(ws *websocket.Conn, message []byte)
	HandlePingMessage(ws *websocket.Conn, message []byte)
	HandlePongMessage(ws *websocket.Conn, message []byte)
	HandleCloseMessage(ws *websocket.Conn, message []byte)
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

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			s.handleClose(ws, messageType, message)
			break
		}

		message := Message{
			Type:    MessageType(messageType),
			Payload: message,
		}
		s.handleMessage(ws, message)
	}
}

func (s *WsServer) handleMessage(ws *websocket.Conn, message Message) {
	switch message.Type {
	case TextMessage:
		s.handler.HandleTextMessage(ws, message.Payload)
	case BinaryMessage:
		s.handler.HandleBinaryMessage(ws, message.Payload)
	case PingMessage:
		s.handler.HandlePingMessage(ws, message.Payload)
	case PongMessage:
		s.handler.HandlePongMessage(ws, message.Payload)
	case CloseMessage:
		s.handleClose(ws, message.Type, message.Payload)
	default:
		log.Printf("Unsupported message type: %v", message.Type)
	}
}

func (s *WsServer) handleClose(ws *websocket.Conn, messageType int, message []byte) {
	// Handle close message here (e.g., call OnClose handler)
	s.handler.HandleCloseMessage(ws, message)
}

// Example usage:
type SimpleWsHandler struct{}

func (h *SimpleWsHandler) HandleTextMessage(ws *websocket.Conn, message []byte) {
	fmt.Printf("Received Text Message: %s\n", string(message))
	// Send a response back to the client
	err := ws.WriteMessage(websocket.TextMessage, []byte("Received your message!"))
	if err != nil {
		log.Printf("Error writing text message: %v", err)
	}
}

func (h *SimpleWsHandler) HandleBinaryMessage(ws *websocket.Conn, message []byte) {
	fmt.Printf("Received Binary Message: %v\n", message)
}

func (h *SimpleWsHandler) HandlePingMessage(ws *websocket.Conn, message []byte) {