import { useContext, useEffect } from "react";
import { WebSocketContext } from "./context";

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

	return {
		connected,
		webSocket,
		createWebSocket,
	};
}

export default useWebSocket;
