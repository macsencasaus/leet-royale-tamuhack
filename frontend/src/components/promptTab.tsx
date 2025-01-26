import useWebSocket from "@/hooks/useWebSocket";
import { Message } from "@/lib/types";
import { useCallback, useState } from "react";

function PromptTab() {
	const [prompt, setPrompt] = useState("No prompt yet.");

	const onMessage = useCallback((message: Message) => {
		switch (message.type) {
			case "ServerMessageRoundStart":
				setPrompt(message.prompt);
				break;
		}
	}, []);

	useWebSocket(onMessage);

	return <p dangerouslySetInnerHTML={{ __html: prompt }}></p>;
}

export default PromptTab;
