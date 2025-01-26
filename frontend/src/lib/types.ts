export type Languages = "python" | "javascript" | "cpp";
export type MessageType = "ClientMessageReady";
export interface Message {
	type: MessageType;
}
