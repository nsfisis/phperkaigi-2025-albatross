import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { useHydrateAtoms } from "jotai/utils";
import { ensureUserLoggedIn } from "../.server/auth";
import {
	ApiAuthTokenContext,
	apiGetGame,
	apiGetGamePlayLatestState,
} from "../api/client";
import GolfPlayApp from "../components/GolfPlayApp";
import {
	setCurrentTimestampAtom,
	setDurationSecondsAtom,
	setGameStartedAtAtom,
	setLatestGameStateAtom,
} from "../states/play";

export const meta: MetaFunction<typeof loader> = ({ data }) => [
	{
		title: data
			? `Golf Playing ${data.game.display_name} | PHPerKaigi 2025 Albatross`
			: "Golf Playing | PHPerKaigi 2025 Albatross",
	},
];

export async function loader({ params, request }: LoaderFunctionArgs) {
	const { token, user } = await ensureUserLoggedIn(request);

	const gameId = Number(params.gameId);

	const fetchGame = async () => {
		return (await apiGetGame(token, gameId)).game;
	};
	const fetchGameState = async () => {
		return (await apiGetGamePlayLatestState(token, gameId)).state;
	};

	const [game, state] = await Promise.all([fetchGame(), fetchGameState()]);

	return {
		apiAuthToken: token,
		game,
		player: user,
		gameState: state,
	};
}

export default function GolfPlay() {
	const { apiAuthToken, game, player, gameState } =
		useLoaderData<typeof loader>();

	useHydrateAtoms([
		[setCurrentTimestampAtom, undefined],
		[setDurationSecondsAtom, game.duration_seconds],
		[setGameStartedAtAtom, game.started_at ?? null],
		[setLatestGameStateAtom, gameState],
	]);

	return (
		<ApiAuthTokenContext.Provider value={apiAuthToken}>
			<GolfPlayApp game={game} player={player} initialCode={gameState.code} />
		</ApiAuthTokenContext.Provider>
	);
}
