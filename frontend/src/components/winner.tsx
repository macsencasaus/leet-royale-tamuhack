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

function Winner() {
	const { gameState, setGameState } = useWebSocket(onMessage);

	function onMessage(message: Message) {
		switch (message.type) {
			case "ServerMessageWinner":
				setGameState("winner");
		}
	}

	function onClick() {
		setGameState("login");
	}

	return (
		<div
			className="flex justify-center h-dvh items-center"
			style={{ display: gameState !== "winner" ? "none" : "" }}
		>
			<Card>
				<CardHeader>
					<CardTitle>Winner</CardTitle>
					<CardDescription>You won!</CardDescription>
				</CardHeader>
				<CardContent>
					<p>Nice victory royale! :swag:</p>
				</CardContent>
				<CardFooter className="flex justify-end">
					<Button
						variant={"outline"}
						onClick={onClick}
					>
						Victory Dance
					</Button>
				</CardFooter>
			</Card>
		</div>
	);
}

export default Winner;
