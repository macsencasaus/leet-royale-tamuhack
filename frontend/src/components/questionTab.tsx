import useWebSocket from "@/hooks/useWebSocket";
import { Message } from "@/lib/types";
import { useState } from "react";

function QuestionTab() {
	useWebSocket(onMessage);

	const [question, setQuestion] = useState("");

	function onMessage(message: Message) {
		switch (message.type) {
			case "ServerMessageRoundStart":
				setQuestion(message.question);
				break;
		}
	}

	return <p>{question}</p>;
}

export default QuestionTab;
