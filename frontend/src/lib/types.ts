export type GameState = "login" | "lobby" | "workspace";

export type Languages = "python" | "javascript" | "cpp";
export type MessageType =
	| "ClientMessageReady"
	| "ClientMessageClientQuit"
	| "ClientMessageSubmit"
	| "ServerMessageHubGreeting"
	| "ServerMessageRoomGreeting"
	| "ServerMessageCountdown"
	| "ServerMessageClientJoined"
	| "ServerMessageClientLeft"
	| "ServerMessageRoundStart"
	| "ServerMessageRoundEnd"
	| "ServerMessageTestPassed"
	| "ServerMessageTestFailed";

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
    question: string;
    tle: boolean;
    cases: {}[];
}

export interface Player {
	id: number;
	name: string;
}

export interface TestCases {
    success: boolean;
    
}