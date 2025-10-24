import useWebSocket from "@/hooks/useWebSocket";
import { Message, Player } from "@/lib/types";
import { useCallback, useState } from "react";

function Breather() {
    const [all, setAll] = useState<Player[]>([]);
    const [dead, setDead] = useState<Player[]>([]);

    const onMessage = useCallback((message: Message) => {
        switch (message.type) {
            case "ServerMessageRoundEnd":
                setAll(
                    shuffle(
                        message.currentPlayers.concat(
                            message.eliminatedPlayers !== null ? message.eliminatedPlayers : []
                        )
                    )
                );
                setDead(message.eliminatedPlayers !== null ? message.eliminatedPlayers : []);
                break;
        }
    }, []);

    function shuffle<T>(array: T[]): T[] {
        for (let i = array.length - 1; i >= 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [array[i], array[j]] = [array[j], array[i]];
        }

        return array;
    }

    function isDead(player: Player): boolean {
        return dead.includes(player);
    }

    useWebSocket(onMessage);

    return (
        <div className="border border-white/10 p-2 rounded h-full flex flex-col gap-2">
            <h1>You made it! Many others didn't, however.</h1>
            <div className="flex flex-wrap gap-2">
                {all.map((value) => (
                    <div
                        className="border border-white/10 p-2 rounded"
                        style={{
                            backgroundColor: isDead(value) ? "rgba(220, 53, 69, 0.5)" : undefined,
                        }}
                    >
                        {value.name}
                    </div>
                ))}
            </div>
        </div>
    );
}

export default Breather;
