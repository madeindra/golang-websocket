package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/madeindra/golang-websocket/handler"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/socket", handler.WebsocketHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
