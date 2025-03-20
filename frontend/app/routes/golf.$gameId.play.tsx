import { Provider as JotaiProvider, createStore } from "jotai";
import { useMemo } from "react";
import type { LoaderFunctionArgs, MetaFunction } from "react-router";
import { useLoaderData } from "react-router";
import { ensureUserLoggedIn } from "../.server/auth";
import {
	ApiAuthTokenContext,
	apiGetGame,
	apiGetGamePlayLatestState,
} from "../api/client";
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

	const store = useMemo(() => {
		void game.game_id;
		void player.user_id;
		return createStore();
	}, [game.game_id, player.user_id]);

	return (
		<JotaiProvider store={store}>
			<ApiAuthTokenContext.Provider value={apiAuthToken}>
				<GolfPlayApp
					key={game.game_id}
					game={game}
					player={player}
					initialGameState={gameState}
				/>
			</ApiAuthTokenContext.Provider>
		</JotaiProvider>
	);
}
