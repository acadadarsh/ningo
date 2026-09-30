package main

import (
	"context"
	"github.com/coder/websocket"
	"sync"
)

type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (h *Hub) Add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = struct{}{}
}

func (h *Hub) Remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
} 

func (h *Hub) Broadcast(ctx context.Context, imageBase64 string){
	h.mu.Lock()
	clients:=make([]*websocket.Conn,0,len(h.clients))

	for conn:= range h.clients{
		clients=append(clients, conn)
	}
	h.mu.Unlock()

	for _,conn:=range clients{
		err:=conn.Write(
			ctx,
			websocket.MessageText,
			[]byte(imageBase64),
		)
		if err!=nil{
			h.Remove(conn)
		}
	}

}

