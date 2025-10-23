package gamelogic

import (
	"fmt"
	"log"
	"time"

	m "leet-guys/messages"
	tr "leet-guys/testrunner"
)

type ClientId = int

const ClientsPerRoom = 40

const (
	RoomWait   = 60 * time.Second  // 1 minute
	Round1Time = 300 * time.Second // 5 minutes
	Round2Time = 300 * time.Second // 5 minutes
	Round3Time = 300 * time.Second // 5 minutes
	Round4Time = 600 * time.Second // 10 minutes
)

const RoundBetweenTime = 5 * time.Second

type Room struct {
	id       int
	register chan *Client

	roomRead chan m.ClientMessage

	// contains all clients even those disconnected
	clients           map[ClientId]*Client
	activeClientCount int

	state RoomState

	roundDone chan struct{}
	done      chan struct{}

	hub *Hub
}

func newRoom(id int, hub *Hub) *Room {
	r := &Room{
		id:       id,
		register: make(chan *Client),

		roomRead: make(chan m.ClientMessage),

		clients: make(map[ClientId]*Client),

		roundDone: make(chan struct{}),
		done:      make(chan struct{}),

		hub: hub,
	}

	return r
}

func (r *Room) run() {
	r.log("running")
	r.setWaitingForPlayers()

	currentRoundTimer := time.NewTimer(RoomWait)

	for {
		select {

		case client := <-r.register:
			if !r.isOpen() {
				// back to the hub
				r.hub.registerClientQueue <- client
				continue
			}
			if r.registerClient(client) {
				// if lobby filled, end countdown early
				if _, ok := r.state.(waitingForPlayers); ok {
					r.endRound()
				}
			}

		case msg := <-r.roomRead:
			switch msg := msg.(type) {

			case m.ClientQuitMessage:
				r.unregisterClient(r.clients[msg.PlayerId])

			case m.SubmitMessage:
				r.state.handleSubmitMessage(msg)
				if r.everyoneDone() {
					r.endRound()
				}

			case m.SkipQuestionMessage, m.SkipLobbyMessage:
				r.endRound()

			}

		case <-r.done:
			r.log("closing")
			r.hub.unregisterRoom(r)
			return

		case <-currentRoundTimer.C:
			r.endRound()

		case <-r.roundDone:
			r.log("Round done: %T", r.state)
			if !currentRoundTimer.Stop() {
				select {
				case <-currentRoundTimer.C:
				default:
				}
			}

			_, inWaitingState := r.state.(waitingForPlayers)

			if !inWaitingState {
				r.handleEliminations()
			}

			nextTimerDuration, round := r.nextState()

			if round == 0 {
				r.end()
				continue
			}

			currentRoundTimer.Reset(nextTimerDuration)

			r.roundDone = make(chan struct{})

			for _, c := range r.clients {
				c.done = false
			}

			if !inWaitingState {
				time.Sleep(RoundBetweenTime)
			}

			r.broadcast(m.NewRoundStartMessage(
				round,
				int(nextTimerDuration/time.Second),
				r.state.getQuestion(),
			))
		}
	}
}

func (r *Room) endRound() {
	select {
	case <-r.roundDone:
	default:
		close(r.roundDone)
	}
}

func (r *Room) end() {
	select {
	case <-r.done:
	default:
		close(r.done)
	}
}

type RoomState interface {
	handleSubmitMessage(msg m.SubmitMessage)
	getQuestion() *tr.QuestionData
}

type waitingForPlayers struct {
	r *Room
}

func (waitingForPlayers) getQuestion() *tr.QuestionData { return nil }

func (waitingForPlayers) handleSubmitMessage(m.SubmitMessage) {
}

type round1Running struct {
	r          *Room
	questionId int
	question   *tr.QuestionData
}

func (s round1Running) getQuestion() *tr.QuestionData { return s.question }

func (s round1Running) handleSubmitMessage(msg m.SubmitMessage) {
	if s.r.runTestRunner(msg, s.questionId) {
		s.r.clients[msg.PlayerId].done = true
	}
}

type round2Running struct {
	r          *Room
	questionId int
	question   *tr.QuestionData
}

func (s round2Running) getQuestion() *tr.QuestionData { return s.question }

func (s round2Running) handleSubmitMessage(msg m.SubmitMessage) {
	if s.r.runTestRunner(msg, s.questionId) {
		s.r.clients[msg.PlayerId].done = true
	}
}

type round3Running struct {
	r                *Room
	timerDone        chan struct{}
	questionId       int
	question         *tr.QuestionData
	clientsSubmitted int
}

func (s round3Running) getQuestion() *tr.QuestionData { return s.question }

func (s round3Running) handleSubmitMessage(msg m.SubmitMessage) {
	if s.r.runTestRunner(msg, s.questionId) {
		s.r.clients[msg.PlayerId].done = true
		s.clientsSubmitted++
		if s.clientsSubmitted == 10 {
			s.r.endRound()
		}
	}
}

type round4Running struct {
	r          *Room
	questionId int
	question   *tr.QuestionData
}

func (s round4Running) getQuestion() *tr.QuestionData { return s.question }

func (s round4Running) handleSubmitMessage(msg m.SubmitMessage) {
	if s.r.runTestRunner(msg, s.questionId) {
		c := s.r.clients[msg.PlayerId]
		c.done = true
		c.send(m.NewWinnerMessage())
		s.r.endRound()
	}
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

	c := r.clients[msg.PlayerId]

	c.send(m.NewTestResultMessage(&res))

	correct, total := res.NCorrect()
	passed := correct == total

	r.broadcast(m.NewUpdateClientStateMessage(
		c.playerInfo(),
		passed,
		correct,
		int(time.Now().Unix()),
	))

	return passed
}

func (r *Room) nextState() (time.Duration, int) {
	switch r.state.(type) {
	case waitingForPlayers:
		return r.setRound1Running(), 1
	case round1Running:
		return r.setRound2Running(), 2
	case round2Running:
		return r.setRound3Running(), 3
	case round3Running:
		return r.setRound4Running(), 4
	}
	return 0, 0
}

func (r *Room) setWaitingForPlayers() {
	s := waitingForPlayers{r: r}
	r.state = s
}

func (r *Room) setRound1Running() time.Duration {
	s := round1Running{
		r:          r,
		questionId: 0,
		question:   &tr.Questions[0],
	}
	r.state = s
	return Round1Time
}

func (r *Room) setRound2Running() time.Duration {
	s := round2Running{
		r:          r,
		questionId: 4,
		question:   &tr.Questions[4],
	}
	r.state = s
	return Round2Time
}

func (r *Room) setRound3Running() time.Duration {
	s := round3Running{
		r:          r,
		questionId: 5,
		question:   &tr.Questions[5],
	}
	r.state = s
	return Round3Time
}

func (r *Room) setRound4Running() time.Duration {
	s := round4Running{
		r:          r,
		questionId: 7,
		question:   &tr.Questions[7],
	}
	r.state = s
	return Round4Time
}

// registers client into lobby, returns whether room is full after register
func (r *Room) registerClient(c *Client) bool {
	r.broadcast(m.NewClientJoinedMessage(c.playerInfo()))

	c.roomRead = r.roomRead
	go c.readPump()

	c.send(m.NewRoomGreetingMessage(r.id, r.playersInfo()))
	r.log("Greeting Message Written")

	c.done = false
	c.closed = false

	r.clients[c.id] = c

	r.activeClientCount++
	full := r.activeClientCount == ClientsPerRoom

	return full
}

func (r *Room) unregisterClient(c *Client) {
	if c.closed {
		return
	}

	c.closed = true
	r.activeClientCount--

	c.log("unregistered")

	r.broadcast(m.NewClientLeftMessage(c.playerInfo()))

	if r.activeClientCount == 0 {
		r.log("all clients left")
		r.end()
	}
}

func (r *Room) broadcast(msg m.ServerMessage) {
	for _, c := range r.clients {
		if !c.closed {
			c.send(msg)
		}
	}
}

func (r *Room) handleEliminations() {
	currentPlayers := []m.PlayerInfo{}
	eliminatedPlayers := []m.PlayerInfo{}

	for _, c := range r.clients {
		if !c.done && !c.closed {
			c.send(m.NewClientEliminatedMessage(
				c.playerInfo(),
				r.activeClientCount,
				len(r.clients),
			))
			r.unregisterClient(c)
		}

		pi := c.playerInfo()
		if !c.done || c.closed {
			eliminatedPlayers = append(eliminatedPlayers, pi)
		} else {
			currentPlayers = append(currentPlayers, pi)
		}
	}

	r.broadcast(m.NewRoundEndMessage(
		0,
		currentPlayers,
		eliminatedPlayers,
	))

	if r.activeClientCount == 0 {
		r.end()
	}
}

func (r *Room) everyoneDone() bool {
	for _, c := range r.clients {
		if !c.closed && !c.done {
			return false
		}
	}
	return true
}

func (r *Room) playersInfo() []m.PlayerInfo {
	playersInfo := make([]m.PlayerInfo, 0, len(r.clients))
	for _, client := range r.clients {
		playersInfo = append(playersInfo, client.playerInfo())
	}
	return playersInfo
}

func (r *Room) isOpen() bool {
	if len(r.clients) == ClientsPerRoom {
		return false
	}
	_, ok := r.state.(waitingForPlayers)
	return ok
}

func (r *Room) log(format string, v ...any) {
	log.Printf("room %d: %s", r.id, fmt.Sprintf(format, v...))
}
