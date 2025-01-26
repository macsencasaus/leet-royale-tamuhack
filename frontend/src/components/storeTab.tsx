import useWebSocket from "@/hooks/useWebSocket";
import { Items, Message } from "@/lib/types";
import { useCallback } from "react";
import { Button } from "./ui/button";
import { Disc3, ShieldAlert, Snowflake } from "lucide-react";

function StoreTab() {
	const onMessage = useCallback((message: Message) => {
		switch (message.type) {
			case "ServerMessageUpdateClientStatus":
				break;
		}
	}, []);

	const { sendMessage } = useWebSocket(onMessage);

	function buyItem(item: Items) {
		sendMessage("ClientMessageBuyItem", { item });
	}

	const items: Record<
		Items,
		{
			name: string;
			price: number;
			desc: string;
			icon: any;
		}
	> = {
		"ad-space": {
			name: "Ad Space",
			icon: <ShieldAlert />,
			price: 10,
			desc: "Play annoying ads on the enemy's screen.",
		},
		"dvd-logo": {
			name: "DvD Logo",
			icon: <Disc3 />,
			price: 10,
			desc: "Deploy a DvD logo to distract your opponent while they code.",
		},
		"light-mode": {
			name: "",
			price: 0,
			desc: "",
			icon: undefined,
		},
		freeze: {
			name: "",
			price: 0,
			desc: "",
			icon: <Snowflake />,
		},
		"remove-all": {
			name: "",
			price: 0,
			desc: "",
			icon: undefined,
		},
		"remove-line": {
			name: "",
			price: 0,
			desc: "",
			icon: undefined,
		},
		"arrow-only": {
			name: "",
			price: 0,
			desc: "",
			icon: undefined,
		},
	};

	return (
		<div className="flex flex-col gap-2 items-start">
			{Object.entries(items).map(([item, data]) => (
				<div>
					<Button
						variant={"outline"}
						className="flex gap-2"
						onClick={() => buyItem(item as Items)}
					>
						{data.icon}
						{data.name}
					</Button>
				</div>
			))}
		</div>
	);
}

export default StoreTab;
