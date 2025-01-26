import { GameState, Player } from "@/lib/types";
import { createContext } from "react";

interface WebSocketContextValue {
	gameState: GameState;
	setGameState: React.Dispatch<React.SetStateAction<GameState>>;
	connected: boolean;
	setConnected: React.Dispatch<React.SetStateAction<boolean>>;
	webSocket: WebSocket | undefined;
	setWebSocket: React.Dispatch<React.SetStateAction<WebSocket | undefined>>;
	player: Player | undefined;
	setPlayer: React.Dispatch<React.SetStateAction<Player | undefined>>;
}

export const WebSocketContext = createContext<WebSocketContextValue>({
	gameState: "login",
	setGameState: () => {},
	connected: false,
	setConnected: () => {},
	webSocket: undefined,
	setWebSocket: () => {},
	player: undefined,
	setPlayer: () => {},
});
