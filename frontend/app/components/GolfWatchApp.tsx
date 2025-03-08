import { useAtomValue, useSetAtom } from "jotai";
import { useContext, useEffect, useState } from "react";
import { useTimer } from "react-use-precision-timer";
import {
	ApiAuthTokenContext,
	apiGetGame,
	apiGetGameWatchLatestStates,
	apiGetGameWatchRanking,
} from "../api/client";
import type { components } from "../api/schema";
import {
	gameStateKindAtom,
	setCurrentTimestampAtom,
	setGameStartedAtAtom,
	setLatestGameStatesAtom,
	setRankingAtom,
} from "../states/watch";
import GolfWatchAppGaming from "./GolfWatchApps/GolfWatchAppGaming";
import GolfWatchAppStarting from "./GolfWatchApps/GolfWatchAppStarting";
import GolfWatchAppWaiting from "./GolfWatchApps/GolfWatchAppWaiting";

type Game = components["schemas"]["Game"];

export type Props = {
	game: Game;
};

export default function GolfWatchApp({ game }: Props) {
	const apiAuthToken = useContext(ApiAuthTokenContext);

	const gameStateKind = useAtomValue(gameStateKindAtom);
	const setGameStartedAt = useSetAtom(setGameStartedAtAtom);
	const setCurrentTimestamp = useSetAtom(setCurrentTimestampAtom);
	const setLatestGameStates = useSetAtom(setLatestGameStatesAtom);
	const setRanking = useSetAtom(setRankingAtom);

	useTimer({ delay: 1000, startImmediately: true }, setCurrentTimestamp);

	const playerA = game.main_players[0]!;
	const playerB = game.main_players[1]!;

	const playerProfileA = {
		id: playerA.user_id,
		displayName: playerA.display_name,
		iconPath: playerA.icon_path ?? null,
	};
	const playerProfileB = {
		id: playerB.user_id,
		displayName: playerB.display_name,
		iconPath: playerB.icon_path ?? null,
	};

	const [isDataPolling, setIsDataPolling] = useState(false);

	useEffect(() => {
		if (isDataPolling) {
			return;
		}
		const timerId = setInterval(async () => {
			if (isDataPolling) {
				return;
			}
			setIsDataPolling(true);

			try {
				if (gameStateKind === "waiting") {
					const { game: g } = await apiGetGame(apiAuthToken, game.game_id);
					if (g.started_at != null) {
						setGameStartedAt(g.started_at);
					}
				} else if (gameStateKind === "gaming") {
					const { states } = await apiGetGameWatchLatestStates(
						apiAuthToken,
						game.game_id,
					);
					setLatestGameStates(states);
					const { ranking } = await apiGetGameWatchRanking(
						apiAuthToken,
						game.game_id,
					);
					setRanking(ranking);
				}
			} catch (error) {
				console.error(error);
			} finally {
				setIsDataPolling(false);
			}
		}, 1000);

		return () => {
			clearInterval(timerId);
		};
	}, [
		isDataPolling,
		apiAuthToken,
		game.game_id,
		gameStateKind,
		setGameStartedAt,
		setLatestGameStates,
		setRanking,
	]);

	if (gameStateKind === "waiting") {
		return (
			<GolfWatchAppWaiting
				gameDisplayName={game.display_name}
				playerProfileA={playerProfileA}
				playerProfileB={playerProfileB}
			/>
		);
	} else if (gameStateKind === "starting") {
		return <GolfWatchAppStarting gameDisplayName={game.display_name} />;
	} else if (gameStateKind === "gaming" || gameStateKind === "finished") {
		return (
			<GolfWatchAppGaming
				gameDisplayName={game.display_name}
				playerProfileA={playerProfileA}
				playerProfileB={playerProfileB}
				problemTitle={game.problem.title}
				problemDescription={game.problem.description}
				gameResult={null /* TODO */}
			/>
		);
	} else {
		return null;
	}
}
