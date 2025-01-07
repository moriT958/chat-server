package hub

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ReadBufferSize  int = 1024
	WriteBufferSize int = 1024
)

type Hub struct {
	sync.Mutex
	connections map[*connection]bool

	upgrader   websocket.Upgrader
	register   chan *connection
	unregister chan *connection

	Broadcast chan *Message
}

func New() *Hub {
	h := &Hub{
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
		Broadcast:   make(chan *Message),
	}

	h.upgrader = websocket.Upgrader{
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
	}

	return h
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.doRegister(c)
		case c := <-h.unregister:
			h.doUnregister(c)
		case m := <-h.Broadcast:
			h.doBroadcast(m)
		}
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] failed to upgrade connection:", err)
		return
	}

	conn := &connection{send: make(chan []byte, 256), ws: ws, hub: h}
	h.register <- conn
	go conn.listenWrite()
	conn.listenRead()
}

func (h *Hub) doRegister(c *connection) {
	h.Lock()
	defer h.Unlock()

	h.connections[c] = false
}

func (h *Hub) doBroadcast(m *Message) {
	h.Lock()
	defer h.Unlock()

	bytes, err := m.bytes()
	if err != nil {
		log.Printf("failed to marshal message: %+v, reason: %s\n", m, err)
		return
	}
	for conn := range h.connections {
		conn.send <- bytes
	}
}

func (h *Hub) doUnregister(conn *connection) {
	h.Lock()
	defer h.Unlock()

	_, ok := h.connections[conn]
	if !ok {
		log.Println("cannot unregister connection, it is not registered.")
		return
	}

	log.Println("closing socket connection")
	conn.close()
	delete(h.connections, conn)
}
