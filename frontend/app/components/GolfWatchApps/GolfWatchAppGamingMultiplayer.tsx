import { useAtomValue } from "jotai";
import type { components } from "../../api/schema";
import { gamingLeftTimeSecondsAtom } from "../../states/watch";
import LeftTime from "../Gaming/LeftTime";
import Problem from "../Gaming/Problem";
import UserLabel from "../UserLabel";

type RankingEntry = components["schemas"]["RankingEntry"];

type Props = {
	gameDisplayName: string;
	ranking: RankingEntry[];
	problemTitle: string;
	problemDescription: string;
	sampleCode: string;
};

export default function GolfWatchAppGamingMultiplayer({
	gameDisplayName,
	ranking,
	problemTitle,
	problemDescription,
	sampleCode,
}: Props) {
	const leftTimeSeconds = useAtomValue(gamingLeftTimeSecondsAtom)!;

	return (
		<div className="min-h-screen bg-gray-100 flex flex-col">
			<div className={`text-white bg-sky-600 grid grid-cols-3 px-4 py-2`}>
				<div className="font-bold flex justify-between my-auto"></div>
				<div className="font-bold text-center">
					<div className="text-gray-100">{gameDisplayName}</div>
					<LeftTime sec={leftTimeSeconds} />
				</div>
				<div className="font-bold flex justify-between my-auto"></div>
			</div>
			<div className="grow grid grid-cols-2 p-4 gap-4">
				<Problem
					title={problemTitle}
					description={problemDescription}
					sampleCode={sampleCode}
				/>
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
											<span className="mx-2">
												<UserLabel label={entry.player.label} />
											</span>
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
