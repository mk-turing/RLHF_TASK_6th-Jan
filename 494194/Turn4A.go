package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// WsMessageHandler interface for handling different types of messages
type WsMessageHandler interface {
	OnOpen(client *Client)
	OnMessage(client *Client, message []byte)
	OnClose(client *Client)
}

// Client struct represents a WebSocket client
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	state  map[string]interface{}
	server *WsServer
}

// WsServer struct manages the WebSocket connections and dynamic handler registration
type WsServer struct {
	clients    map[*Client]bool
	handlers   map[string]WsMessageHandler
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.Mutex
	upgrader   websocket.Upgrader
}

// NewWsServer creates a new WebSocket server
func NewWsServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		handlers:   make(map[string]WsMessageHandler),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// RegisterHandler registers a handler dynamically
func (s *WsServer) RegisterHandler(messageType string, handler WsMessageHandler) {
	s.handlers[messageType] = handler
}

// Run the server
func (s *WsServer) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()
			for _, handler := range s.handlers {
				handler.OnOpen(client)
			}
		case client := <-s.unregister:
			s.mu.Lock()
			delete(s.clients, client)
			close(client.send)
			for _, handler := range s.handlers {
				handler.OnClose(client)
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

// ServeHTTP is the WebSocket handler
func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte), server: s, state: make(map[string]interface{})}
	s.register <- client

	go client.readPump(s.handlers)
	go client.writePump()
}

// Client readPump to read messages from WebSocket
func (c *Client) readPump(handlers map[string]WsMessageHandler) {
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

		// Convert messageType to string for map key
		messageTypeStr := fmt.Sprintf("%d", messageType)

		handler, exists := handlers[messageTypeStr]
		if exists {
			handler.OnMessage(c, message)
		} else {
			log.Printf("No handler for message type: %d", messageType)
		}
	}
}

// Client writePump to send messages to WebSocket
func (c *Client) writePump() {
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("writePump error:", err)
			break
		}
	}
}

// TextMessageHandler handles text messages
type TextMessageHandler struct{}

func (h *TextMessageHandler) OnOpen(client *Client) {
	log.Println("TextMessageHandler: WebSocket connection opened for client:", client)
}

func (h *TextMessageHandler) OnMessage(client *Client, message []byte) {
	log.Printf("TextMessageHandler: Received message: %s", string(message))
	client.send <- message
}

func (h *TextMessageHandler) OnClose(client *Client) {
	log.Println("TextMessageHandler: WebSocket connection closed for client:", client)
}

// BinaryMessageHandler handles binary messages
type BinaryMessageHandler struct{}

func (h *BinaryMessageHandler) OnOpen(client *Client) {
	log.Println("BinaryMessageHandler: WebSocket connection opened for client:", client)
}

func (h *BinaryMessageHandler) OnMessage(client *Client, message []byte) {
	log.Printf("BinaryMessageHandler: Received binary message: %x", message)
}

func (h *BinaryMessageHandler) OnClose(client *Client) {
	log.Println("BinaryMessageHandler: WebSocket connection closed for client:", client)
}

// LoadHandlerFromConfig loads handlers dynamically based on the configuration
func loadHandlersFromConfig(configPath string, server *WsServer) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	var config struct {
		Handlers []struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			MessageType string `json:"message_type"`
		} `json:"handlers"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	for _, handlerConfig := range config.Handlers {
		// Corrected reflection logic for handler instantiation
		var handler WsMessageHandler
		switch handlerConfig.Type {
		case "TextMessageHandler":
			handler = &TextMessageHandler{}
		case "BinaryMessageHandler":
			handler = &BinaryMessageHandler{}
		default:
			return fmt.Errorf("unknown handler type: %s", handlerConfig.Type)
		}
		server.RegisterHandler(handlerConfig.MessageType, handler)
	}

	return nil
}

func main() {
	server := NewWsServer()

	// Load handlers dynamically
	if err := loadHandlersFromConfig("handlers_config.json", server); err != nil {
		log.Fatalf("Error loading handlers: %v", err)
	}

	go server.run()

	http.Handle("/ws", server)
	log.Println("Starting WebSocket server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
