import { useAtomValue } from "jotai";
import {
	gamingLeftTimeSecondsAtom,
	latestGameStatesAtom,
} from "../../states/watch";
import type { PlayerProfile } from "../../types/PlayerProfile";
import CodeBlock from "../Gaming/CodeBlock";
import LeftTime from "../Gaming/LeftTime";
import Problem from "../Gaming/Problem";
import ScoreBar from "../Gaming/ScoreBar";
import SubmitResult from "../Gaming/SubmitResult";
import UserIcon from "../UserIcon";

type Props = {
	gameDisplayName: string;
	playerProfileA: PlayerProfile;
	playerProfileB: PlayerProfile;
	problemTitle: string;
	problemDescription: string;
	sampleCode: string;
	gameResult: "winA" | "winB" | "draw" | null;
};

export default function GolfWatchAppGaming1v1({
	gameDisplayName,
	playerProfileA,
	playerProfileB,
	problemTitle,
	problemDescription,
	sampleCode,
	gameResult,
}: Props) {
	const leftTimeSeconds = useAtomValue(gamingLeftTimeSecondsAtom)!;
	const latestGameStates = useAtomValue(latestGameStatesAtom);

	const stateA = latestGameStates[`${playerProfileA.id}`];
	const codeA = stateA?.code ?? "";
	const scoreA = stateA?.score ?? null;
	const statusA = stateA?.status ?? "none";
	const stateB = latestGameStates[`${playerProfileB.id}`];
	const codeB = stateB?.code ?? "";
	const scoreB = stateB?.score ?? null;
	const statusB = stateB?.status ?? "none";

	const topBg = gameResult
		? gameResult === "winA"
			? "bg-orange-400"
			: gameResult === "winB"
				? "bg-purple-400"
				: "bg-pink-500"
		: "bg-sky-600";

	return (
		<div className="min-h-screen bg-gray-100 flex flex-col">
			<div className={`text-white ${topBg} grid grid-cols-3 px-4 py-2`}>
				<div className="font-bold flex justify-between my-auto">
					<div className="flex gap-6">
						{playerProfileA.iconPath && (
							<UserIcon
								iconPath={playerProfileA.iconPath}
								displayName={playerProfileA.displayName}
								className="w-12 h-12 my-auto"
							/>
						)}
						<div>
							<div className="text-gray-100">Player 1</div>
							<div className="text-2xl">{playerProfileA.displayName}</div>
						</div>
					</div>
					<div className="text-6xl">{scoreA}</div>
				</div>
				<div className="font-bold text-center">
					<div className="text-gray-100">{gameDisplayName}</div>
					{gameResult ? (
						<div className="text-3xl">
							{gameResult === "winA"
								? `勝者 ${playerProfileA.displayName}`
								: gameResult === "winB"
									? `勝者 ${playerProfileB.displayName}`
									: "引き分け"}
						</div>
					) : (
						<LeftTime sec={leftTimeSeconds} />
					)}
				</div>
				<div className="font-bold flex justify-between my-auto">
					<div className="text-6xl">{scoreB}</div>
					<div className="flex gap-6 text-end">
						<div>
							<div className="text-gray-100">Player 2</div>
							<div className="text-2xl">{playerProfileB.displayName}</div>
						</div>
						{playerProfileB.iconPath && (
							<UserIcon
								iconPath={playerProfileB.iconPath}
								displayName={playerProfileB.displayName}
								className="w-12 h-12 my-auto"
							/>
						)}
					</div>
				</div>
			</div>
			<ScoreBar
				scoreA={scoreA}
				scoreB={scoreB}
				bgA="bg-orange-400"
				bgB="bg-purple-400"
			/>
			<div className="grow grid grid-cols-3 p-4 gap-4">
				<CodeBlock code={codeA} language="php" />
				<div className="flex flex-col gap-4">
					<div className="grid grid-cols-2 gap-4">
						<SubmitResult status={statusA} />
						<SubmitResult status={statusB} />
					</div>
					<Problem
						title={problemTitle}
						description={problemDescription}
						sampleCode={sampleCode}
					/>
				</div>
				<CodeBlock code={codeB} language="php" />
			</div>
		</div>
	);
}
