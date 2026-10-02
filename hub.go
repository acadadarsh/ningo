package main

import (
	"context"
	"github.com/coder/websocket"
	"sync"
)

type Client struct {
	conn *websocket.Conn
	send chan string
	role string
}

func (c *Client) WriteLoop(ctx context.Context) {

	for {
		select {
		case queuedImage := <-c.send:
			err := c.conn.Write(
				ctx, 
				websocket.MessageText,
				[]byte(queuedImage),
			)
			if err != nil {
				return
			}

		case <-ctx.Done():
			return

		}
	}
}

type Hub struct {
	mu      sync.Mutex
	clients map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
	}
}

func (h *Hub) Add(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = struct{}{}
}

func (h *Hub) Remove(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client)
}

func (h *Hub) Broadcast(ctx context.Context, imageBase64 string) {
	h.mu.Lock()
	clients := make([]*Client, 0, len(h.clients))

	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.Unlock()

	for _, client := range clients {

		if client.role != "extension" {
			continue
		}
		select {
		case client.send <- imageBase64:
		case <-ctx.Done():
			return
		}
	}
}
