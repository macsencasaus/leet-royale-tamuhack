export type GameState = "login" | "lobby" | "workspace" | "eliminated" | "winner";

export type Languages = "python" | "javascript" | "cpp";
export type MessageType =
    | "ClientMessageReady"
    | "ClientMessageClientQuit"
    | "ClientMessageSubmit"
    | "ClientMessageSkipLobby"
    | "ClientMessageSkipQuestion"
    | "ClientMessageBuyItem"
    | "ServerMessageClientItem"
    | "ServerMessageHubGreeting"
    | "ServerMessageRoomGreeting"
    | "ServerMessageCountdown"
    | "ServerMessageClientJoined"
    | "ServerMessageClientLeft"
    | "ServerMessageRoundStart"
    | "ServerMessageRoundEnd"
    | "ServerMessageTestResult"
    | "ServerMessageUpdateClientStatus"
    | "ServerMessageClientEliminated"
    | "ServerMessageWinner";

export interface Message {
    type: MessageType;

    player: Player;
    lobbyId: number;
    otherPlayers: Player[];
    count: number;

    questionId: number;
    round: number;
    time: number; // seconds
    prompt: string;
    templates: {
        python: string;
        javascript: string;
        cpp: string;
    };
    numTestCases: number;
    visibleTestCases: VisibleCases[];
    tle: boolean;
    cases: TestCase[];
    finished: boolean;
    casesCompleted: number;
    place: number;
    totalPlayers: number;
    timestamp: number;
    item: Items;
    currentPlayers: Player[];
    eliminatedPlayers: Player[];
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
    timestamp: number;
}

export interface VisibleCases {
    input: string;
    output: string;
}

export type Items =
    | "light-mode"
    | "dvd-logo"
    | "ad-space"
    | "freeze"
    | "remove-all"
    | "remove-line"
    | "arrow-only";
