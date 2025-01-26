import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { debug } from "@/config";
import useWebSocket from "@/hooks/useWebSocket";
import { Message } from "@/lib/types";
import { useState } from "react";
import QuestionTab from "./questionTab";

function InfoPanel() {
	useWebSocket(onMessage);
	const [messages, setMessages] = useState<string[]>([]);

	function onMessage(message: Message) {
		const next = messages.concat([JSON.stringify(message)]);
		setMessages(next);
	}

	return (
		<Tabs
			defaultValue={debug ? "debug" : "question"}
			className="flex flex-col"
		>
			<div className="bg-border overflow-hidden -m-2 p-1">
				<TabsList className="self-start">
					<TabsTrigger value="question">Question</TabsTrigger>
					<TabsTrigger value="submissions">Submissions</TabsTrigger>
					<TabsTrigger value="leaderboard">Leaderboard</TabsTrigger>
					<TabsTrigger value="store">Store</TabsTrigger>

					{debug && <TabsTrigger value="debug">Debug</TabsTrigger>}
				</TabsList>
			</div>

			<div className="p-2">
				<TabsContent value="question">
					<QuestionTab />
				</TabsContent>
				<TabsContent value="submissions">Submissions</TabsContent>
				<TabsContent value="leaderboard">Leaderboard</TabsContent>
				<TabsContent value="store">Store</TabsContent>

				{debug && <TabsContent value="debug">{messages}</TabsContent>}
			</div>
		</Tabs>
	);
}

export default InfoPanel;
