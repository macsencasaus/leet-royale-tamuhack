import {
	ResizableHandle,
	ResizablePanel,
	ResizablePanelGroup,
} from "@/components/ui/resizeable";
import Editor from "./editor";
import InfoPanel from "./infoPanel";

function Workspace() {
	return (
		<ResizablePanelGroup
			direction="horizontal"
			className="grow"
		>
			<ResizablePanel>
				<div className="border border-white/10 rounded mr-1 p-2 h-full w-full">
					<InfoPanel />
				</div>
			</ResizablePanel>
			<ResizableHandle className="bg-transparent" />
			<ResizablePanel>
				<div className="border border-white/10 rounded ml-1 p-2 h-full">
					<Editor />
				</div>
			</ResizablePanel>
		</ResizablePanelGroup>
	);
}

export default Workspace;
