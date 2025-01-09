package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) 
var broadcast = make(chan Message)         


type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Println("WebSocket sunucusu 8080 portunda çalışıyor...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html") 
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Hata: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Hata: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
