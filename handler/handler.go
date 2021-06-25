package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/madeindra/golang-websocket/model"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var ps = &model.PubSub{}

func WebsocketHandler(ctx *gin.Context) {
	// trust all origin to avoid CORS
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// upgrades connection to websocket
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// create new client & add to client list
	client := model.Client{
		ID:         uuid.Must(uuid.NewRandom()).String(),
		Connection: conn,
	}
	ps.AddClient(client)

	// message handling
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			ps.RemoveClient(client)
			return
		}
		ps.ProcessMessage(client, messageType, p)
	}
}
