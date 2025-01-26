export type GameState = "login" | "lobby" | "workspace";

export type Languages = "python" | "javascript" | "cpp";
export type MessageType =
	| "ClientMessageReady"
	| "ClientMessageClientQuit"
	| "ClientMessageSubmit"
	| "ClientMessageSkipLobby"
	| "ClientMessageSkipQuestion"
	| "ClientMessageBuyItem"
	| "ServerMessageHubGreeting"
	| "ServerMessageRoomGreeting"
	| "ServerMessageCountdown"
	| "ServerMessageClientJoined"
	| "ServerMessageClientLeft"
	| "ServerMessageRoundStart"
	| "ServerMessageRoundEnd"
	| "ServerMessageTestResult"
	| "ServerMessageUpdateClientStatus"
	| "ServerMessageClientEliminated";

export interface Message {
	type: MessageType;

	player: Player;
	lobbyId: number;
	otherPlayers: Player[];
	count: number;

	round: number;
	time: number; // seconds
	prompt: string;
	templates: {
		python: string;
		javascript: string;
		cpp: string;
	};
	numTestCases: number;
	visibleTestCases: { input: string; output: string }[];
	tle: boolean;
	cases: TestCase[];
	finished: boolean;
	casesCompleted: number;
	place: number;
	totalPlayers: number;
}

export interface Player {
	id: number;
	name: string;
}

export interface TestCase {
	success: boolean;
	stdout: string;
}

export interface Templates {
	python: string;
	javascript: string;
	cpp: string;
}

export interface PlayerStatus {
	player: Player;
	finished: boolean;
	casesCompleted: number;
}
