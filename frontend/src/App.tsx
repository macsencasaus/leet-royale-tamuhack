import { useState } from "react";
import Workspace from "./components/workspace";
import { WebSocketContext } from "./hooks/context";
import Login from "./components/login";
// import Lobby from "./components/lobby";

function App() {
	const [connected, setConnected] = useState(false);
	const [webSocket, setWebSocket] = useState<WebSocket | undefined>(
		undefined
	);

	return (
		<WebSocketContext.Provider
			value={{
				connected,
				setConnected,
				webSocket,
				setWebSocket,
			}}
		>
			<div className="flex flex-col w-dvh h-dvh p-2 gap-2">
				{!connected && <Login />}
				{/* <Lobby /> */}
				{connected && <Workspace />}
				{/* <Workspace /> */}
			</div>
		</WebSocketContext.Provider>
	);
}

export default App;
