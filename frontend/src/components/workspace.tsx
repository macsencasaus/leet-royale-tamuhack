import {
	ResizableHandle,
	ResizablePanel,
	ResizablePanelGroup,
} from "@/components/ui/resizeable";
import EditorPanel from "./editorPanel";
import InfoPanel from "./infoPanel";
import HeaderPanel from "./headerPanel";
import useWebSocket from "@/hooks/useWebSocket";
// import { debug } from "@/config";

function Workspace() {
	const { gameState } = useWebSocket();

	return (
		<div
			className={`flex flex-col gap-1 h-full ${
				gameState !== "workspace" ? "hidden" : ""
			}`}
		>
			<HeaderPanel />
			<ResizablePanelGroup
				direction="horizontal"
				className="grow"
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
