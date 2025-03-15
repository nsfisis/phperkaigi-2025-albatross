import { useAtom, useAtomValue, useSetAtom } from "jotai";
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
	rankingAtom,
	setCurrentTimestampAtom,
	setGameStartedAtAtom,
	setLatestGameStatesAtom,
} from "../states/watch";
import GolfWatchAppGaming1v1 from "./GolfWatchApps/GolfWatchAppGaming1v1";
import GolfWatchAppGamingMultiplayer from "./GolfWatchApps/GolfWatchAppGamingMultiplayer";
import GolfWatchAppStarting from "./GolfWatchApps/GolfWatchAppStarting";
import GolfWatchAppWaiting1v1 from "./GolfWatchApps/GolfWatchAppWaiting1v1";
import GolfWatchAppWaitingMultiplayer from "./GolfWatchApps/GolfWatchAppWaitingMultiplayer";

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
	const [ranking, setRanking] = useAtom(rankingAtom);

	useTimer({ delay: 1000, startImmediately: true }, setCurrentTimestamp);

	const playerA = game.main_players[0];
	const playerB = game.main_players[1];

	const playerProfileA = playerA
		? {
				id: playerA.user_id,
				displayName: playerA.display_name,
				iconPath: playerA.icon_path ?? null,
			}
		: null;
	const playerProfileB = playerB
		? {
				id: playerB.user_id,
				displayName: playerB.display_name,
				iconPath: playerB.icon_path ?? null,
			}
		: null;

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
		return game.game_type === "1v1" ? (
			<GolfWatchAppWaiting1v1
				gameDisplayName={game.display_name}
				playerProfileA={playerProfileA!}
				playerProfileB={playerProfileB!}
			/>
		) : (
			<GolfWatchAppWaitingMultiplayer gameDisplayName={game.display_name} />
		);
	} else if (gameStateKind === "starting") {
		return <GolfWatchAppStarting gameDisplayName={game.display_name} />;
	} else if (gameStateKind === "gaming" || gameStateKind === "finished") {
		return game.game_type === "1v1" ? (
			<GolfWatchAppGaming1v1
				gameDisplayName={game.display_name}
				playerProfileA={playerProfileA!}
				playerProfileB={playerProfileB!}
				problemTitle={game.problem.title}
				problemDescription={game.problem.description}
				sampleCode={game.problem.sample_code}
				gameResult={null /* TODO */}
			/>
		) : (
			<GolfWatchAppGamingMultiplayer
				gameDisplayName={game.display_name}
				ranking={ranking}
				problemTitle={game.problem.title}
				problemDescription={game.problem.description}
				sampleCode={game.problem.sample_code}
				gameResult={null /* TODO */}
			/>
		);
	} else {
		return null;
	}
}
