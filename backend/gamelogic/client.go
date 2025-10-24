package gamelogic

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	m "leet-guys/messages"
	tr "leet-guys/testrunner"

	"github.com/gorilla/websocket"
)

type Client struct {
	id   int
	name string

	conn *websocket.Conn

	hub *Hub

	roomWrite chan m.ServerMessage
	roomRead  chan ClientRoomMessage

	// room things
	done   bool
	closed bool
}

type PendingMessage struct {
	msg  m.ServerMessage
	done chan struct{}
}

func (c *Client) readPump() {
	for {
		var v json.RawMessage

		err := c.conn.ReadJSON(&v)
		if err != nil {
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

		var w m.ClientMessageWrapper

		err = json.Unmarshal(v, &w)
		if err != nil {
			c.log("Error unmarshaling data 1: %v", err)
			continue
		}

		var cm ClientRoomMessage

		switch w.Type {

		case m.ClientMessageTypeClientQuit:
			qm := m.ClientQuitMessage{}
			err = json.Unmarshal(w.Data, &qm)
			cm = Quit{c.id}

		case m.ClientMessageTypeSubmit:
			sm := m.SubmitMessage{}
			err = json.Unmarshal(w.Data, &sm)
			if err != nil {
				break
			}

			var l tr.Language
			switch sm.Language {
			case "python":
				l = tr.Python
			case "javascript":
				l = tr.Javascript
			case "cpp":
				l = tr.CPP
			}
			res, err := tr.RunTest([]byte(sm.Code), l, sm.QuestionId)
			if err != nil {
				c.log("Error running test: %v", err)
				continue
			}
			c.send(m.NewTestResultMessage(res))
			cm = Submit{
				playerId:   c.id,
				questionId: sm.QuestionId,
				results:    res,
			}

		case m.ClientMessageTypeSkipLobby, m.ClientMessageTypeSkipQuestion:
			cm = Skip{}

		}

		if err != nil {
			c.log("Error unmarshaling data 2: %v", err)
			continue
		}

		c.roomRead <- cm
	}

	if c.roomRead != nil && !c.closed {
		c.roomRead <- Quit{c.id}
	}
	close(c.roomWrite)
	c.log("readPump closed")
}

func (c *Client) send(msg m.ServerMessage) {
	c.roomWrite <- msg
}

func (c *Client) writePump() {
	for msg := range c.roomWrite {
		err := c.conn.WriteJSON(msg)

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.log("Tried to write, websocket closed")
				break
			}
			c.log("error writing server message to json %v", err)
			break
		}

		switch msg.(type) {
		case m.ClientEliminatedMessage, m.WinnerMessage:
			c.conn.Close()
		}
	}

	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	c.log("writePump closed")
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
