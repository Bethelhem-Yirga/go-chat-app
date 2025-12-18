package chat

import "sync"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.Lock()
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					delete(h.clients, c)
					close(c.send)
				}
			}
			h.mu.Unlock()
		}
	}
}

func anyChannel(h *Hub) interface{} {
	select {
	case c := <-h.register:
		h.mu.Lock()
		h.clients[c] = true
		h.mu.Unlock()
	case c := <-h.unregister:
		h.mu.Lock()
		delete(h.clients, c)
		close(c.send)
		h.mu.Unlock()
	case msg := <-h.broadcast:
		h.mu.Lock()
		for c := range h.clients {
			select {
			case c.send <- msg:
			default:
				delete(h.clients, c)
				close(c.send)
			}
		}
		h.mu.Unlock()
	}
	return nil
}
