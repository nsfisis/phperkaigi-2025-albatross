import type { PlayerProfile } from "../../types/PlayerProfile";
import PlayerNameAndIcon from "../PlayerNameAndIcon";

type Props = {
	gameDisplayName: string;
	playerProfileA: PlayerProfile;
	playerProfileB: PlayerProfile;
};

export default function GolfWatchAppWaiting1v1({
	gameDisplayName,
	playerProfileA,
	playerProfileB,
}: Props) {
	return (
		<div className="min-h-screen bg-gray-100 flex flex-col font-bold text-center">
			<div className="text-white bg-sky-600 p-10">
				<div className="text-4xl">{gameDisplayName}</div>
			</div>
			<div className="grow grid grid-cols-3 gap-10 mx-auto text-black">
				<PlayerNameAndIcon profile={playerProfileA} />
				<div className="text-8xl my-auto">vs.</div>
				<PlayerNameAndIcon profile={playerProfileB} />
			</div>
		</div>
	);
}
