package messages

type ServerMessageType string

const (
	ServerMessageTypeHubGreeting   = "ServerMessageHubGreeting"
	ServerMessageTypeLobbyGreeting = "ServerMessageLobbyGreeting"
	ServerMessageCountdown         = "ServerMessageCountdown"
	ServerMessageRoundStart        = "ServerMessageRoundStart"
	ServerMessageRoundEnd          = "ServerMessageRoundEnd"
	ServerMessageTestPassed        = "ServerMessageTestPassed"
	ServerMessageTestFailed        = "ServerMessageTestFailed"
)

type ServerMessage interface {
	serverMessage()
}

type HubGreetingMessage struct {
	Type ServerMessageType `json:"type"`
}

func NewHubGreetingMessage() HubGreetingMessage {
	return HubGreetingMessage{
		Type: ServerMessageTypeHubGreeting,
	}
}
func (m HubGreetingMessage) serverMessage() {}

type LobbyGreetingMessage struct {
	Type    ServerMessageType `json:"type"`
	LobbyId int               `json:"lobbyId"`
	// TODO: add other players
}

func NewLobbyGreetingMessage(lobbyId int) LobbyGreetingMessage {
	return LobbyGreetingMessage{
		Type:    ServerMessageTypeLobbyGreeting,
		LobbyId: lobbyId,
	}
}
func (m LobbyGreetingMessage) serverMessage() {}

type CountdownMessage struct {
	Type  ServerMessageType `json:"type"`
	Count int               `json:"count"`
}

func NewMessageCountdownMessage(count int) CountdownMessage {
	return CountdownMessage{
		Type:  ServerMessageCountdown,
		Count: count,
	}
}
func (m CountdownMessage) serverMessage() {}

type RoundStartMessage struct {
	Type  ServerMessageType `json:"type"`
	Round int               `json:"round"`
}

func NewRoundStartMessage(round int) RoundStartMessage {
	return RoundStartMessage{
		Type:  ServerMessageRoundStart,
		Round: round,
	}
}
func (m RoundStartMessage) serverMessage() {}

type RoundEndMessage struct {
	Type  ServerMessageType `json:"type"`
	Round int               `json:"round"`
}

func NewRoundEndMessage(round int) RoundEndMessage {
	return RoundEndMessage{
		Type:  ServerMessageRoundEnd,
		Round: round,
	}
}
func (m RoundEndMessage) serverMessage() {}

type TestPassedMessage struct {
	Type     ServerMessageType `json:"type"`
	Question string            `json:"question"`
	// TODO: maybe add more info about length taken
}

func NewTestPassedMessage(question string) TestPassedMessage {
	return TestPassedMessage{
		Type:     ServerMessageTestPassed,
		Question: question,
	}
}
func (m TestPassedMessage) serverMessage() {}

type TestFailedMessage struct {
	Type     ServerMessageType `json:"type"`
	Question string            `json:"question"`
    // TODO: add failure reason
}

func NewTestFailedMessage(question string) TestFailedMessage {
	return TestFailedMessage{
		Type:     ServerMessageTestFailed,
		Question: question,
	}
}
func (m TestFailedMessage) servermessage() {}
