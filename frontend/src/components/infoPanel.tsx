import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"


function InfoPanel() {
    return (
		<Tabs
			defaultValue="question"
			className="flex flex-col"
		>
			<div className="bg-border overflow-hidden -m-2 p-1">
				<TabsList className="self-start">
					<TabsTrigger value="question">Question</TabsTrigger>
					<TabsTrigger value="leaderboard">Leaderboard</TabsTrigger>
					<TabsTrigger value="store">Store</TabsTrigger>
					<TabsTrigger value="submissions">Submissions</TabsTrigger>
				</TabsList>
			</div>

			<div className="p-2">
				<TabsContent value="question">
					Questionsasd ioasnd ajsiodj asojd aisod
					oiasdiojasoidjaoishduiadasn oudnoias dh oasjdsaj
					iodhaifojaesdhiamsknocjb dundiaksjbuns dobuai
				</TabsContent>
				<TabsContent value="leaderboard">Leaderboard</TabsContent>
				<TabsContent value="store">Store</TabsContent>
				<TabsContent value="submissions">Submissions</TabsContent>
			</div>
		</Tabs>
	);
}

export default InfoPanel;