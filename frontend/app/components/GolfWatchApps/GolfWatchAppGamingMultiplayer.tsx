import { useAtomValue } from "jotai";
import type { components } from "../../api/schema";
import { gamingLeftTimeSecondsAtom } from "../../states/watch";
import BorderedContainer from "../BorderedContainer";
import UserLabel from "../UserLabel";

type RankingEntry = components["schemas"]["RankingEntry"];

type Props = {
	gameDisplayName: string;
	ranking: RankingEntry[];
	problemTitle: string;
	problemDescription: string;
	gameResult: "winA" | "winB" | "draw" | null;
};

export default function GolfWatchAppGamingMultiplayer({
	gameDisplayName,
	ranking,
	problemTitle,
	problemDescription,
	gameResult,
}: Props) {
	const leftTimeSeconds = useAtomValue(gamingLeftTimeSecondsAtom)!;

	const leftTime = (() => {
		const m = Math.floor(leftTimeSeconds / 60);
		const s = leftTimeSeconds % 60;
		return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
	})();

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
				<div className="font-bold flex justify-between my-auto"></div>
				<div className="font-bold text-center">
					<div className="text-gray-100">{gameDisplayName}</div>
					<div className="text-3xl">{leftTime}</div>
				</div>
				<div className="font-bold flex justify-between my-auto"></div>
			</div>
			<div className="grow grid grid-cols-2 p-4 gap-4">
				<div className="flex flex-col gap-4">
					<div>
						<div className="mb-2 text-center text-xl font-bold">
							{problemTitle}
						</div>
						<BorderedContainer>
							<pre className="text-gray-700 whitespace-pre-wrap break-words">
								{problemDescription}
							</pre>
						</BorderedContainer>
					</div>
				</div>
				<div>
					<table className="min-w-full divide-y divide-gray-200">
						<thead className="bg-gray-50">
							<tr>
								<th
									scope="col"
									className="px-6 py-3 text-left font-medium text-gray-800 uppercase tracking-wider"
								>
									順位
								</th>
								<th
									scope="col"
									className="px-6 py-3 text-left font-medium text-gray-800 uppercase tracking-wider"
								>
									名前
								</th>
								<th
									scope="col"
									className="px-6 py-3 text-left font-medium text-gray-800 uppercase tracking-wider"
								>
									スコア
								</th>
							</tr>
						</thead>
						<tbody className="bg-white divide-y divide-gray-200">
							{ranking.map((entry, index) => (
								<tr key={entry.player.user_id}>
									<td className="px-6 py-4 whitespace-nowrap text-gray-900">
										{index + 1}
									</td>
									<td className="px-6 py-4 whitespace-nowrap text-gray-900">
										{entry.player.display_name}
										{entry.player.label && (
											<UserLabel label={entry.player.label} />
										)}
									</td>
									<td className="px-6 py-4 whitespace-nowrap text-gray-900">
										{entry.score}
									</td>
								</tr>
							))}
						</tbody>
					</table>
				</div>
			</div>
		</div>
	);
}
