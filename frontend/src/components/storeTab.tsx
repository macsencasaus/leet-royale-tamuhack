import useWebSocket from "@/hooks/useWebSocket";
import { Message } from "@/lib/types";
import { useCallback } from "react";
import { Button } from "./ui/button";
import { ShieldAlert } from "lucide-react";

function StoreTab() {
	const onMessage = useCallback((message: Message) => {
		switch (message.type) {
			case "ServerMessageUpdateClientStatus":
				break;
		}
	}, []);

	const { sendMessage } = useWebSocket(onMessage);

	function buyItem() {
		sendMessage("ClientMessageBuyItem", { item: "malware" });
	}

	return (
		<div>
			<Button
				variant={"outline"}
				onClick={buyItem}
			>
				<ShieldAlert />
				Malware
			</Button>
		</div>
	);
}

export default StoreTab;
