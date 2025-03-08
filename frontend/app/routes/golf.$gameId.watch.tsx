import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { useHydrateAtoms } from "jotai/utils";
import { ensureUserLoggedIn } from "../.server/auth";
import {
	ApiAuthTokenContext,
	apiGetGame,
	apiGetGameWatchLatestStates,
	apiGetGameWatchRanking,
} from "../api/client";
import GolfWatchApp from "../components/GolfWatchApp";
import {
	rankingAtom,
	setCurrentTimestampAtom,
	setDurationSecondsAtom,
	setGameStartedAtAtom,
	setLatestGameStatesAtom,
} from "../states/watch";

export const meta: MetaFunction<typeof loader> = ({ data }) => [
	{
		title: data
			? `Golf Watching ${data.game.display_name} | PHPerKaigi 2025 Albatross`
			: "Golf Watching | PHPerKaigi 2025 Albatross",
	},
];

export async function loader({ params, request }: LoaderFunctionArgs) {
	const { token } = await ensureUserLoggedIn(request);

	const gameId = Number(params.gameId);

	const fetchGame = async () => {
		return (await apiGetGame(token, gameId)).game;
	};
	const fetchRanking = async () => {
		return (await apiGetGameWatchRanking(token, gameId)).ranking;
	};
	const fetchGameStates = async () => {
		return (await apiGetGameWatchLatestStates(token, gameId)).states;
	};

	const [game, ranking, gameStates] = await Promise.all([
		fetchGame(),
		fetchRanking(),
		fetchGameStates(),
	]);

	return {
		apiAuthToken: token,
		game,
		ranking,
		gameStates,
	};
}

export default function GolfWatch() {
	const { apiAuthToken, game, ranking, gameStates } =
		useLoaderData<typeof loader>();

	useHydrateAtoms([
		[rankingAtom, ranking],
		[setCurrentTimestampAtom, undefined],
		[setDurationSecondsAtom, game.duration_seconds],
		[setGameStartedAtAtom, game.started_at ?? null],
		[setLatestGameStatesAtom, gameStates],
	]);

	return (
		<ApiAuthTokenContext.Provider value={apiAuthToken}>
			<GolfWatchApp game={game} />
		</ApiAuthTokenContext.Provider>
	);
}
