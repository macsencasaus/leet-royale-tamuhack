import {
	ResizableHandle,
	ResizablePanel,
	ResizablePanelGroup,
} from "@/components/ui/resizeable";
import EditorPanel from "./editorPanel";
import InfoPanel from "./infoPanel";
import HeaderPanel from "./headerPanel";
import useWebSocket from "@/hooks/useWebSocket";
import { useCallback, useState } from "react";
import { Message } from "@/lib/types";
import Breather from "./breather";

function Workspace({ force }: { force?: boolean }) {
	const [waiting, setWaiting] = useState(false);

	const onMessage = useCallback((message: Message) => {
		switch (message.type) {
			case "ServerMessageRoundStart":
				setWaiting(false);
				break;
			case "ServerMessageRoundEnd":
				setWaiting(true);
				break;
		}
	}, []);

	const { gameState } = useWebSocket(onMessage);

	return (
		<div
			className={`flex flex-col gap-1 h-full ${
				gameState !== "workspace" && !force ? "hidden" : ""
			}`}
		>
			<HeaderPanel />
			<div className="grow" style={{ display: waiting ? "" : "none" }}>
				<Breather />
			</div>
			<ResizablePanelGroup
				direction="horizontal"
				className="grow"
				style={{ display: waiting ? "none" : "" }}
			>
				<ResizablePanel className="border border-white/10 rounded mr-1 p-2 h-full w-full">
					<InfoPanel />
				</ResizablePanel>
				<ResizableHandle className="bg-transparent" />
				<ResizablePanel className="border border-white/10 rounded p-2 h-full">
					<EditorPanel />
				</ResizablePanel>
			</ResizablePanelGroup>
		</div>
	);
}

export default Workspace;
