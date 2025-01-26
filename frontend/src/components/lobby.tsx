import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
// import { debug } from "@/config";
import useWebSocket from "@/hooks/useWebSocket";
import { Message, Player } from "@/lib/types";
import { useEffect, useState } from "react";
import { ScrollArea } from "./ui/scroll-area";

function Lobby() {
	const { gameState, sendMessage } = useWebSocket(onMessage);
	const [waitTime, setWaitTime] = useState<number | undefined>(undefined);
	const [others, setOthers] = useState<Player[]>([]);

	useEffect(() => {
		sendMessage("ClientMessageReady");
	}, []);

	function onMessage(message: Message) {
		switch (message.type) {
			case "ServerMessageRoomGreeting":
				setOthers(message.otherPlayers);
				break;
			case "ServerMessageClientJoined":
				setOthers((prev) => [message.player, ...prev]);
				break;
			case "ServerMessageClientLeft":
				setOthers((prev) =>
					prev.filter((player) => player.id != message.player.id)
				);
				break;
			case "ServerMessageCountdown":
				setWaitTime(message.count);
				break;
		}
	}

	return (
		<div
			className={`flex justify-center h-full items-center ${
				gameState !== "lobby" ? "hidden" : ""
			}`}
		>
			<Card className="w-screen max-w-xs">
				<CardHeader>
					<CardTitle>Lobby</CardTitle>
					<CardDescription>
						{waitTime
							? `Starting game in ${waitTime}`
							: "Waiting to begin countdown..."}
					</CardDescription>
				</CardHeader>
				<CardContent>
					<ScrollArea className="h-[300px]">
						{others.length > 0 ? (
							<ol
								className="list-decimal list-inside"
								reversed
							>
								{others.map((player) => (
									<li>{player.name}</li>
								))}
							</ol>
						) : (
							<p>Nobody else has joined yet.</p>
						)}
					</ScrollArea>
				</CardContent>
			</Card>
		</div>
	);
}

export default Lobby;
