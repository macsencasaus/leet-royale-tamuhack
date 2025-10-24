package gamelogic

import (
	"log"
	"net/http"

	m "leet-guys/messages"

	"github.com/gorilla/websocket"
)

var cId int = 0
var rId int = 0

type Hub struct {
	registerClientQueue chan *Client

	room *Room
}

func NewHub() *Hub {
	h := &Hub{
		registerClientQueue: make(chan *Client, ClientsPerRoom*10),
	}

	return h
}

func (h *Hub) Run() {
	log.Println("Hub Started")

	for client := range h.registerClientQueue {
		log.Println("Client Recieved in Hub")

		if h.room == nil || !h.room.isOpen() {
			h.room = newRoom(rId, h)
			go h.room.run()
			rId++
		}

		h.room.register <- client
	}
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

	log.Println("New WS Connection")

	queryParams := r.URL.Query()

	name := queryParams.Get("name")

	c := &Client{
		id:   cId,
		name: name,
		conn: conn,

		roomWrite: make(chan m.ServerMessage),
	}

	var cmw m.ClientMessageWrapper

	err = c.conn.ReadJSON(&cmw)
	if err != nil {
		c.log("error: %v", err)
		return
	}

	if cmw.Type != m.ClientMessageTypeReady {
		c.log("Not ready type")
		return
	}

	log.Println("Ready Message Recieved")

	conn.WriteJSON(m.NewHubGreetingMessage(c.playerInfo()))

	log.Println("Sent Hub Greeting")

	go c.writePump()

	h.registerClientQueue <- c

	cId++
}
