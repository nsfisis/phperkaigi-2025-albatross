import { Provider as JotaiProvider, createStore } from "jotai";
import { useMemo } from "react";
import type { LoaderFunctionArgs, MetaFunction } from "react-router";
import { useLoaderData } from "react-router";
import { ensureUserLoggedIn } from "../.server/auth";
import { ApiClientContext, createApiClient } from "../api/client";
import GolfPlayApp from "../components/GolfPlayApp";

export const meta: MetaFunction<typeof loader> = ({ data }) => [
	{
		title: data
			? `Golf Playing ${data.game.display_name} | PHPerKaigi 2025 Albatross`
			: "Golf Playing | PHPerKaigi 2025 Albatross",
	},
];

export async function loader({ params, request }: LoaderFunctionArgs) {
	const { token, user } = await ensureUserLoggedIn(request);
	const apiClient = createApiClient(token);

	const gameId = Number(params.gameId);

	const [{ game }, { state: gameState }] = await Promise.all([
		apiClient.getGame(gameId),
		apiClient.getGamePlayLatestState(gameId),
	]);

	return {
		apiToken: token,
		game,
		player: user,
		gameState,
	};
}

export default function GolfPlay() {
	const { apiToken, game, player, gameState } = useLoaderData<typeof loader>();

	const store = useMemo(() => {
		void game.game_id;
		void player.user_id;
		return createStore();
	}, [game.game_id, player.user_id]);

	return (
		<JotaiProvider store={store}>
			<ApiClientContext.Provider value={createApiClient(apiToken)}>
				<GolfPlayApp
					key={game.game_id}
					game={game}
					player={player}
					initialGameState={gameState}
				/>
			</ApiClientContext.Provider>
		</JotaiProvider>
	);
}
