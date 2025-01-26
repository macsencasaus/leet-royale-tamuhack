import useWebSocket from "@/hooks/useWebSocket";
import { Message, PlayerStatus } from "@/lib/types";
import { useCallback, useState } from "react";

// const _ = {
// 	0: {
// 		player: {
// 			id: 0,
// 			name: "Macsen",
// 		},
// 		finished: true,
// 		casesCompleted: 50,
// 	},
// 	1: {
// 		player: {
// 			id: 1,
// 			name: "Owen",
// 		},
// 		finished: false,
// 		casesCompleted: 10,
// 	},
// 	2: {
// 		player: {
// 			id: 2,
// 			name: "Alex",
// 		},
// 		finished: false,
// 		casesCompleted: 0,
// 	},
// };

function LeaderboardTab() {
	const [totalCases, setTotalCases] = useState(50);
	const [players, setPlayers] = useState<Record<number, PlayerStatus>>({});

	const onMessage = useCallback((message: Message) => {
		switch (message.type) {
			case "ServerMessageRoundStart":
				setTotalCases(message.numTestCases);
				break;
			case "ServerMessageUpdateClientStatus":
				const next: any = {};
				next[message.player.id] = {
					player: message.player,
					finished: message.finished,
					casesCompleted: message.casesCompleted,
				};

				setPlayers((prev) => ({
					...prev,
					...next,
				}));

				break;
		}
	}, []);

	useWebSocket(onMessage);

	function orderedPlayers(): PlayerStatus[] {
		return Object.values(
			Object.fromEntries(
				Object.entries(players).sort(([, a], [, b]) => {
					return b.casesCompleted - a.casesCompleted;
				})
			)
		);
	}

	return Object.keys(players).length > 0 ? (
		<div className="flex flex-col gap-2">
			{orderedPlayers().map((value) => (
				<div
					className="border border-white/10 rounded p-2 flex justify-between"
					style={{
						backgroundColor: value.finished
							? "rgba(25, 135, 84, .5)"
							: undefined,
					}}
				>
					<p className="font-bold">{value.player.name}</p>
					<p>
						<span className="font-bold">
							{value.casesCompleted}
						</span>
						/{totalCases}
					</p>
				</div>
			))}
		</div>
	) : (
		<p>Nobody has made any progress yet.</p>
	);
}

export default LeaderboardTab;
