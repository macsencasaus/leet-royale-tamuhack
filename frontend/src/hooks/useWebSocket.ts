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
        player,
        setPlayer,
    } = useContext(WebSocketContext);

    useEffect(() => {
        if (onMessage && webSocket && connected) {
            const handleMessage = (ev: MessageEvent) => {
                const ms: Message = JSON.parse(ev.data);

                if (ms.type === "ServerMessageHubGreeting") {
                    setPlayer(ms.player);
                }

                if (gameState == "lobby" && ms.type === "ServerMessageRoundStart") {
                    setGameState("workspace");
                }

                if (ms.type === "ServerMessageClientEliminated") {
                    setGameState("eliminated");
                }

                onMessage(ms);
            };

            webSocket.addEventListener("message", handleMessage);
            return () => {
                webSocket.removeEventListener("message", handleMessage);
            };
        }
    }, [onMessage, webSocket, connected]);

    function createWebSocket(address: string) {
        const ws = new WebSocket(address);

        ws.onopen = () => {
            setConnected(true);
            setGameState("lobby");
        };

        ws.onclose = () => {
            setConnected(false);
            setGameState((prev) => (prev === "eliminated" ? "eliminated" : "login"));
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
        setGameState,
        player,
        connected,
        webSocket,
        sendMessage,
        createWebSocket,
    };
}

export default useWebSocket;
