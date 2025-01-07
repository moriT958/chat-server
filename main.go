package main

import (
	"net/http"
	"simple-message-server/hub"
)

func main() {
	h := hub.New()
	go h.Run()

	http.Handle("/ws", h)
	http.ListenAndServe(":8080", nil)
}
