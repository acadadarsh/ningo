package main

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func main() {
	hub := NewHub()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Println("websocket accept error:", err)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "")
		role := r.URL.Query().Get("role")
		client := &Client{
			conn: conn,
			send: make(chan string, 16),
			role: role,
		}
		hub.Add(client)
		defer hub.Remove(client)
		if role == "extension" {
			go client.WriteLoop(r.Context())
		}

		for {
			messageType, data, err := conn.Read(r.Context())
			if err != nil {
				log.Println("websocket read error:", err)
				return
			}
			if messageType != websocket.MessageText {
				log.Println("ignoring non -text message")
				continue
			}

			hub.Broadcast(r.Context(), string(data))
		}
	})

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
