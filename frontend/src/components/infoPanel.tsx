import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import useWebSocket from "@/hooks/useWebSocket";
import { useEffect, useState } from "react";

function InfoPanel() {
	const debug = true;
	const { sendMessage } = useWebSocket(onMessage);
	const [message, setMessage] = useState<string[]>([]);

	useEffect(() => {
		sendMessage("ClientMessageReady");
	});

	function onMessage(message: any) {
		console.log(message);
		setMessage(message);
	}

	return (
		<Tabs
			defaultValue={debug ? "debug" : "question"}
			className="flex flex-col"
		>
			<div className="bg-border overflow-hidden -m-2 p-1">
				<TabsList className="self-start">
					<TabsTrigger value="question">Question</TabsTrigger>
					<TabsTrigger value="leaderboard">Leaderboard</TabsTrigger>
					<TabsTrigger value="store">Store</TabsTrigger>
					<TabsTrigger value="submissions">Submissions</TabsTrigger>
					{debug && <TabsTrigger value="debug">Debug</TabsTrigger>}
				</TabsList>
			</div>

			<div className="p-2">
				<TabsContent value="question">Questions</TabsContent>
				<TabsContent value="leaderboard">Leaderboard</TabsContent>
				<TabsContent value="store">Store</TabsContent>
				<TabsContent value="submissions">Submissions</TabsContent>
				{debug && <TabsContent value="debug">{message}</TabsContent>}
			</div>
		</Tabs>
	);
}

export default InfoPanel;
