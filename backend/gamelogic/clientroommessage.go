package gamelogic

import tr "leet-guys/testrunner"

type ClientRoomMessage interface {
	clientRoomMessage()
}

type Quit struct {
	playerId int
}

func (Quit) clientRoomMessage() {}

type Submit struct {
	playerId   int
	questionId int
	results    tr.Result
}

func (Submit) clientRoomMessage() {}

type Skip struct {
}

func (Skip) clientRoomMessage() {}
