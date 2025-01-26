import { Languages } from "@/lib/types";
import _Editor from "@monaco-editor/react";
import { useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Send } from "lucide-react";
import { Button } from "./ui/button";

const solution = `function solution() {
    // TODO
}`;

function EditorPanel() {
	const [language, setLanguage] = useState<Languages>("javascript");

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
				<Button
					variant={"outline"}
					style={{ backgroundColor: "rgb(25, 135, 84)" }}
				>
					<Send /> Submit
				</Button>
			</div>

			<div className="py-2 overflow-hidden grow">
				<div className="-ml-8 h-full">
					<TabsContent
						value="javascript"
						className="h-full"
					>
						<_Editor
							height="100%"
							width="100%"
							defaultLanguage="javascript"
							theme="vs-dark"
							defaultValue={solution}
							className="rounded"
						/>
					</TabsContent>
					<TabsContent
						value="python"
						className="h-full"
					>
						<_Editor
							height="100%"
							width="100%"
							defaultLanguage="python"
							theme="vs-dark"
							defaultValue={solution}
							className="rounded"
						/>
					</TabsContent>
					<TabsContent
						value="cpp"
						className="h-full"
					>
						<_Editor
							height="100%"
							width="100%"
							defaultLanguage="c++"
							language={language}
							theme="vs-dark"
							defaultValue={solution}
							className="rounded"
						/>
					</TabsContent>
				</div>
			</div>
		</Tabs>
	);
}

export default EditorPanel;
