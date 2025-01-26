import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Button } from "./ui/button";
import useWebSocket from "@/hooks/useWebSocket";
import { Message } from "@/lib/types";

function Eliminated() {
	const { gameState, setGameState } = useWebSocket(onMessage);

	function onMessage(message: Message) {
		switch (message.type) {
			case "ServerMessageClientEliminated":
				setGameState("eliminated");
		}
	}

	function onClick() {
		setGameState("login");
	}

	return (
		<div
			className="flex justify-center h-dvh items-center"
			style={{ display: gameState !== "eliminated" ? "none" : "" }}
		>
			<Card>
				<CardHeader>
					<CardTitle>Eliminated</CardTitle>
					<CardDescription>You've been eliminated.</CardDescription>
				</CardHeader>
				<CardContent>
					<p>Better 1337 next time.</p>
				</CardContent>
				<CardFooter className="flex justify-end">
					<Button
						variant={"outline"}
						onClick={onClick}
					>
						Accept Defeat
					</Button>
				</CardFooter>
			</Card>
		</div>
	);
}

export default Eliminated;
