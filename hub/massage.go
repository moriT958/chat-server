package hub

import "encoding/json"

type Message struct {
	Title string      `json:"title"`
	Data  interface{} `json:"data"`
}

func (m *Message) bytes() ([]byte, error) {
	return json.Marshal(m)
}
