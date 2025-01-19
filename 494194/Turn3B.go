package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientState struct {
	// Add any client-specific state here
	// For example, you can store user information or application-specific data
}

type ClientManager struct {
	clients    map[*websocket.Conn]*ClientState
	mutex      sync.Mutex
	broadcast chan []byte
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*websocket.Conn]*ClientState),
		broadcast: make(chan []byte),
	}
}

func (cm *ClientManager) AddClient(ws *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.clients[ws] = &ClientState{}
}

func (cm *ClientManager) RemoveClient(ws *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.clients, ws)
}

func (cm *ClientManager) Broadcast(message []byte) {
	cm.broadcast <- message
}

func (cm *ClientManager) Run() {
	for {
		message := <-cm.broadcast
		cm.mutex.Lock()
		for ws := range cm.clients {
			if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error writing message to client: %v", err)
				delete(cm.clients, ws)
			}
		}
		cm.mutex.Unlock()
	}
}

// Define an interface for handling different message types
type WsMessageHandler interface {
	OnOpen(ws *websocket.Conn, cm *ClientManager)
	OnTextMessage(ws *websocket.Conn, cm *ClientManager, message []byte)
	OnBinaryMessage(ws *websocket.Conn, cm *ClientManager, message []byte)
	OnPing(ws *websocket.Conn, cm *ClientManager, message []byte)
	OnPong(ws *websocket.Conn, cm *ClientManager, message []byte)
	OnClose(ws *websocket.Conn, cm *ClientManager)
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

	cm := NewClientManager()
	cm.AddClient(ws)
	defer cm.RemoveClient(ws)

	s.handler.OnOpen(ws, cm)
	defer s.handler.OnClose(ws, cm)

	go cm.Run()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		switch messageType {
		case websocket.TextMessage:
			s.handler.OnTextMessage(ws, cm, message)
		case websocket.BinaryMessage:
			s.handler.OnBinaryMessage(ws, cm, message)