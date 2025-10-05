package gamelogic

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	m "leet-guys/messages"

	"github.com/gorilla/websocket"
)

type Client struct {
	id   int
	name string

	conn       *websocket.Conn
	connClosed bool

	hub *Hub

	roomWrite chan PendingMessage
	roomRead  chan m.ClientMessage

	// room things
	done   bool
	closed bool

	closedMu sync.Mutex
}

type PendingMessage struct {
	msg  m.ServerMessage
	done chan struct{}
}

func (c *Client) readPump() {
	defer func() {
		if c.roomRead != nil {
			c.roomRead <- m.ClientQuitMessage{PlayerId: c.id}
		}
		c.close()
	}()

	for {
		var v json.RawMessage

		err := c.conn.ReadJSON(&v)
		if err != nil {
			c.connClosed = true
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) || errors.Is(err, net.ErrClosed) {
				c.log("Tried to read, websocket closed")
				break
			}
			c.log("error reading connection json: %+v", err.Error())
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

// send message (non blocking)
func (c *Client) send(msg m.ServerMessage) {
	c.roomWrite <- PendingMessage{msg: msg, done: nil}
}

// send message and returns done channel
func (c *Client) sendAsync(msg m.ServerMessage) chan struct{} {
	done := make(chan struct{})
	c.roomWrite <- PendingMessage{msg: msg, done: done}
	return done
}

func (c *Client) writePump() {
	for pendingMsg := range c.roomWrite {
		c.closedMu.Lock()
		err := c.conn.WriteJSON(pendingMsg.msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.log("Tried to write, websocket closed")
				break
			}
			c.log("error writing server message to json %v", err)
			return
		}
		c.closedMu.Unlock()

		if pendingMsg.done != nil {
			close(pendingMsg.done)
		}
	}

	if !c.connClosed {
		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	}
	c.log("writePump closed")
}

func (c *Client) close() {
	c.closedMu.Lock()
	defer c.closedMu.Unlock()

	if !c.closed {
		close(c.roomWrite)
	}
	if !c.connClosed {
		c.conn.Close()
	}
	c.closed = true
	c.connClosed = true
}

func (c *Client) playerInfo() m.PlayerInfo {
	return m.PlayerInfo{
		Id:   c.id,
		Name: c.name,
	}
}

func (c *Client) log(format string, v ...any) {
	log.Printf("client %d: %s", c.id, fmt.Sprintf(format, v...))
}
