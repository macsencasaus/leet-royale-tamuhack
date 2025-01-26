import { useContext, useEffect } from "react";
import { WebSocketContext } from "./context";
import { Message, MessageType } from "@/lib/types";

function useWebSocket(onMessage?: (message: any) => void) {
	const {
		gameState,
		setGameState,
		connected,
		setConnected,
		webSocket,
		setWebSocket,
		setPlayer,
	} = useContext(WebSocketContext);

	useEffect(() => {
		if (onMessage && webSocket && connected) {
			webSocket.addEventListener("message", _onMessage);
			return () => {
				webSocket.removeEventListener("message", _onMessage);
			};
		}
	}, [onMessage, webSocket, connected]);

	function _onMessage(ev: MessageEvent) {
		const ms: Message = JSON.parse(ev.data);
		console.log(ms);

		if (ms.type === "ServerMessageHubGreeting") {
			setPlayer(ms.player);
		}

		if (gameState == "lobby" && ms.type === "ServerMessageRoundStart") {
			setGameState("workspace");
		}

		onMessage!(ms);
	}

	function createWebSocket(address: string) {
		const ws = new WebSocket(address);

		ws.onopen = () => {
			setConnected(true);
			setGameState("lobby");
		};

		ws.onclose = () => {
			setConnected(false);
			setGameState("login");
		};

		setWebSocket(ws);
	}

	function sendMessage(type: MessageType, data?: Record<string, any>) {
		const message = {
			type,
			data,
		};

		webSocket?.send(JSON.stringify(message));
	}

	return {
		gameState,
		connected,
		webSocket,
		sendMessage,
		createWebSocket,
	};
}

export default useWebSocket;
