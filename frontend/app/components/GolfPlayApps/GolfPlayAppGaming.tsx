import { Link } from "@remix-run/react";
import React, { useRef } from "react";
import SubmitButton from "../../components/SubmitButton";
import type { PlayerInfo } from "../../models/PlayerInfo";
import BorderedContainer from "../BorderedContainer";
import SubmitStatusLabel from "../SubmitStatusLabel";
import ExecStatusIndicatorIcon from "../ExecStatusIndicatorIcon";
import { faArrowDown } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

type Props = {
	gameDisplayName: string;
	playerInfo: Omit<PlayerInfo, "code">;
	problemTitle: string;
	problemDescription: string;
	onCodeChange: (code: string) => void;
	onCodeSubmit: (code: string) => void;
};

export default function GolfPlayAppGaming({
	gameDisplayName,
	playerInfo,
	problemTitle,
	problemDescription,
	onCodeChange,
	onCodeSubmit,
}: Props) {
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
			<div className="text-white bg-iosdc-japan flex flex-row justify-between px-4 py-2">
				<div className="font-bold">
					<div className="text-gray-100">{gameDisplayName}</div>
					<div className="text-2xl">03:21</div>
				</div>
				<div className="font-bold text-end">
					<Link to={"/dashboard"} className="text-gray-100">
						{playerInfo.displayName}
					</Link>
					<div className="text-2xl">{playerInfo.score}</div>
				</div>
			</div>
			<div className="grow grid grid-cols-3 divide-x divide-gray-300">
				<div className="p-4">
					<div className="mb-2 text-xl font-bold">{problemTitle}</div>
					<div className="p-2">
						<BorderedContainer>
							<div className="text-gray-700">{problemDescription}</div>
						</BorderedContainer>
					</div>
				</div>
				<div className="p-4">
					<textarea
						ref={textareaRef}
						onChange={handleTextChange}
						className="resize-none h-full w-full rounded-lg border border-gray-300 p-2 focus:outline-none focus:ring-2 focus:ring-gray-400 transition duration-300"
					></textarea>
				</div>
				<div className="p-4 flex flex-col gap-4">
					<div className="flex">
						<SubmitButton onClick={handleSubmitButtonClick}>提出</SubmitButton>
						<div className="grow font-bold text-xl text-center m-1">
							<SubmitStatusLabel status={playerInfo.submitResult.status} />
						</div>
					</div>
					<ul className="flex flex-col gap-2">
						{playerInfo.submitResult.execResults.map((r, idx) => (
							<li key={r.testcase_id ?? -1} className="flex gap-2">
								<div className="flex flex-col gap-2 p-2">
									<div className="w-6">
										<ExecStatusIndicatorIcon status={r.status} />
									</div>
									{idx !== playerInfo.submitResult.execResults.length - 1 && (
										<div>
											<FontAwesomeIcon
												icon={faArrowDown}
												fixedWidth
												className="text-gray-500"
											/>
										</div>
									)}
								</div>
								<div className="grow p-2 overflow-x-scroll">
									<BorderedContainer>
										<div className="font-semibold">{r.label}</div>
										<div>
											<code>
												{r.stdout}
												{r.stderr}
											</code>
										</div>
									</BorderedContainer>
								</div>
							</li>
						))}
					</ul>
				</div>
			</div>
		</div>
	);
}
