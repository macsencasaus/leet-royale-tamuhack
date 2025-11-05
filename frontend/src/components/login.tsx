import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { Input } from "./ui/input";
import useWebSocket from "@/hooks/useWebSocket";
import { Button } from "./ui/button";
import { useState } from "react";
import Rules from "./rules";
import { ws_scheme } from "@/config";

function Login() {
    const { gameState, createWebSocket } = useWebSocket();
    const [name, setName] = useState("");

    function connect(ev: any) {
        ev.preventDefault();
        const address = `${ws_scheme}://${window.location.host}/ws?name=${name}`;
        console.log(`Sending connection request to ${address}`);
        createWebSocket(address);
    }

    return (
        <div
            className="flex justify-center h-dvh items-center gap-2 rounded"
            style={{ display: gameState === "login" ? undefined : "none" }}
        >
            <form onSubmit={connect}>
                <Card className="w-screen max-w-xs">
                    <CardHeader>
                        <CardTitle className="flex gap-2 items-center">
                            <img src="/leetroyale.png" className="h-8 w-8" />
                            Leet Royale
                        </CardTitle>
                        <CardDescription>LeetCode has never been more fun.</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <Input
                            placeholder="Display Name"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                        />
                    </CardContent>
                    <CardFooter className="flex justify-end">
                        <Button onClick={connect} variant={"outline"} disabled={name === ""}>
                            Connect
                        </Button>
                    </CardFooter>
                </Card>
            </form>

            <Rules />
        </div>
    );
}

export default Login;
