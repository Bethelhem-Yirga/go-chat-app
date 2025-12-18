package main

import (
	"log"
	"net/http"

	"chat-app/internal/chat"
)

func main() {
	hub := chat.NewHub()
	go hub.Run()

	http.HandleFunc("/", chat.ServeHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(hub, w, r)
	})

	log.Println("Chat server running on http://localhost:8880")
	log.Fatal(http.ListenAndServe(":8880", nil))
}
