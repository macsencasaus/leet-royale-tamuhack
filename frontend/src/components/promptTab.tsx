import useWebSocket from "@/hooks/useWebSocket";
import { Message, VisibleCases } from "@/lib/types";
import { useCallback, useState } from "react";

function PromptTab() {
    const [cases, setCases] = useState<VisibleCases[]>([]);
    const [prompt, setPrompt] = useState("No prompt yet.");

    const onMessage = useCallback((message: Message) => {
        switch (message.type) {
            case "ServerMessageRoundStart":
                setPrompt(message.prompt);
                setCases(message.visibleTestCases);
                break;
        }
    }, []);

    useWebSocket(onMessage);

    return (
        <div className="flex flex-col gap-2">
            <div dangerouslySetInnerHTML={{ __html: prompt }}></div>
            {cases.map((value) => (
                <div className="border border-white/10 p-2 rounded">
                    <p>
                        <span className="font-bold">Input</span>: <code>{value.input}</code>
                    </p>
                    <p>
                        <span className="font-bold">Expected Output</span>:{" "}
                        <code>{value.output}</code>
                    </p>
                </div>
            ))}
        </div>
    );
}

export default PromptTab;
