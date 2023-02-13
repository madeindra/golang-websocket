package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/madeindra/golang-websocket/internal/websocket"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/socket", websocket.HandleWS)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
