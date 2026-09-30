package main

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func main(){
	hub:=NewHub()

	http.HandleFunc("/ws",func(w http.ResponseWriter,r *http.Request){
		conn,err:=websocket.Accept(w,r,nil)
		if err!=nil{
			log.Println("websocket accept error:", err)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure,"")
		hub.Add(conn)
		defer hub.Remove(conn)

		for {
			messageType, data, err:=conn.Read(r.Context())
			if err!=nil{
				log.Println("websocket read error:", err)
				return
			}
			if messageType!=websocket.MessageText{
				log.Println("ignoring non -text message")
				continue
			}

			hub.Broadcast(r.Context(),string(data))
		}
	})

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080",nil))
}