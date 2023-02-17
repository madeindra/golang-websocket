package websocket

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Initialize server with empty subscription
var server = &Server{Subscriptions: make(Subscription)}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	// trust all origin to avoid CORS
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// upgrades connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed upgrading connection"))
		return
	}
	defer conn.Close()

	// create new client id
	clientID := uuid.New().String()

	// greet the new client
	server.Send(conn, fmt.Sprintf("Server: Welcome! Your ID is %s", clientID))

	// message handling
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			server.RemoveClient(clientID)
			return
		}
		fmt.Println("message type")
		server.ProcessMessage(conn, clientID, msg)
	}
}
