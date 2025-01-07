package main

import (
	"chat-server/hub"
	"net/http"
)

func main() {
	h := hub.New()
	go h.Run()

	http.Handle("/ws", h)
	http.ListenAndServe(":8080", nil)
}
