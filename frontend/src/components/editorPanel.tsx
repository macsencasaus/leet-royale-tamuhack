import { Languages, Message, Templates } from "@/lib/types";
import _Editor from "@monaco-editor/react";
import { useCallback, useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Send, SkipForward } from "lucide-react";
import { Button } from "./ui/button";
import useWebSocket from "@/hooks/useWebSocket";
import { demo } from "@/config";

// const template: Templates = {
// 	javascript: "// js",
// 	cpp: "// cpp",
// 	python: "# py",
// };

function EditorPanel() {
	const [templates, setTemplates] = useState<Templates | undefined>(
		undefined
	);
	const [language, setLanguage] = useState<Languages>("javascript");
	const [code, setCode] = useState<Record<Languages, string>>({
		javascript: "",
		python: "",
		cpp: "",
	});

	const onMessage = useCallback((message: Message) => {
		switch (message.type) {
			case "ServerMessageRoundStart":
				setTemplates(message.templates);
				setCode(message.templates);
				break;
			case "ServerMessageRoundEnd":
				setTemplates(undefined);
				break;
		}
	}, []);

	const { player, sendMessage } = useWebSocket(onMessage);

	function onCode(_c: string | undefined) {
		const next: any = {};
		next[language] = _c;

		setCode((prev) => {
			return {
				...prev,
				...next,
			};
		});
	}

	function submitCode() {
		console.log(code);

		sendMessage("ClientMessageSubmit", {
			playerId: player?.id,
			language,
			code: code[language],
		});
	}

	function skipQuestion() {
		sendMessage("ClientMessageSkipQuestion");
	}

	return (
		<Tabs
			defaultValue="javascript"
			className="flex flex-col h-full"
			onValueChange={(value) => setLanguage(value as Languages)}
		>
			<div className="bg-border overflow-hidden -m-2 p-1 flex justify-between">
				<TabsList className="self-start">
					<TabsTrigger value="javascript">JavaScript</TabsTrigger>
					<TabsTrigger value="python">Python</TabsTrigger>
					<TabsTrigger value="cpp">C++</TabsTrigger>
				</TabsList>
				<div className="flex gap-2">
					{demo && (
						<Button
							variant={"outline"}
							onClick={skipQuestion}
						>
							<SkipForward /> Skip
						</Button>
					)}
					<Button
						variant={"outline"}
						style={{ backgroundColor: "rgba(25, 135, 84, .5)" }}
						onClick={submitCode}
					>
						<Send /> Submit
					</Button>
				</div>
			</div>

			<div className="py-4 px-2 overflow-hidden grow">
				{templates ? (
					<div className="-ml-8 h-full">
						<TabsContent
							value="javascript"
							className="h-full"
							forceMount
							style={{
								display:
									language !== "javascript"
										? "none"
										: undefined,
							}}
						>
							<_Editor
								height="100%"
								width="100%"
								defaultLanguage="javascript"
								theme="vs-dark"
								defaultValue={templates.javascript}
								className="rounded"
								value={code["javascript"]}
								onChange={onCode}
							/>
						</TabsContent>
						<TabsContent
							value="python"
							className="h-full"
							forceMount
							style={{
								display:
									language !== "python" ? "none" : undefined,
							}}
						>
							<_Editor
								height="100%"
								width="100%"
								defaultLanguage="python"
								theme="vs-dark"
								defaultValue={templates.python}
								className="rounded"
								value={code["python"]}
								onChange={onCode}
							/>
						</TabsContent>
						<TabsContent
							value="cpp"
							className="h-full"
							forceMount
							style={{
								display:
									language !== "cpp" ? "none" : undefined,
							}}
						>
							<_Editor
								height="100%"
								width="100%"
								defaultLanguage="c++"
								language={language}
								theme="vs-dark"
								defaultValue={templates.cpp}
								className="rounded"
								value={code["cpp"]}
								onChange={onCode}
							/>
						</TabsContent>
					</div>
				) : (
					<p>
						Hold your horses,{" "}
						<span className="font-bold">{player?.name}</span>. The
						round hasn't started yet.
					</p>
				)}
			</div>
		</Tabs>
	);
}

export default EditorPanel;
