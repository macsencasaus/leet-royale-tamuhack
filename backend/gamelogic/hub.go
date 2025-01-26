package gamelogic

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	m "leet-guys/messages"
)

var cId int = 0

const MaxRooms = 10

type Hub struct {
	register chan *client

	rooms [MaxRooms]*room
}

func NewHub() *Hub {
	h := &Hub{
		register: make(chan *client, ClientsPerRoom*10),
	}

	for i := 0; i < MaxRooms; i++ {
		h.rooms[i] = newRoom(i)
	}

	return h
}

func (h *Hub) Run() {
	defer func() {
		close(h.register)
	}()
	log.Println("Hub Started")
	for {
		client := <-h.register
		log.Println("Client Recieved in Hub")

		for _, room := range h.rooms {
			if !room.isRunning() {
				go room.run()
			} else if !room.isOpen() {
				continue
			}
			client.log("registered to room %d", room.id)
			room.register <- client
			break
		}

		time.Sleep(100 * time.Millisecond)
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

	c := &client{
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

	h.register <- c

	cId++
}
