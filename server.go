package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Player struct {
	ID     string
	Conn   *websocket.Conn
	IsReady bool
}

var players = make(map[string]*Player)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	player := &Player{
		ID:   generatePlayerID(),
		Conn: conn,
	}
	players[player.ID] = player

	for {
		// Handle incoming messages from the player
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(players, player.ID)
			break
		}
		fmt.Printf("Received: %s\n", msg)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.ListenAndServe(":8080", nil)
}

