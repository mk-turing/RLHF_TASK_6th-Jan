package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// MessageHandler defines an interface for handling WebSocket messages.
type MessageHandler interface {
	HandleMessage(conn *websocket.Conn, messageType int, data []byte)
}

// StateHandler defines an interface for handling WebSocket connection state changes.
type StateHandler interface {
	OnOpen(conn *websocket.Conn)
	OnClose(conn *websocket.Conn)
	OnError(conn *websocket.Conn, err error)
}

type EchoMessageHandler struct{}

func (e *EchoMessageHandler) HandleMessage(conn *websocket.Conn, messageType int, data []byte) {
	fmt.Printf("Received message: %s\n", data)
	if err := conn.WriteMessage(messageType, data); err != nil {
		fmt.Println("Error writing message:", err)
	}
}

type LogStateHandler struct{}

func (l *LogStateHandler) OnOpen(conn *websocket.Conn) {
	fmt.Println("Connection opened")
}

func (l *LogStateHandler) OnClose(conn *websocket.Conn) {
	fmt.Println("Connection closed")
}

func (l *LogStateHandler) OnError(conn *websocket.Conn, err error) {
	fmt.Println("Connection error:", err)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request, msgHandler MessageHandler, stateHandler StateHandler) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(w, "Failed to upgrade to websocket: %v\n", err)
		return
	}
	defer conn.Close()

	stateHandler.OnOpen(conn)
	defer stateHandler.OnClose(conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			stateHandler.OnError(conn, err)
			break
		}
		msgHandler.HandleMessage(conn, messageType, message)
	}
}

func main() {
	messageHandler := &EchoMessageHandler{}
	stateHandler := &LogStateHandler{}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, messageHandler, stateHandler)
	})

	fmt.Println("WebSocket server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
