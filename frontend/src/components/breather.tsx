import useWebSocket from "@/hooks/useWebSocket";
import { Message } from "@/lib/types";
import { useCallback } from "react";

function Breather() {
	const onMessage = useCallback((message: Message) => {
		switch (message.type) {
			case "ServerMessageClientEliminated":
				break;
		}
	}, []);

	useWebSocket(onMessage);

	return (
		<div className="border border-white/10 p-2 rounded h-full">
			<h1>Take a breather...</h1>
		</div>
	);
}

export default Breather;
