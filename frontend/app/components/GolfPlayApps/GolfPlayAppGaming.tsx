import { useAtomValue } from "jotai";
import React, { useRef } from "react";
import { Link } from "react-router";
import SubmitButton from "../../components/SubmitButton";
import {
	gamingLeftTimeSecondsAtom,
	scoreAtom,
	statusAtom,
} from "../../states/play";
import type { PlayerProfile } from "../../types/PlayerProfile";
import LeftTime from "../Gaming/LeftTime";
import Problem from "../Gaming/Problem";
import SubmitResult from "../Gaming/SubmitResult";
import UserIcon from "../UserIcon";

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

	const textareaRef = useRef<HTMLTextAreaElement>(null);

	const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
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
			<div className="grow grid grid-cols-3 divide-x divide-gray-300">
				<Problem
					title={problemTitle}
					description={problemDescription}
					sampleCode={sampleCode}
				/>
				<div className="p-4">
					<textarea
						ref={textareaRef}
						defaultValue={initialCode}
						onChange={handleTextChange}
						className="resize-none h-full w-full rounded-lg border border-gray-300 p-2 focus:outline-hidden focus:ring-2 focus:ring-gray-400 transition duration-300"
					/>
				</div>
				<div className="p-4">
					<SubmitResult
						status={status}
						submitButton={
							<SubmitButton onClick={handleSubmitButtonClick}>
								提出
							</SubmitButton>
						}
					/>
				</div>
			</div>
		</div>
	);
}
