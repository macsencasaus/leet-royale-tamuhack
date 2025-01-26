import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { debug } from "@/config";
import useWebSocket from "@/hooks/useWebSocket";
import { Message } from "@/lib/types";
import { useCallback, useState } from "react";
import PromptTab from "./promptTab";

type Tab = "prompt" | "submissions" | "leaderboard" | "store";

function InfoPanel() {
	const [tab, setTab] = useState<Tab>("prompt");
	const [messages, setMessages] = useState<string[]>([]);

	const onMessage = useCallback((message: Message) => {
		const next = messages.concat([JSON.stringify(message)]);
		setMessages(next);
	}, []);

	useWebSocket(onMessage);

	return (
		<Tabs
			defaultValue={debug ? "prompt" : "prompt"}
			className="flex flex-col"
			onValueChange={(value) => setTab(value as Tab)}
		>
			<div className="bg-border overflow-hidden -m-2 p-1">
				<TabsList className="self-start">
					<TabsTrigger value="prompt">Prompt</TabsTrigger>
					<TabsTrigger value="submissions">Submissions</TabsTrigger>
					<TabsTrigger value="leaderboard">Leaderboard</TabsTrigger>
					<TabsTrigger value="store">Store</TabsTrigger>

					{debug && <TabsTrigger value="debug">Debug</TabsTrigger>}
				</TabsList>
			</div>

			<div className="p-2">
				<TabsContent
					value="prompt"
					forceMount
					style={{ display: tab !== "prompt" ? "none" : undefined }}
				>
					<PromptTab />
				</TabsContent>
				<TabsContent
					value="submissions"
					style={{ display: tab !== "submissions" ? "none" : undefined }}
				>
					Submissions
				</TabsContent>
				<TabsContent
					value="leaderboard"
					style={{ display: tab !== "leaderboard" ? "none" : undefined }}
				>
					Leaderboard
				</TabsContent>
				<TabsContent
					value="store"
					style={{ display: tab !== "store" ? "none" : undefined }}
				>
					Store
				</TabsContent>

				{debug && <TabsContent value="debug">{messages}</TabsContent>}
			</div>
		</Tabs>
	);
}

export default InfoPanel;
