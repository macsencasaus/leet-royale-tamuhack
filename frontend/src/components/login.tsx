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

function Login() {
	const { createWebSocket } = useWebSocket();
	const [name, setName] = useState("");

	function connect() {
		const address = `ws://${window.location.host}/ws?name=${name}`;
		console.log(`Sending connection request to ${address}`);
		createWebSocket(address);
	}

	return (
		<div className="flex justify-center h-dvh items-center">
			<Card className="max-w-screen-lg">
				<CardHeader>
					<CardTitle>LeetGuys</CardTitle>
					<CardDescription>
						LeetCode has never been more fun.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<Input
						placeholder="Display Name"
						value={name}
						onChange={(e) => setName(e.target.value)}
					/>
				</CardContent>
				<CardFooter className="flex justify-end">
					<Button
						onClick={connect}
						variant={"outline"}
                        disabled={name === ""}
					>
						Connect
					</Button>
				</CardFooter>
			</Card>
		</div>
	);
}

export default Login;
