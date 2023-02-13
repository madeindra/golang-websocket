package websocket

import "github.com/gorilla/websocket"

// Subscription is a struct for each topic and which client subscribe to it
type Subscription struct {
	Topic   string
	Clients *[]Client
}

// Client is a struct that describe the clients' ID and their connection
type Client struct {
	ID         string
	Connection *websocket.Conn
}

// Message is a struct for message to be sent by the client
type Message struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}
