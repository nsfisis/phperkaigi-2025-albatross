import { useAtomValue } from "jotai";
import React, { useRef, useState } from "react";
import { Link } from "react-router";
import {
	gamingLeftTimeSecondsAtom,
	scoreAtom,
	statusAtom,
} from "../../states/play";
import type { PlayerProfile } from "../../types/PlayerProfile";
import BorderedContainer from "../BorderedContainer";
import LeftTime from "../Gaming/LeftTime";
import Problem from "../Gaming/Problem";
import SubmitButton from "../SubmitButton";
import SubmitStatusLabel from "../SubmitStatusLabel";
import ThreeColumnLayout from "../ThreeColumnLayout";
import UserIcon from "../UserIcon";

function calcCodeSize(code: string): number {
	const trimmed = code
		.replace(/\s+/g, "")
		.replace(/^<\?php/, "")
		.replace(/^<\?/, "")
		.replace(/\?>$/, "");
	const utf8Encoded = new TextEncoder().encode(trimmed);
	return utf8Encoded.length;
}

type Props = {
	gameDisplayName: string;
	playerProfile: PlayerProfile;
	problemTitle: string;
	problemDescription: string;
	sampleCode: string;
	initialCode: string;
	onCodeChange: (code: string) => void;
	onCodeSubmit: (code: string) => void;
};

export default function GolfPlayAppGaming({
	gameDisplayName,
	playerProfile,
	problemTitle,
	problemDescription,
	sampleCode,
	initialCode,
	onCodeChange,
	onCodeSubmit,
}: Props) {
	const leftTimeSeconds = useAtomValue(gamingLeftTimeSecondsAtom)!;
	const score = useAtomValue(scoreAtom);
	const status = useAtomValue(statusAtom);

	const [codeSize, setCodeSize] = useState(calcCodeSize(initialCode));
	const textareaRef = useRef<HTMLTextAreaElement>(null);

	const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
		setCodeSize(calcCodeSize(e.target.value));
		onCodeChange(e.target.value);
	};

	const handleSubmitButtonClick = () => {
		if (textareaRef.current) {
			onCodeSubmit(textareaRef.current.value);
		}
	};

	return (
		<div className="min-h-screen bg-gray-100 flex flex-col">
			<div className="text-white bg-sky-600 flex flex-row justify-between px-4 py-2">
				<div className="font-bold">
					<div className="text-gray-100">{gameDisplayName}</div>
					<LeftTime sec={leftTimeSeconds} />
				</div>
				<Link to={"/dashboard"}>
					<div className="flex gap-4 my-auto font-bold">
						<div className="text-6xl">{score}</div>
						<div className="text-end">
							<div className="text-gray-100">Player 1</div>
							<div className="text-2xl">{playerProfile.displayName}</div>
						</div>
						{playerProfile.iconPath && (
							<UserIcon
								iconPath={playerProfile.iconPath}
								displayName={playerProfile.displayName}
								className="w-12 h-12 my-auto"
							/>
						)}
					</div>
				</Link>
			</div>
			<ThreeColumnLayout>
				<Problem
					title={problemTitle}
					description={problemDescription}
					sampleCode={sampleCode}
				/>
				<div className="p-4 flex flex-col gap-4">
					<div className="text-center text-xl font-bold">ソースコード</div>
					<BorderedContainer className="grow flex flex-col gap-4">
						<div className="flex flex-row gap-2 items-center">
							<div className="grow font-semibold text-lg">
								コードサイズ: {codeSize}
							</div>
							<SubmitButton onClick={handleSubmitButtonClick}>
								提出
							</SubmitButton>
						</div>
						<textarea
							ref={textareaRef}
							defaultValue={initialCode}
							onChange={handleTextChange}
							className="grow resize-none h-full w-full p-2 bg-gray-50 rounded-lg border border-gray-300 focus:outline-hidden focus:ring-2 focus:ring-gray-400 transition duration-300"
						/>
					</BorderedContainer>
				</div>
				<div className="p-4 flex flex-col gap-4">
					<div className="text-center text-xl font-bold">提出結果</div>
					<div className="overflow-hidden border-2 border-blue-600 rounded-xl">
						<table className="min-w-full divide-y divide-gray-400 border-collapse">
							<thead className="bg-gray-50">
								<tr>
									<th
										scope="col"
										className="px-6 py-3 text-left font-medium text-gray-800"
									>
										ステータス
									</th>
								</tr>
							</thead>
							<tbody className="bg-white divide-y divide-gray-300">
								{[status].map((status) => (
									<tr key={99999}>
										<td className="px-6 py-4 whitespace-nowrap text-gray-900">
											<SubmitStatusLabel status={status} />
										</td>
									</tr>
								))}
							</tbody>
						</table>
					</div>
					<p>
						NOTE:
						過去の提出結果を閲覧する機能は現在実装中です。それまでは提出コードをお手元に保管しておいてください。
					</p>
				</div>
			</ThreeColumnLayout>
		</div>
	);
}
