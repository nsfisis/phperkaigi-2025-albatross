import { Provider as JotaiProvider, createStore } from "jotai";
import { useMemo } from "react";
import type { LoaderFunctionArgs, MetaFunction } from "react-router";
import { useLoaderData } from "react-router";
import { ensureUserLoggedIn } from "../.server/auth";
import { ApiClientContext, createApiClient } from "../api/client";
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
	const apiClient = createApiClient(token);

	const gameId = Number(params.gameId);

	const [{ game }, { ranking }, { states: gameStates }] = await Promise.all([
		await apiClient.getGame(gameId),
		await apiClient.getGameWatchRanking(gameId),
		await apiClient.getGameWatchLatestStates(gameId),
	]);

	return {
		apiToken: token,
		game,
		ranking,
		gameStates,
	};
}

export default function GolfWatch() {
	const { apiToken, game, ranking, gameStates } =
		useLoaderData<typeof loader>();

	const store = useMemo(() => {
		void game.game_id;
		return createStore();
	}, [game.game_id]);

	return (
		<JotaiProvider store={store}>
			<ApiClientContext.Provider value={createApiClient(apiToken)}>
				<GolfWatchApp
					key={game.game_id}
					game={game}
					initialGameStates={gameStates}
					initialRanking={ranking}
				/>
			</ApiClientContext.Provider>
		</JotaiProvider>
	);
}
