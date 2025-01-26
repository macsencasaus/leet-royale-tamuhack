package messages

import (
	"encoding/json"
)

type ClientMessageType string

const (
	ClientMessageTypeReady      = "ClientMessageReady"
	ClientMessageTypeClientQuit = "ClientMessageClientQuit"
	ClientMessageTypeSubmit     = "ClientMessageSubmit"

	ClientMessageTypeSkipLobby    = "ClientMessageSkipLobby"
	ClientMessageTypeSkipQuestion = "ClientMessageSkipQuestion"
)

type ClientMessage interface {
	clientMessage()
}

type ClientMessageWrapper struct {
	Type ClientMessageType `json:"type"`
	Data json.RawMessage   `json:"data"`
}

type ReadyMessage struct {
	Type ClientMessageType `json:"type"`
}

func (m ReadyMessage) clientMessage() {}

type ClientQuitMessage struct {
	PlayerId int `json:"playerId"`
}

func (m ClientQuitMessage) clientMessage() {}

type SubmitMessage struct {
	PlayerId int    `json:"playerId"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

func (m SubmitMessage) clientMessage() {}

type SkipLobbyMessage struct{}

func (m SkipLobbyMessage) clientMessage() {}

type SkipQuestionMessage struct{}

func (m SkipQuestionMessage) clientMessage() {}
