import { useContext, useEffect } from "react";
import { WebSocketContext } from "./context";

function useWebSocket(onMessage?: (json: any) => void) {
	const { connected, setConnected, webSocket, setWebSocket } =
		useContext(WebSocketContext);

	useEffect(() => {
		if (onMessage && webSocket && connected) {
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
