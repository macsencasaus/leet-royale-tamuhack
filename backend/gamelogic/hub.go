package gamelogic

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
    m "leet-guys/messages"
)

var cId int = 0

type Hub struct {
	register   chan *client
	unregister chan *client
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

    queryParams := r.URL.Query()

    name := queryParams.Get("name")

    c := &client{
        id: cId,
        name: name,
        conn: conn,
    }

	conn.WriteJSON(m.NewHubGreetingMessage(c.playerInfo()))

    cId++
}
