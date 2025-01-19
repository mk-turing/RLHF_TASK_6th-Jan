package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Define an interface for handling different message types
type WsMessageHandler interface {
	OnOpen(client *Client)
	OnTextMessage(client *Client, message []byte)
	OnBinaryMessage(client *Client, message []byte)
	OnClose(client *Client)
}

// Client represents a single connection and its state
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	state  map[string]interface{}
	server *WsServer
}

// WsServer manages multiple client connections
type WsServer struct {
	clients    map[*Client]bool
	handler    WsMessageHandler
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.Mutex
	upgrader   websocket.Upgrader
}

// NewWsServer creates a new WebSocket server with the given handler
func NewWsServer(handler WsMessageHandler) *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		handler:    handler,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for simplicity
			},
		},
	}
}

func (s *WsServer) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()
			s.handler.OnOpen(client)
		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
				s.handler.OnClose(client)
			}
			s.mu.Unlock()
		case message := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte), server: s, state: make(map[string]interface{})}
	s.register <- client

	go client.readPump(s.handler)
	go client.writePump()
}

func (c *Client) readPump(handler WsMessageHandler) {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("readPump error:", err)
			break
		}
		switch messageType {
		case websocket.TextMessage:
			handler.OnTextMessage(c, message)
		case websocket.BinaryMessage:
			handler.OnBinaryMessage(c, message)
		default:
			log.Printf("Unsupported message type: %d", messageType)
		}
	}
}

func (c *Client) writePump() {
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("writePump error:", err)
			break
		}
	}
}

type MyWsHandler struct{}

func (h *MyWsHandler) OnOpen(client *Client) {
	log.Println("WebSocket connection opened for client:", client)
	client.state["exampleKey"] = "exampleValue"
}

func (h *MyWsHandler) OnTextMessage(client *Client, message []byte) {
	log.Printf("Client %v sent text message: %s", client, string(message))
	// Echo the message back
	client.send <- message
}

func (h *MyWsHandler) OnBinaryMessage(client *Client, message []byte) {
	log.Printf("Client %v sent binary message: %x", client, message)
}

func (h *MyWsHandler) OnClose(client *Client) {
	log.Println("WebSocket connection closed for client:", client)
}

func main() {
	handler := &MyWsHandler{}
	server := NewWsServer(handler)
	go server.run()

	http.Handle("/ws", server)
	log.Println("Starting WebSocket server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
