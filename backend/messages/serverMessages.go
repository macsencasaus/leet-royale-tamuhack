package messages

import tr "leet-guys/testrunner"

type ServerMessageType string

const (
	ServerMessageTypeHubGreeting    = "ServerMessageHubGreeting"
	ServerMessageTypeRoomGreeting   = "ServerMessageRoomGreeting"
	ServerMessageCountdown          = "ServerMessageCountdown"
	ServerMessageClientJoined       = "ServerMessageClientJoined"
	ServerMessageClientLeft         = "ServerMessageClientLeft"
	ServerMessageRoundStart         = "ServerMessageRoundStart"
	ServerMessageRoundEnd           = "ServerMessageRoundEnd"
	ServerMessageTestResult         = "ServerMessageTestResult"
	ServerMessageUpdateClientStatus = "ServerMessageUpdateClientStatus"
	ServerMessageClientEliminated   = "ServerMessageClientEliminated"
)

type ServerMessage interface {
	serverMessage()
}

type HubGreetingMessage struct {
	Type   ServerMessageType `json:"type"`
	Player PlayerInfo        `json:"player"`
}

func NewHubGreetingMessage(p PlayerInfo) HubGreetingMessage {
	return HubGreetingMessage{
		Type:   ServerMessageTypeHubGreeting,
		Player: p,
	}
}
func (m HubGreetingMessage) serverMessage() {}

type RoomGreetingMessage struct {
	Type         ServerMessageType `json:"type"`
	LobbyId      int               `json:"lobbyId"`
	OtherPlayers []PlayerInfo      `json:"otherPlayers"`
}

func NewRoomGreetingMessage(lobbyId int, otherPlayers []PlayerInfo) RoomGreetingMessage {
	return RoomGreetingMessage{
		Type:         ServerMessageTypeRoomGreeting,
		LobbyId:      lobbyId,
		OtherPlayers: otherPlayers,
	}
}
func (m RoomGreetingMessage) serverMessage() {}

type CountdownMessage struct {
	Type  ServerMessageType `json:"type"`
	Count int               `json:"count"`
}

func NewCountdownMessage(count int) CountdownMessage {
	return CountdownMessage{
		Type:  ServerMessageCountdown,
		Count: count,
	}
}
func (m CountdownMessage) serverMessage() {}

type ClientJoinedMessage struct {
	Type   ServerMessageType `json:"type"`
	Player PlayerInfo        `json:"player"`
}

func NewClientJoinedMessage(player PlayerInfo) ClientJoinedMessage {
	return ClientJoinedMessage{
		Type:   ServerMessageClientJoined,
		Player: player,
	}
}
func (m ClientJoinedMessage) serverMessage() {}

type ClientLeftMessage struct {
	Type   ServerMessageType `json:"type"`
	Player PlayerInfo        `json:"player"`
}

func NewClientLeftMessage(player PlayerInfo) ClientLeftMessage {
	return ClientLeftMessage{
		Type:   ServerMessageClientLeft,
		Player: player,
	}
}
func (m ClientLeftMessage) serverMessage() {}

type RoundStartMessage struct {
	Type             ServerMessageType            `json:"type"`
	Round            int                          `json:"round"`
	Time             int                          `json:"time"`
	Prompt           string                       `json:"prompt"`
	Templates        tr.LanguageFunctionTemplates `json:"templates"`
	NumTestCases     int                          `json:"numTestCases"`
	VisibleTestCases []TestCase                   `json:"visibleTestCases"`
}

func NewRoundStartMessage(
	round int,
	sec int,
	rd *tr.QuestionData,
) RoundStartMessage {
	return RoundStartMessage{
		Type:         ServerMessageRoundStart,
		Round:        round,
		Time:         sec,
		Prompt:       rd.Prompt,
		Templates:    rd.Templates,
		NumTestCases: rd.NumCases,
		// TODO:
		VisibleTestCases: []TestCase{
			{
				Input:  "1, 2",
				Output: "3",
			},
			{
				Input:  "-1, 1",
				Output: "0",
			},
			{
				Input:  "77, 33",
				Output: "110",
			},
		},
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

type TestResultMessage struct {
	Type  ServerMessageType `json:"type"`
	TLE   bool              `json:"tle"`
	Cases []ResultCase      `json:"cases"`
}

// TODO:
func NewTestResultMessage(res *tr.Result) TestResultMessage {
	var tle bool

	tle = res.Issue == tr.RunTime

	resCases := make([]ResultCase, 0, res.NCasesRun)

	for i, status := range res.PFStatus {
		resCases = append(resCases, ResultCase{
			Success: status == tr.AC,
			Stdout:  string(res.Stdout[i]),
		})
	}

	return TestResultMessage{
		Type:  ServerMessageTestResult,
		TLE:   tle,
		Cases: resCases,
	}
}
func (m TestResultMessage) serverMessage() {}

type UpdateClientStatus struct {
	Type           ServerMessageType `json:"type"`
	Player         PlayerInfo        `json:"player"`
	Finished       bool              `json:"finished"`
	CasesCompleted int               `json:"casesCompleted"`
	Timestamp      int               `json:"timestamp"`
}

func NewUpdateClientStateMessage(
	player PlayerInfo,
	finished bool,
	casesCompleted int,
    timestamp int,
) UpdateClientStatus {
	return UpdateClientStatus{
		Type:           ServerMessageUpdateClientStatus,
		Player:         player,
		Finished:       finished,
		CasesCompleted: casesCompleted,
        Timestamp: timestamp,
	}
}
func (m UpdateClientStatus) serverMessage() {}

type ClientEliminatedMessage struct {
	Type         ServerMessageType `json:"type"`
	Player       PlayerInfo        `json:"player"`
	Place        int               `json:"place"`
	TotalPlayers int               `json:"totalPlayers"`
}

func NewClientEliminatedMessage(
	player PlayerInfo,
	place int,
	totalPlayers int,
) ClientEliminatedMessage {
	return ClientEliminatedMessage{
		Type:         ServerMessageClientEliminated,
		Player:       player,
		Place:        place,
		TotalPlayers: totalPlayers,
	}
}

func (m ClientEliminatedMessage) serverMessage() {}

type PlayerInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type LanguageFunctionTemplate struct {
	Python     string `json:"python"`
	Javascript string `json:"javascript"`
	Cpp        string `json:"cpp"`
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type ResultCase struct {
	Success bool   `json:"success"`
	Stdout  string `json:"stdout"`
}
