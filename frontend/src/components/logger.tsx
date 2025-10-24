import useWebSocket from "@/hooks/useWebSocket";
import { Message } from "@/lib/types";
import { useCallback } from "react";

function Logger() {
    const onMessage = useCallback((message: Message) => {
        console.log(message);
    }, []);

    useWebSocket(onMessage);

    return <></>;
}

export default Logger;
