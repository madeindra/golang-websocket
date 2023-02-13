package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/madeindra/golang-websocket/model"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var server = &model.Server{}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
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

	// create new client & add to client list
	client := model.Client{
		ID:         uuid.Must(uuid.NewRandom()).String(),
		Connection: conn,
	}

	// greet the new client
	server.Send(&client, "Server: Welcome! Your ID is "+client.ID)

	// message handling
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			server.RemoveClient(client)
			return
		}
		server.ProcessMessage(client, messageType, p)
	}
}
