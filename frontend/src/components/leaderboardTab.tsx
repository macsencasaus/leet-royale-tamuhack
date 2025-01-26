import useWebSocket from "@/hooks/useWebSocket";
import { Message, PlayerStatus } from "@/lib/types";
import { Award, Crown } from "lucide-react";
import { useCallback, useState } from "react";

// const temp: Record<number, PlayerStatus> = {
// 	0: {
// 		player: {
// 			id: 0,
// 			name: "Macsen",
// 		},
// 		finished: true,
// 		casesCompleted: 50,
// 		timestamp: 10,
// 	},
// 	1: {
// 		player: {
// 			id: 1,
// 			name: "Owen",
// 		},
// 		finished: false,
// 		casesCompleted: 10,
// 		timestamp: 20,
// 	},
// 	2: {
// 		player: {
// 			id: 2,
// 			name: "Alex",
// 		},
// 		finished: false,
// 		casesCompleted: 0,
// 		timestamp: 1,
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
					timestamp: message.timestamp,
				};

				setPlayers((prev) => ({
					...prev,
					...next,
				}));

				break;
		}
	}, []);

	useWebSocket(onMessage);

	function orderedPlayers(
		players: Record<number, PlayerStatus>
	): PlayerStatus[] {
		return Object.values(players).sort((a, b) => {
			if (a.finished && b.finished) {
				return a.timestamp - b.timestamp;
			} else if (a.finished) {
				return -1;
			} else if (b.finished) {
				return 1;
			} else if (b.casesCompleted !== a.casesCompleted) {
				return b.casesCompleted - a.casesCompleted;
			} else if (b.timestamp !== a.timestamp) {
				return a.timestamp - b.timestamp;
			} else {
				return a.player.name.localeCompare(b.player.name);
			}
		});
	}

	return Object.keys(players).length > 0 ? (
		<div className="flex flex-col gap-2">
			{orderedPlayers(players).map((value, index) => (
				<div
					className="border border-white/10 rounded p-2 flex justify-between"
					style={{
						backgroundColor: value.finished
							? "rgba(25, 135, 84, .5)"
							: undefined,
					}}
				>
					<p className="font-bold flex gap-2">
						{value.finished &&
							(index === 0 ? (
								<Crown />
							) : index === 1 || index === 2 ? (
								<Award />
							) : undefined)}
						{value.player.name}
					</p>
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
