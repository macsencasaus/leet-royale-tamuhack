package gamelogic

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
}

func NewHub() *Hub {
	return nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusUpgradeRequired)
		log.Printf("upgrading error: %s\n", err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("Poop fart"))
}
