package main

import (
	"net/http"

	"github.com/madeindra/golang-websocket/internal/websocket"
)

func main() {
	http.HandleFunc("/socket", websocket.HandleWS)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
