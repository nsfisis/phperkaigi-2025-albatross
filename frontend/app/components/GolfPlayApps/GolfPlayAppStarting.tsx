import { useAtomValue } from "jotai";
import { startingLeftTimeSecondsAtom } from "../../states/play";

type Props = {
	gameDisplayName: string;
};

export default function GolfPlayAppStarting({ gameDisplayName }: Props) {
	const leftTimeSeconds = useAtomValue(startingLeftTimeSecondsAtom)!;

	return (
		<div className="min-h-screen bg-gray-100 flex flex-col">
			<div className="text-white bg-sky-600 p-10 text-center">
				<div className="text-4xl font-bold">{gameDisplayName}</div>
			</div>
			<div className="text-center text-black font-black text-10xl">
				{leftTimeSeconds}
			</div>
		</div>
	);
}
