import useWebSocket from "@/hooks/useWebSocket";
import { Message } from "@/lib/types";
import { useCallback, useEffect, useState } from "react";

function HeaderPanel() {
	const [round, setRound] = useState(0);
	const [timeLeft, setTimeLeft] = useState<number | undefined>(100);

	const onMessage = useCallback((message: Message) => {
		switch (message.type) {
			case "ServerMessageRoundStart":
				setRound(message.round);
				setTimeLeft(message.time);
				break;
			case "ServerMessageRoundEnd":
				setTimeLeft(undefined);
				break;
		}
	}, []);

	useWebSocket(onMessage);

	useEffect(() => {
		if (timeLeft === undefined) {
			return;
		}

		const timer = setTimeout(() => {
			setTimeLeft((prev) => (prev !== undefined ? prev - 1 : undefined));
		}, 1000);

		return () => clearTimeout(timer);
	}, [timeLeft]);

	function formatTime(seconds: number | undefined) {
		if (seconds === undefined) {
			return "";
		} else if (seconds < 0) {
			return "Last Seconds!";
		}

		const minutes = Math.floor(seconds / 60);
		const remainingSeconds = seconds % 60;
		const paddedSeconds = String(remainingSeconds).padStart(2, "0");

		return (
			<p>
				Time Left:{" "}
				<span className="font-bold">
					{minutes}:{paddedSeconds}
				</span>
			</p>
		);
	}

	return (
		<div className="border border-white/10 bg-border rounded p-2 px-4 flex justify-between">
			<h1 className="font-bold">Leet Royale</h1>
			<p className="font-bold">Round {round}</p>
			{formatTime(timeLeft)}
		</div>
	);
}

export default HeaderPanel;
