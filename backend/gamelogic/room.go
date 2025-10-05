package gamelogic

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	m "leet-guys/messages"
	tr "leet-guys/testrunner"
)

// TODO: make locking vs unlocked API consistent
// probably make everything unlocked by default

type ClientId = int

const ClientsPerRoom = 40

const RoomWait = 60    // 1 minute
const Round1Time = 300 // 5 minutes
const Round2Time = 300 // 5 minutes
const Round3Time = 300 // 5 minutes
const Round4Time = 600 // 10 minutes

const RoundBetweenTime = 5

type Room struct {
	id       int
	register chan *Client

	roomRead chan m.ClientMessage

	clientsMu sync.RWMutex

	// contains all clients even those disconnected
	clients map[ClientId]*Client

	activeClientCount int

	stateMu sync.RWMutex
	state   RoomState

	// states
	waitingForPlayers
	round1Running
	round2Running
	round3Running
	round4Running
	gameEnded

	hub *Hub

	done          chan struct{}
	shutdownTimer chan struct{}
}

func newRoom(id int, hub *Hub) *Room {
	r := &Room{
		id:       id,
		register: make(chan *Client),

		roomRead: make(chan m.ClientMessage),

		clients: make(map[ClientId]*Client),

		hub: hub,
	}

	var x atomic.Int32

	r.waitingForPlayers = waitingForPlayers{r, make(chan struct{})}
	r.round1Running = round1Running{r, make(chan struct{}), 0, &tr.Questions[0]}
	r.round2Running = round2Running{r, make(chan struct{}), 4, &tr.Questions[4]}
	r.round3Running = round3Running{r, make(chan struct{}), 5, &tr.Questions[5], &x}
	r.round4Running = round4Running{r, make(chan struct{}), 7, &tr.Questions[7]}
	r.gameEnded = gameEnded{r}

	return r
}

func (r *Room) run() {
	r.log("running")
	r.setState(r.waitingForPlayers)

	r.done = make(chan struct{})
	r.shutdownTimer = make(chan struct{})

	defer func() {
		close(r.register)
		close(r.roomRead)
	}()

	for {
		select {
		case client := <-r.register:
			if !r.isOpen() {
				log.Fatalf("unable to register with state %T", r.state)
			}
			if r.registerClient(client) {
				r.countdownDone <- struct{}{}
			}
		case msg := <-r.roomRead:
			r.stateMu.RLock()
			go r.state.handleClientMessage(msg)
			r.stateMu.RUnlock()
		case <-r.done:
			r.log("closing")
			r.hub.unregisterRoom(r)
			return
		}
	}

}

type RoomState interface {
	handleClientMessage(msg m.ClientMessage)
}

type waitingForPlayers struct {
	r             *Room
	countdownDone chan struct{}
}

func (s waitingForPlayers) handleClientMessage(msg m.ClientMessage) {
	switch msg := msg.(type) {
	case m.ClientQuitMessage:
		s.r.unregisterClient(msg.PlayerId)
	case m.SkipLobbyMessage:
		s.countdownDone <- struct{}{}
	}
}

type round1Running struct {
	r          *Room
	timerDone  chan struct{}
	questionId int
	question   *tr.QuestionData
}

func (s round1Running) handleClientMessage(msg m.ClientMessage) {
	switch msg := msg.(type) {
	case m.ClientQuitMessage:
		s.r.unregisterClient(msg.PlayerId)
	case m.SubmitMessage:
		if s.r.runTestRunner(msg, s.questionId) {
			s.r.setClientDone(msg.PlayerId)
		}
	case m.SkipQuestionMessage:
		s.timerDone <- struct{}{}
	}
}

type round2Running struct {
	r          *Room
	timerDone  chan struct{}
	questionId int
	question   *tr.QuestionData
}

func (s round2Running) handleClientMessage(msg m.ClientMessage) {
	switch msg := msg.(type) {
	case m.ClientQuitMessage:
		s.r.unregisterClient(msg.PlayerId)
	case m.SubmitMessage:
		if s.r.runTestRunner(msg, s.questionId) {
			s.r.setClientDone(msg.PlayerId)
		}
	case m.SkipQuestionMessage:
		s.timerDone <- struct{}{}
	}
}

type round3Running struct {
	r                *Room
	timerDone        chan struct{}
	questionId       int
	question         *tr.QuestionData
	clientsSubmitted *atomic.Int32
}

func (s round3Running) handleClientMessage(msg m.ClientMessage) {
	switch msg := msg.(type) {
	case m.ClientQuitMessage:
		s.r.unregisterClient(msg.PlayerId)
	case m.SubmitMessage:
		if s.r.runTestRunner(msg, s.questionId) {
			if !s.r.isClientDone(msg.PlayerId) {
				s.clientsSubmitted.Add(1)
				s.r.setClientDone(msg.PlayerId)
				if int(s.clientsSubmitted.Load()) == 10 {
					s.timerDone <- struct{}{}
				}
			}
		}
	case m.SkipQuestionMessage:
		s.timerDone <- struct{}{}
	}
}

type round4Running struct {
	r          *Room
	timerDone  chan struct{}
	questionId int
	question   *tr.QuestionData
}

func (s round4Running) handleClientMessage(msg m.ClientMessage) {
	switch msg := msg.(type) {
	case m.ClientQuitMessage:
		s.r.unregisterClient(msg.PlayerId)
	case m.SubmitMessage:
		// TODO: let players continue to play and communicate place
		if s.r.runTestRunner(msg, s.questionId) {
			s.r.clientsMu.Lock()

			c := s.r.clients[msg.PlayerId]
			c.done = true

			done := c.sendAsync(m.NewWinnerMessage())

			s.r.clientsMu.Unlock()

			// Make sure we send the message before we close the room.
			// There is a chance the client's connection closes
			// before the message gets in its message queue, but
			// client's writePump should handle the error
			<-done

			// FIXME: race condition can cause write on closed channel
			s.timerDone <- struct{}{}
		}
	case m.SkipQuestionMessage:
		s.timerDone <- struct{}{}
	}
}

type gameEnded struct {
	r *Room
}

func (s gameEnded) handleClientMessage(msg m.ClientMessage) {

}

// runs tests for given submit message, returns whether test passed
func (r *Room) runTestRunner(msg m.SubmitMessage, question int) bool {
	var l tr.Language
	switch msg.Language {
	case "python":
		l = tr.Python
	case "javascript":
		l = tr.Javascript
	case "cpp":
		l = tr.CPP
	}
	res, err := tr.RunTest([]byte(msg.Code), l, question)
	if err != nil {
		r.log(err.Error())
	}

	r.clientsMu.RLock()
	defer r.clientsMu.RUnlock()

	c := r.clients[msg.PlayerId]

	c.send(m.NewTestResultMessage(&res))

	correct, total := res.NCorrect()
	passed := correct == total

	r.broadcastUnlocked(m.NewUpdateClientStateMessage(
		c.playerInfo(),
		passed,
		correct,
		int(time.Now().Unix()),
	))

	return passed
}

func (r *Room) setState(s RoomState) {
	switch s.(type) {
	case waitingForPlayers:
		go r.startCountdown(
			RoomWait,
			r.round1Running,
			m.NewRoundStartMessage(1, Round1Time, r.round1Running.question),
			r.waitingForPlayers.countdownDone,
		)
	case round1Running:
		go r.startRoundTimer(
			Round1Time,
			r.round2Running,
			m.NewRoundEndMessage(1),
			m.NewRoundStartMessage(2, Round2Time, r.round2Running.question),
			r.round1Running.timerDone,
		)
	case round2Running:
		go r.startRoundTimer(
			Round2Time,
			r.round3Running,
			m.NewRoundEndMessage(2),
			m.NewRoundStartMessage(3, Round3Time, r.round3Running.question),
			r.round2Running.timerDone,
		)
	case round3Running:
		go r.startRoundTimer(
			Round3Time,
			r.round4Running,
			m.NewRoundEndMessage(3),
			m.NewRoundStartMessage(4, Round4Time, r.round4Running.question),
			r.round3Running.timerDone,
		)
	case round4Running:
		go r.startRoundTimer(
			Round4Time,
			r.gameEnded,
			m.NewRoundEndMessage(4),
			nil,
			r.round4Running.timerDone,
		)
	case gameEnded:
		r.clientsMu.Lock()
		for _, c := range r.clients {
			if !c.closed {
				c.close()
			}
		}
		r.clientsMu.Unlock()
		close(r.done)
	}

	r.stateMu.Lock()
	defer r.stateMu.Unlock()
	r.state = s
}

// registers client into lobby, returns whether room is full after register
func (r *Room) registerClient(c *Client) bool {
	r.clientsMu.Lock()
	defer r.clientsMu.Unlock()

	r.broadcastUnlocked(m.NewClientJoinedMessage(c.playerInfo()))

	c.roomRead = r.roomRead
	go c.readPump()

	c.send(m.NewRoomGreetingMessage(r.id, r.playersInfoUnlocked()))
	r.log("Greeting Message Written")

	c.done = false
	c.closed = false

	r.clients[c.id] = c

	r.activeClientCount++
	full := r.activeClientCount == ClientsPerRoom

	return full
}

// unregister client from lobby
// sends done message if all clients have left the lobby
// is the only function response for closing a client
func (r *Room) unregisterClient(clientId ClientId) {
	r.clientsMu.Lock()
	defer r.clientsMu.Unlock()
	r.unregisterClientUnlocked(clientId)
}

// requires clientsMu
func (r *Room) unregisterClientUnlocked(clientId ClientId) {
	c, ok := r.clients[clientId]
	if !ok {
		return
	}

	if c.closed {
		return
	}

	c.close()
	r.activeClientCount--

	c.log("unregistered")

	r.broadcastUnlocked(m.NewClientLeftMessage(c.playerInfo()))

	if r.activeClientCount == 0 {
		r.log("all clients left")
		close(r.done)
	}
}

// func (r *Room) broadcast(msg m.ServerMessage) {
// 	r.clientsMu.RLock()
// 	defer r.clientsMu.RUnlock()
// 	r.broadcastUnlocked(msg)
// }

// Requires: clientsMu
func (r *Room) broadcastUnlocked(msg m.ServerMessage) {
	for _, c := range r.clients {
		if !c.closed {
			c.send(msg)
		}
	}
}

func (r *Room) startRoundTimer(
	sec int,
	nextState RoomState,
	endMessage m.RoundEndMessage, // message to send when round ends (immediate)
	nextMessage m.ServerMessage, // message to send when going to next round (after cool down timer)
	done chan struct{},
) {
	defer close(done)

	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	defer ticker.Stop()

	select {

	case <-done:
		break

	case <-ticker.C:
		break

	case <-r.shutdownTimer:
		return
	}

	r.handleEliminations()

	r.clientsMu.Lock()

	for _, c := range r.clients {
		if c.closed {
			endMessage.EliminatedPlayers = append(endMessage.EliminatedPlayers, c.playerInfo())
		} else {
			endMessage.CurrentPlayers = append(endMessage.CurrentPlayers, c.playerInfo())
		}
	}

	r.broadcastUnlocked(endMessage)

	for _, c := range r.clients {
		if !c.closed {
			c.done = false
		}
	}

	r.clientsMu.Unlock()

	time.Sleep(RoundBetweenTime * time.Second)

	r.setState(nextState)

	r.clientsMu.Lock()
	if nextMessage != nil {
		r.broadcastUnlocked(nextMessage)
	}
	r.clientsMu.Unlock()
}

func (r *Room) startCountdown(sec int, nextState RoomState, broadcastMessage m.ServerMessage, done chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

timerLoop:
	for i := sec; i >= 1; i-- {
		select {
		case <-done:
			break timerLoop

		case <-ticker.C:
			r.clientsMu.Lock()
			r.broadcastUnlocked(m.NewCountdownMessage(i))
			r.clientsMu.Unlock()

		case <-r.shutdownTimer:
			return
		}
	}

	r.setState(nextState)

	r.clientsMu.Lock()
	r.broadcastUnlocked(broadcastMessage)
	r.clientsMu.Unlock()
}

func (r *Room) handleEliminations() {
	r.clientsMu.Lock()
	defer r.clientsMu.Unlock()

	type ClientDone struct {
		done chan struct{}
		c    *Client
	}

	dones := make([]ClientDone, 0, len(r.clients))

	for _, c := range r.clients {
		if !c.done && !c.closed {
			msg := m.NewClientEliminatedMessage(
				c.playerInfo(),
				r.activeClientCount,
				len(r.clients),
			)
			done := c.sendAsync(msg)
			dones = append(dones, ClientDone{done: done, c: c})
		}
	}

	for _, cd := range dones {
		<-cd.done
		r.unregisterClientUnlocked(cd.c.id)
		cd.c.log("eliminated")
	}
}

func (r *Room) setClientDone(clientId ClientId) {
	r.clientsMu.Lock()
	defer r.clientsMu.Unlock()
	r.clients[clientId].done = true
}

func (r *Room) isClientDone(clientId ClientId) bool {
	r.clientsMu.RLock()
	defer r.clientsMu.RUnlock()
	return r.clients[clientId].done
}

func (r *Room) playersInfoUnlocked() []m.PlayerInfo {
	playersInfo := make([]m.PlayerInfo, 0, len(r.clients))
	for _, client := range r.clients {
		playersInfo = append(playersInfo, client.playerInfo())
	}
	return playersInfo
}

func (r *Room) isOpen() bool {
	r.clientsMu.RLock()
	if len(r.clients) == ClientsPerRoom {
		return false
	}
	r.clientsMu.RUnlock()
	r.stateMu.RLock()
	defer r.stateMu.RUnlock()
	_, ok := r.state.(waitingForPlayers)
	return ok
}

func (r *Room) clientCount() int {
	r.clientsMu.RLock()
	defer r.clientsMu.RUnlock()
	return r.activeClientCount
}

func (r *Room) log(format string, v ...any) {
	log.Printf("room %d: %s", r.id, fmt.Sprintf(format, v...))
}
