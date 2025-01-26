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

	return <div dangerouslySetInnerHTML={{ __html: prompt }}></div>;
}

export default PromptTab;
