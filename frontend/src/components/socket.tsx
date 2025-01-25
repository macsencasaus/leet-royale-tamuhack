import useWebSocket from "@/hooks/useWebSocket";

function Socket() {
	const { createWebSocket } = useWebSocket();

	function connect() {
		const address = `ws://${window.location.host}/ws`;
		console.log(`Sending connection request to ${address}`);
		createWebSocket(address);
	}

	return (
		<button
			onClick={connect}
			className="border border-white/10 rounded p-2"
		>
			Connect
		</button>
	);
}

export default Socket;
