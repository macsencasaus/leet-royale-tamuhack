import { createContext } from "react";

interface WebSocketContextValue {
	connected: boolean;
	setConnected: React.Dispatch<React.SetStateAction<boolean>>;
	webSocket: WebSocket | undefined;
	setWebSocket: React.Dispatch<React.SetStateAction<WebSocket | undefined>>;
}

export const WebSocketContext = createContext<WebSocketContextValue>({
	connected: false,
	setConnected: () => {}, // Default no-op function
	webSocket: undefined,
	setWebSocket: () => {}, // Default no-op function
});