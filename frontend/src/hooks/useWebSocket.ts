import { useContext, useEffect } from "react";
import { WebSocketContext } from "./context";
import { MessageType } from "@/lib/types";

function useWebSocket(onMessage?: (message: any) => void) {
	const { connected, setConnected, webSocket, setWebSocket } =
		useContext(WebSocketContext);

	useEffect(() => {
		if (onMessage && webSocket && connected) {
			console.trace("Listening for messages");

			webSocket.addEventListener("message", onMessage);
			return () => {
				webSocket.removeEventListener("message", onMessage);
			};
		}
	}, [onMessage, webSocket, connected]);

	function createWebSocket(address: string) {
		const ws = new WebSocket(address);

		ws.onopen = () => {
			setConnected(true);
		};

		ws.onclose = () => {
			setConnected(false);
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
		connected,
		webSocket,
		sendMessage,
		createWebSocket,
	};
}

export default useWebSocket;
