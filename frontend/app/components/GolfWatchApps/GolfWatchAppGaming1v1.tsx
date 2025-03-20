import { useAtomValue } from "jotai";
import {
	calcCodeSize,
	gamingLeftTimeSecondsAtom,
	latestGameStatesAtom,
} from "../../states/watch";
import type { PlayerProfile } from "../../types/PlayerProfile";
import FoldableBorderedContainerWithCaption from "../FoldableBorderedContainerWithCaption";
import CodeBlock from "../Gaming/CodeBlock";
import LeftTime from "../Gaming/LeftTime";
import ProblemColumn from "../Gaming/ProblemColumn";
import ScoreBar from "../Gaming/ScoreBar";
import SubmitStatusLabel from "../SubmitStatusLabel";
import ThreeColumnLayout from "../ThreeColumnLayout";
import TitledColumn from "../TitledColumn";
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

	const codeSizeA = calcCodeSize(codeA);
	const codeSizeB = calcCodeSize(codeB);

	const topBg = gameResult
		? gameResult === "winA"
			? "bg-orange-400"
			: gameResult === "winB"
				? "bg-purple-400"
				: "bg-sky-600"
		: "bg-sky-600";

	return (
		<div className="min-h-screen bg-gray-100 flex flex-col">
			<div className={`text-white ${topBg} grid grid-cols-3 px-4 py-2`}>
				<div className="font-bold flex justify-between my-auto">
					<div className="flex gap-6 items-center">
						{playerProfileA.iconPath && (
							<UserIcon
								iconPath={playerProfileA.iconPath}
								displayName={playerProfileA.displayName}
								className="w-12 h-12 my-auto"
							/>
						)}
						<div className="text-4xl">{playerProfileA.displayName}</div>
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
					<div className="flex gap-6 items-center text-end">
						<div className="text-4xl">{playerProfileB.displayName}</div>
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
			<ThreeColumnLayout>
				<TitledColumn title={<SubmitStatusLabel status={statusA} />}>
					<FoldableBorderedContainerWithCaption
						caption={`コードサイズ: ${codeSizeA}`}
					>
						<CodeBlock code={codeA} language="php" />
					</FoldableBorderedContainerWithCaption>
				</TitledColumn>
				<ProblemColumn
					title={problemTitle}
					description={problemDescription}
					sampleCode={sampleCode}
				/>
				<TitledColumn title={<SubmitStatusLabel status={statusB} />}>
					<FoldableBorderedContainerWithCaption
						caption={`コードサイズ: ${codeSizeB}`}
					>
						<CodeBlock code={codeB} language="php" />
					</FoldableBorderedContainerWithCaption>
				</TitledColumn>
			</ThreeColumnLayout>
		</div>
	);
}
