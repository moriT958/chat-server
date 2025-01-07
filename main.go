package main

import (
	"chat-server/hub"
	"chat-server/server"
)

func main() {
	h := hub.New()
	go h.Run()

	s := server.New(h)
	s.Run()
}
