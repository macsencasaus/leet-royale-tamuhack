import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

function Rules() {
    return (
        <Card className="w-screen max-w-xs">
            <CardHeader>
                <CardTitle>Rules of the Game</CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-2">
                <p>There are 4 rounds.</p>
                <p>
                    The first two are timed: complete the problem within the time limit, and you'll
                    move on.
                </p>
                <p>
                    The last two are elimination: only the top move on, and only the first wins in
                    the final round.
                </p>
            </CardContent>
        </Card>
    );
}

export default Rules;
