package gamelogic

import (
	"encoding/json"
	"fmt"
	"log"

	m "leet-guys/messages"

	"github.com/gorilla/websocket"
)

type client struct {
	id   int
	name string

	conn *websocket.Conn
	hub  *Hub

	roomWrite chan m.ServerMessage
	roomRead  chan m.ClientMessage
}

func (c *client) readPump() {
	defer func() {
		if c.roomRead != nil {
			c.roomRead <- m.ClientQuitMessage{PlayerId: c.id}
		}
		c.conn.Close()
		c.log("readPump closed")
	}()

	for {
		var v json.RawMessage

		err := c.conn.ReadJSON(&v)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.log("Tried to read, websocket closed")
				break
			}
			c.log("error reading connection json: %v", err)
			break
		}

		c.log("message received: %s", v)

		var w m.ClientMessageWrapper

		err = json.Unmarshal(v, &w)
		if err != nil {
			c.log("Error unmarshaling data 1: %v", err)
			continue
		}

		var cm m.ClientMessage

		switch w.Type {
		case m.ClientMessageTypeClientQuit:
            qm := m.ClientQuitMessage{}
			err = json.Unmarshal(w.Data, &qm)
            cm = qm
		case m.ClientMessageTypeSubmit:
            sm := m.SubmitMessage{}
			err = json.Unmarshal(w.Data, &sm)
            cm = sm
		case m.ClientMessageTypeSkipLobby:
			cm = m.SkipLobbyMessage{}
		case m.ClientMessageTypeSkipQuestion:
			cm = m.SkipQuestionMessage{}
		}

		if err != nil {
			c.log("Error unmarshaling data 2: %v", err)
			continue
		}

		c.roomRead <- cm
	}
}

func (c *client) writePump() {
	defer func() {
		c.conn.Close()
		c.log("writePump closed")
	}()

	for {
		msg, ok := <-c.roomWrite
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		err := c.conn.WriteJSON(msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.log("Tried to write, websocket closed")
				break
			}
			c.log("error writing server message to json")
			return
		}
	}
}

func (c *client) playerInfo() m.PlayerInfo {
	return m.PlayerInfo{
		Id:   c.id,
		Name: c.name,
	}
}

func (c *client) log(format string, v ...any) {
	log.Printf("client %d: %s", c.id, fmt.Sprintf(format, v...))
}
