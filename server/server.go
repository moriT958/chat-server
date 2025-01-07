package server

import (
	"chat-server/hub"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	Addr = "0.0.0.0:8080"
)

type ChatServer struct {
	http.Server
	hub *hub.Hub
}

func New(h *hub.Hub) *ChatServer {
	s := new(ChatServer)

	mux := http.NewServeMux()
	mux.Handle("/ws", h)

	s.Addr = Addr
	s.Handler = mux
	s.hub = h
	return s
}

func (s *ChatServer) Run() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	<-shutdown
	if err := s.Close(); err != nil {
		log.Println(err)
	}
}
