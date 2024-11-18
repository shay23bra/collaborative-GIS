package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "http://localhost:8080"
		},
	}
)

type Message struct {
	Type        string      `json:"type"`
	Username    string      `json:"username"`
	Coordinates [][]float64 `json:"coordinates"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		msg.Username = r.URL.Query().Get("username")
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		messageData, _ := json.Marshal(msg)
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, messageData)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
