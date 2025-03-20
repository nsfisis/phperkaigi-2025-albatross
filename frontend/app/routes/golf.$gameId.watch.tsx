import { Provider as JotaiProvider, createStore } from "jotai";
import { useMemo } from "react";
import type { LoaderFunctionArgs, MetaFunction } from "react-router";
import { useLoaderData } from "react-router";
import { ensureUserLoggedIn } from "../.server/auth";
import {
	ApiAuthTokenContext,
	apiGetGame,
	apiGetGameWatchLatestStates,
	apiGetGameWatchRanking,
} from "../api/client";
import GolfWatchApp from "../components/GolfWatchApp";

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

	const store = useMemo(() => {
		void game.game_id;
		return createStore();
	}, [game.game_id]);

	return (
		<JotaiProvider store={store}>
			<ApiAuthTokenContext.Provider value={apiAuthToken}>
				<GolfWatchApp
					key={game.game_id}
					game={game}
					initialGameStates={gameStates}
					initialRanking={ranking}
				/>
			</ApiAuthTokenContext.Provider>
		</JotaiProvider>
	);
}
