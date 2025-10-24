import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import useWebSocket from "@/hooks/useWebSocket";
import { Message } from "@/lib/types";
import { useCallback, useState } from "react";
import PromptTab from "./promptTab";
import SubmissionsTab from "./submissionsTab";
// import StoreTab from "./storeTab";
import { ScrollArea } from "./ui/scroll-area";
import LeaderboardTab from "./leaderboardTab";

type Tab = "prompt" | "submissions" | "leaderboard" | "store";

function InfoPanel() {
    const [tab, setTab] = useState<Tab>("prompt");

    const onMessage = useCallback((message: Message) => {
        switch (message.type) {
            case "ServerMessageTestResult":
                for (const c of message.cases) {
                    if (!c.success) {
                        setTab("submissions");
                        return;
                    }
                }

                setTab("leaderboard");
                break;
        }
    }, []);

    useWebSocket(onMessage);

    return (
        <Tabs
            defaultValue={"prompt"}
            className="flex flex-col h-full"
            value={tab}
            onValueChange={(value) => setTab(value as Tab)}
        >
            <div className="bg-border overflow-hidden -m-2 p-1 flex-none">
                <TabsList className="self-start h-min">
                    <TabsTrigger value="prompt">Prompt</TabsTrigger>
                    <TabsTrigger value="submissions">Submission</TabsTrigger>
                    <TabsTrigger value="leaderboard">Leaderboard</TabsTrigger>
                    {/* <TabsTrigger value="store">Store</TabsTrigger> */}
                </TabsList>
            </div>

            <div className="mt-4 overflow-hidden -mr-1.5">
                <ScrollArea className="h-full pr-3">
                    <TabsContent
                        value="prompt"
                        forceMount
                        style={{
                            display: tab !== "prompt" ? "none" : undefined,
                        }}
                    >
                        <PromptTab />
                    </TabsContent>
                    <TabsContent
                        value="submissions"
                        forceMount
                        style={{
                            display: tab !== "submissions" ? "none" : undefined,
                        }}
                    >
                        <SubmissionsTab />
                    </TabsContent>
                    <TabsContent
                        value="leaderboard"
                        forceMount
                        style={{
                            display: tab !== "leaderboard" ? "none" : undefined,
                        }}
                    >
                        <LeaderboardTab />
                    </TabsContent>
                    {/* <TabsContent
						value="store"
						forceMount
						style={{
							display: tab !== "store" ? "none" : undefined,
						}}
					>
						<StoreTab />
					</TabsContent> */}
                </ScrollArea>
            </div>
        </Tabs>
    );
}

export default InfoPanel;
