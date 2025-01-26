package gamelogic

import (
	"log"
	"sync"
	"time"

	m "leet-guys/messages"
)

type clientId = int

const ClientsPerLobby = 40

const Round1Time = 300 // 5 minutes
const Round2Time = 300 // 5 minutes
const Round3Time = 300 // 5 minutes
const Round4Time = 600 // 10 minutes

type room struct {
	id       int
	register chan *client

	roomRead chan m.ClientMessage

	clientsMu sync.RWMutex
	clients   map[clientId]*client
	clientsDone   map[clientId]bool

	stateMu sync.Mutex
	state   roomState

	// states
	waitingForPlayers
	round1Running
	round2Running
	round3Running
	round4Running
	gameEnded
}

func newRoom(id int) *room {
	r := &room{
		id:       id,
		register: make(chan *client),

		clients: make(map[clientId]*client),
	}

	r.waitingForPlayers = waitingForPlayers{r, make(chan bool)}
	r.round1Running = round1Running{r, make(chan bool)}
	r.round2Running = round2Running{r, make(chan bool)}
	r.round3Running = round3Running{r, make(chan bool)}
	r.round4Running = round4Running{r, make(chan bool)}
	r.gameEnded = gameEnded{r}

	return r
}

func (r *room) run() {
	defer func() {
		close(r.register)
		close(r.roomRead)
	}()

	r.state = r.waitingForPlayers

	for {
		select {
		case client := <-r.register:
			_, ok := r.state.(waitingForPlayers)
			if !ok {
				log.Fatalf("unable to register with state %T", r.state)
			}
			if r.registerClient(client) {
				r.broadcast(m.NewRoundStartMessage(1, Round1Time))
				r.setState(r.round1Running)
				go r.startTimer(Round1Time, r.round2Running, r.round1Running.timerStop)
			}
		case msg := <-r.roomRead:
			r.state.handleClientMessage(msg)
		}
	}
}

type roomState interface {
	handleClientMessage(msg m.ClientMessage)
}

type waitingForPlayers struct {
	r             *room
	countDownStop chan bool
}

func (s waitingForPlayers) handleClientMessage(msg m.ClientMessage) {
	switch msg := msg.(type) {
	case m.ClientQuitMessage:
		s.r.unregisterClient(msg.PlayerId)
	}
}

type round1Running struct {
	r         *room
	timerStop chan bool
}

func (s round1Running) handleClientMessage(msg m.ClientMessage) {
    // TODO:
}

type round2Running struct {
	r         *room
	timerStop chan bool
}

func (s round2Running) handleClientMessage(msg m.ClientMessage) {

}

type round3Running struct {
	r         *room
	timerStop chan bool
}

type round4Running struct {
	r         *room
	timerStop chan bool
}

type gameEnded struct {
	r *room
}

func (r *room) setState(s roomState) {
	r.stateMu.Lock()
	r.state = s
	r.stateMu.Unlock()
}

func (r *room) registerClient(c *client) bool {
	r.broadcast(m.NewClientJoinedMessage(c.playerInfo()))
	roomWrite := make(chan m.ServerMessage)
	c.roomWrite = roomWrite
	r.clientsMu.Lock()
	r.clients[c.id] = c
	full := len(r.clients) == ClientsPerLobby
	r.clientsMu.Unlock()
	return full
}

func (r *room) unregisterClient(clientId clientId) {
	r.clientsMu.Lock()
	c := r.clients[clientId]
	delete(r.clients, clientId)
	r.clientsMu.Unlock()
	close(c.roomWrite)
	r.broadcast(m.NewClientLeftMessage(c.playerInfo()))
}

func (r *room) broadcast(msg m.ServerMessage) {
	r.clientsMu.RLock()
	for _, ch := range r.clients {
		ch.roomWrite <- msg
	}
	r.clientsMu.RUnlock()
}

func (r *room) sentMessageTo(clientId clientId, msg m.ServerMessage) {
	r.clientsMu.RLock()
	r.clients[clientId].roomWrite <- msg
	r.clientsMu.RUnlock()
}

func (r *room) startTimer(sec int, successState roomState, stop chan bool) {
	defer func() {
		close(stop)
	}()

	ticker := time.NewTicker(time.Duration(sec) * 1000)

	for {
		select {
		case success := <-stop:
			if success {
				r.setState(successState)
			}
			return
		case <-ticker.C:
			r.setState(successState)
			return
		}
	}
}

func (r *room) startCountDown(sec int, successState roomState, stop chan bool) {
	defer func() {
		close(stop)
	}()

	ticker := time.NewTicker(1000)

	for i := sec; i >= 1; i++ {
		select {
		case success := <-stop:
			if success {
				r.setState(successState)
			}
			return
		case <-ticker.C:
			r.broadcast(m.NewCountdownMessage(i))
			return
		}
	}
}
