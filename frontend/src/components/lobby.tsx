import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import useWebSocket from "@/hooks/useWebSocket";
import { Message, Player } from "@/lib/types";
import { useCallback, useEffect, useState } from "react";
import { ScrollArea } from "./ui/scroll-area";
import { Button } from "./ui/button";
import { SkipForward } from "lucide-react";
import Rules from "./rules";
import { debug } from "@/config";

function Lobby({ force }: { force?: boolean }) {
	const [lobbyId, setLobbyId] = useState(-1);
	const [waitTime, setWaitTime] = useState<number | undefined>(undefined);
	const [others, setOthers] = useState<Player[]>([]);

	useEffect(() => {
		sendMessage("ClientMessageReady");
	}, []);

	const onMessage = useCallback((message: Message) => {
		switch (message.type) {
			case "ServerMessageRoomGreeting":
				setLobbyId(message.lobbyId);
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
	}, []);

	const { gameState, player, sendMessage } = useWebSocket(onMessage);

	function skipLobby() {
		sendMessage("ClientMessageSkipLobby");
	}

	return (
		<div
			className={`flex gap-2 justify-center h-full items-center ${
				gameState !== "lobby" && !force ? "hidden" : ""
			}`}
		>
			<Card className="w-screen max-w-xs">
				<CardHeader>
					<CardTitle>Lobby {lobbyId != -1 ? lobbyId : ""}</CardTitle>
					{player && (
						<CardDescription>
							{others.length > 0 ? (
								<p>
									Settle in,{" "}
									<span className="font-bold">
										{player.name}
									</span>
									. This is your fierce competition:
								</p>
							) : (
								<p>
									Welcome,{" "}
									<span className="font-bold">
										{player.name}
									</span>
									.
								</p>
							)}
						</CardDescription>
					)}
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
				<CardFooter>
					<div className="flex flex-col w-full gap-2">
						{waitTime
							? `Starting game in ${waitTime}`
							: "Waiting to begin countdown..."}
                        {debug && (
                            <Button
                                className="self-end"
                                variant={"outline"}
                                onClick={skipLobby}
                            >
                                <SkipForward />
                                Skip Wait
                            </Button>
                        )}
					</div>
				</CardFooter>
			</Card>
			<Rules />
		</div>
	);
}

export default Lobby;
