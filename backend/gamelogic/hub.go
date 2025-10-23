package gamelogic

import (
	"log"
	"net/http"
	"sync"

	m "leet-guys/messages"

	"github.com/gorilla/websocket"
)

var cId int = 0

type Hub struct {
	registerClientQueue chan *Client
	unregisterRoomQueue chan *Room

	roomsMu sync.Mutex
	rooms   []*Room
}

func NewHub() *Hub {
	h := &Hub{
		registerClientQueue: make(chan *Client, ClientsPerRoom*10),
		rooms:               []*Room{},
	}

	return h
}

func (h *Hub) Run() {
	log.Println("Hub Started")

	for client := range h.registerClientQueue {
		log.Println("Client Recieved in Hub")

		r := h.findBestRoom()

		r.register <- client
	}
}

func (h *Hub) findBestRoom() *Room {
	maxFill := -1
	var r *Room

	for _, room := range h.rooms {
		if room.isOpen() && len(room.clients) > maxFill {
			r = room
			maxFill = len(room.clients)
		}
	}

	if r != nil {
		return r
	}

	log.Println("New room created")

	r = newRoom(len(h.rooms), h)
	h.rooms = append(h.rooms, r)

	go r.run()

	return r
}

func (h *Hub) unregisterRoom(r *Room) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	for i, room := range h.rooms {
		if room == r {
			h.rooms[i] = h.rooms[len(h.rooms)-1]
			h.rooms = h.rooms[:len(h.rooms)-1]
			return
		}
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
