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

type clientId = int

const ClientsPerRoom = 40

const RoomWait = 60    // 1 minute
const Round1Time = 300 // 5 minutes
const Round2Time = 300 // 5 minutes
const Round3Time = 300 // 5 minutes
const Round4Time = 600 // 10 minutes

const RoundBetweenTime = 5

type room struct {
	id       int
	register chan *client

	roomRead chan m.ClientMessage

	clientsMu   sync.RWMutex
	clients     map[clientId]*client
	clientsDone map[clientId]bool

	stateMu sync.Mutex
	state   roomState

	totalPlayers atomic.Int32
	place        atomic.Int32

	// states
	waitingForPlayers
	round1Running
	round2Running
	round3Running
	round4Running
	gameEnded

	hub *Hub
}

func newRoom(id int) *room {
	r := &room{
		id:       id,
		register: make(chan *client),

		roomRead: make(chan m.ClientMessage),

		clients:     make(map[clientId]*client),
		clientsDone: make(map[clientId]bool),
	}

	var x atomic.Int32

	r.waitingForPlayers = waitingForPlayers{r, make(chan struct{})}
	r.round1Running = round1Running{r, make(chan struct{}), 0, &tr.Questions[0]}
	r.round2Running = round2Running{r, make(chan struct{}), 0, &tr.Questions[0]}
	r.round3Running = round3Running{r, make(chan struct{}), 0, &tr.Questions[0], &x}
	r.round4Running = round4Running{r, make(chan struct{}), 0, &tr.Questions[0]}
	r.gameEnded = gameEnded{r}

	return r
}

func (r *room) run() {
	defer func() {
		close(r.register)
		close(r.roomRead)
	}()

	r.setState(r.waitingForPlayers)

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
			r.state.handleClientMessage(msg)
		}
	}
}

type roomState interface {
	handleClientMessage(msg m.ClientMessage)
}

type waitingForPlayers struct {
	r             *room
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
	r          *room
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
	r          *room
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
	r                *room
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
		// TODO:
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
	r          *room
	timerDone  chan struct{}
	questionId int
	question   *tr.QuestionData
}

func (s round4Running) handleClientMessage(msg m.ClientMessage) {
	switch msg := msg.(type) {
	case m.ClientQuitMessage:
		s.r.unregisterClient(msg.PlayerId)
	case m.SubmitMessage:
		if s.r.runTestRunner(msg, s.questionId) {
			s.r.setClientDone(msg.PlayerId)
			s.timerDone <- struct{}{}
		}
	case m.SkipQuestionMessage:
		s.timerDone <- struct{}{}
	}
}

type gameEnded struct {
	r *room
}

func (s gameEnded) handleClientMessage(msg m.ClientMessage) {

}

func (r *room) runTestRunner(msg m.SubmitMessage, question int) bool {
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
		log.Fatal(err)
	}

	r.sendMessageTo(msg.PlayerId, m.NewTestResultMessage(&res))

	correct, total := res.NCorrect()

	c := r.getClient(msg.PlayerId)
	r.broadcast(m.NewUpdateClientStateMessage(
		c.playerInfo(),
		correct == total,
		correct,
		int(time.Now().Unix()),
	))

	return correct == total
}

func (r *room) setState(s roomState) {
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
	}
	r.stateMu.Lock()
	r.state = s
	r.stateMu.Unlock()
}

func (r *room) registerClient(c *client) bool {
	r.broadcast(m.NewClientJoinedMessage(c.playerInfo()))

	c.roomRead = r.roomRead
	go c.readPump()

	c.roomWrite <- m.NewRoomGreetingMessage(r.id, r.playersInfo())
	log.Println("Greeting Message Written")

	// r.clientsMu.Lock()
	r.clients[c.id] = c
	r.clientsDone[c.id] = false
	full := len(r.clients) == ClientsPerRoom
	// r.clientsMu.Unlock()

	r.totalPlayers.Add(1)
	r.place.Add(1)

	return full
}

func (r *room) unregisterClient(clientId clientId) {
    var c *client
    var ok bool
	{
		r.clientsMu.Lock()
		c, ok = r.clients[clientId]
		if !ok {
			return
		}
		delete(r.clients, clientId)
		delete(r.clientsDone, clientId)
		r.clientsMu.Unlock()
	}

	close(c.roomWrite)

	r.place.Add(-1)

	c.log("unregistered")

	r.broadcast(m.NewClientLeftMessage(c.playerInfo()))
}

func (r *room) broadcast(msg m.ServerMessage) {
	r.clientsMu.RLock()
	for _, ch := range r.clients {
		ch.roomWrite <- msg
	}
	r.clientsMu.RUnlock()
}

func (r *room) sendMessageTo(clientId clientId, msg m.ServerMessage) {
	r.clientsMu.RLock()
	r.clients[clientId].roomWrite <- msg
	r.clientsMu.RUnlock()
}

func (r *room) startRoundTimer(
	sec int,
	nextState roomState,
	endMessage m.ServerMessage,
	nextMessage m.ServerMessage,
	done chan struct{},
) {
	defer func() {
		r.setState(nextState)
		if nextMessage != nil {
			r.broadcast(nextMessage)
		}
		close(done)
	}()

	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	defer ticker.Stop()

	select {
	case <-done:
		break
	case <-ticker.C:
		break
	}

	r.broadcast(endMessage)

	r.handleEliminations()

	r.resetDone()

	time.Sleep(RoundBetweenTime * time.Second)
}

func (r *room) startCountdown(sec int, nextState roomState, broadcastMessage m.ServerMessage, done chan struct{}) {
	defer func() {
		r.setState(nextState)
		r.broadcast(broadcastMessage)
		close(done)
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for i := sec; i >= 1; i-- {
		select {
		case <-done:
			return
		case <-ticker.C:
			r.broadcast(m.NewCountdownMessage(i))
		}
	}
}

func (r *room) handleEliminations() {
	r.clientsMu.Lock()
	defer r.clientsMu.Unlock()
	toDelete := []int{}
	for cId, finished := range r.clientsDone {
		if !finished {
			c := r.clients[cId]
			r.clients[cId].roomWrite <- m.NewClientEliminatedMessage(
				c.playerInfo(),
				int(r.place.Load()),
				int(r.totalPlayers.Load()),
			)

			delete(r.clients, cId)
			toDelete = append(toDelete, cId)
			close(c.roomWrite)
			r.place.Add(-1)
			c.log("eliminated")
			for _, ch := range r.clients {
				ch.roomWrite <- m.NewClientLeftMessage(c.playerInfo())
			}
		}
	}
	for _, cId := range toDelete {
		delete(r.clientsDone, cId)
	}
}

func (r *room) setClientDone(clientId int) {
	r.clientsMu.Lock()
	r.clientsDone[clientId] = true
	r.clientsMu.Unlock()
}

func (r *room) isClientDone(clientId int) bool {
	r.clientsMu.RLock()
	done := r.clientsDone[clientId]
	r.clientsMu.RUnlock()
	return done
}

func (r *room) resetDone() {
	r.clientsMu.Lock()
	for clientsId := range r.clientsDone {
		r.clientsDone[clientsId] = false
	}
	r.clientsMu.Unlock()
}

func (r *room) playersInfo() []m.PlayerInfo {
	r.clientsMu.RLock()
	playersInfo := make([]m.PlayerInfo, 0, len(r.clients))
	for _, client := range r.clients {
		playersInfo = append(playersInfo, client.playerInfo())
	}
	r.clientsMu.RUnlock()
	return playersInfo
}

func (r *room) isOpen() bool {
	r.clientsMu.RLock()
	if len(r.clients) == ClientsPerRoom {
		return false
	}
	r.clientsMu.RUnlock()
	r.stateMu.Lock()
	_, ok := r.state.(waitingForPlayers)
	r.stateMu.Unlock()
	return ok
}

func (r *room) isRunning() bool {
	r.stateMu.Lock()
	running := r.state != nil
	r.stateMu.Unlock()
	return running
}

func (r *room) getClient(clientId clientId) *client {
	r.clientsMu.RLock()
	client := r.clients[clientId]
	r.clientsMu.RUnlock()
	return client
}

func (r *room) log(format string, v ...any) {
	log.Printf("room %d: %s", r.id, fmt.Sprintf(format, v...))
}
