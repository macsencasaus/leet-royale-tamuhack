import { useState } from "react";
import Workspace from "./components/workspace";
import { WebSocketContext } from "./hooks/context";
import Login from "./components/login";
import Lobby from "./components/lobby";
import { GameState, Player } from "./lib/types";

function App() {
	const [gameState, setGameState] = useState<GameState>("login");
	const [connected, setConnected] = useState(false);
	const [webSocket, setWebSocket] = useState<WebSocket | undefined>(
		undefined
	);
	const [player, setPlayer] = useState<Player | undefined>(undefined);

	return (
		<WebSocketContext.Provider
			value={{
				gameState,
				setGameState,
				connected,
				setConnected,
				webSocket,
				setWebSocket,
				player,
				setPlayer,
			}}
		>
			<div className="flex flex-col w-dvh h-dvh p-2 gap-2">
				{!connected && <Login />}
				{connected && (
					<>
						<Lobby />
						<Workspace />
					</>
				)}

				<Lobby />
				{/* <Workspace /> */}
			</div>
		</WebSocketContext.Provider>
	);
}

export default App;
