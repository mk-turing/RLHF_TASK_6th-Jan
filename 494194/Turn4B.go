package main

import (
	"fmt"
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
	handlers   []WsMessageHandler
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.Mutex
	upgrader   websocket.Upgrader
}

// HandlerConfig holds the configuration for a WsMessageHandler
type HandlerConfig struct {
	Name   string   `json:"name"`   // Handler name
	Module string   `json:"module"` // Module path (Go package name)
	Funcs  []string `json:"funcs"`  // Functions to register as handlers
}

// NewWsServer creates a new WebSocket server with given handlers
func NewWsServer(config []HandlerConfig) *WsServer {
	s := &WsServer{
		clients:    make(map[*Client]bool),
		handlers:   []WsMessageHandler{},
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for simplicity
			},
		},
	}
	if err := s.LoadHandlers(config); err != nil {
		log.Fatalf("Error loading handlers: %v", err)
	}
	return s
}

// LoadHandlers dynamically loads handlers based on the configuration
func (s *WsServer) LoadHandlers(config []HandlerConfig) error {
	for _, c := range config {
		if err := s.loadHandler(c); err != nil {
			return fmt.Errorf("error loading handler '%s': %w", c.Name, err)
		}
	}
	return nil
}

func (s *WsServer) loadHandler(config HandlerConfig) error {
	module, err := loadModule(config.Module)
	if err != nil {
		return err
	}
	for _, funcName := range config.Funcs {
		h, err := newHandler(module, funcName)
		if err != nil {
			return err
		}
		s.handlers = append(s.handlers, h)
	}
	return nil
}

// loadModule loads a Go module from a given path and returns its exports
func loadModule(path string) (map[string]interface{}, error) {
	// Implement loading logic using packages like "plugin" for Go plugins
	// For demonstration, we'll use a simple map as an example
	exports := map[string]interface{}{
		"MyHandler": &MyHandler{},
	}
	return exports, nil
}

// newHandler creates a new WsMessageHandler instance
