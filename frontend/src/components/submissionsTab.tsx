import useWebSocket from "@/hooks/useWebSocket";
import { Message, TestCase } from "@/lib/types";
import { useCallback, useState } from "react";
import {
	Accordion,
	AccordionContent,
	AccordionItem,
	AccordionTrigger,
} from "@/components/ui/accordion";
import { Ban, Check, TriangleAlert } from "lucide-react";

function SubmissionsTab() {
	const [tle, setTle] = useState(false);
	const [results, setResults] = useState<TestCase[]>([]);

	const onMessage = useCallback((message: Message) => {
		switch (message.type) {
			case "ServerMessageRoundStart":
				
				break;
			case "ServerMessageTestResult":
				setTle(message.tle);
				setResults(message.cases);
				break;
		}
	}, []);

	useWebSocket(onMessage);

	return (
		<div className="flex flex-col gap-2 h-full">
			{tle && (
				<div>
					<div
						className="flex gap-2 p-2 rounded-t items-center"
						style={{ backgroundColor: "rgba(255, 193, 7, .5)" }}
					>
						<TriangleAlert />
						Time Limit Exceeded
					</div>
					<div className="border border-white/10 p-2 rounded-b">
						Not all tests cases may have finished.
					</div>
				</div>
			)}
			{results.length > 0 && (
				<Accordion
					type="multiple"
					className="flex flex-col gap-2"
				>
					{results.map((value, index) => (
						<AccordionItem
							value={`case-${index}`}
							key={index}
						>
							<AccordionTrigger
								className="border border-white/10 p-2 rounded data-[state=open]:rounded-b-none"
								style={{
									backgroundColor: value.success
										? "rgba(25, 135, 84, .5)"
										: "rgba(220, 53, 69, .5)",
								}}
							>
								<div className="flex gap-2 items-center">
									{value.success ? <Check /> : <Ban />} Test
									Case {index + 1}
								</div>
							</AccordionTrigger>
							<AccordionContent className="border border-white/10 p-2 rounded-b">
								{value.stdout ? (
									<>
										<p className="font-bold">Output:</p>
										<pre>{value.stdout}</pre>
									</>
								) : (
									<p>No output was produced.</p>
								)}
							</AccordionContent>
						</AccordionItem>
					))}
				</Accordion>
			)}
			{!tle && results.length <= 0 && (
				<p>You haven't submitted anything yet.</p>
			)}
		</div>
	);
}

export default SubmissionsTab;
