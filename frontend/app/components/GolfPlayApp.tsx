import { useAtomValue, useSetAtom } from "jotai";
import { useContext, useEffect, useState } from "react";
import { useTimer } from "react-use-precision-timer";
import { useDebouncedCallback } from "use-debounce";
import {
	ApiAuthTokenContext,
	apiGetGame,
	apiGetGamePlayLatestState,
	apiPostGamePlayCode,
	apiPostGamePlaySubmit,
} from "../api/client";
import type { components } from "../api/schema";
import {
	gameStateKindAtom,
	handleSubmitCodePostAtom,
	handleSubmitCodePreAtom,
	setCurrentTimestampAtom,
	setGameStartedAtAtom,
	setLatestGameStateAtom,
} from "../states/play";
import GolfPlayAppFinished from "./GolfPlayApps/GolfPlayAppFinished";
import GolfPlayAppGaming from "./GolfPlayApps/GolfPlayAppGaming";
import GolfPlayAppStarting from "./GolfPlayApps/GolfPlayAppStarting";
import GolfPlayAppWaiting from "./GolfPlayApps/GolfPlayAppWaiting";

type Game = components["schemas"]["Game"];
type User = components["schemas"]["User"];

type Props = {
	game: Game;
	player: User;
	initialCode: string;
};

export default function GolfPlayApp({ game, player, initialCode }: Props) {
	const apiAuthToken = useContext(ApiAuthTokenContext);

	const gameStateKind = useAtomValue(gameStateKindAtom);
	const setGameStartedAt = useSetAtom(setGameStartedAtAtom);
	const setCurrentTimestamp = useSetAtom(setCurrentTimestampAtom);
	const handleSubmitCodePre = useSetAtom(handleSubmitCodePreAtom);
	const handleSubmitCodePost = useSetAtom(handleSubmitCodePostAtom);
	const setLatestGameState = useSetAtom(setLatestGameStateAtom);

	useTimer({ delay: 1000, startImmediately: true }, setCurrentTimestamp);

	const playerProfile = {
		id: player.user_id,
		displayName: player.display_name,
		iconPath: player.icon_path ?? null,
	};

	const onCodeChange = useDebouncedCallback(async (code: string) => {
		console.log("player:c2s:code");
		if (game.game_type === "1v1") {
			await apiPostGamePlayCode(apiAuthToken, game.game_id, code);
		}
	}, 1000);

	const onCodeSubmit = useDebouncedCallback(async (code: string) => {
		if (code === "") {
			return;
		}
		console.log("player:c2s:submit");
		handleSubmitCodePre();
		await apiPostGamePlaySubmit(apiAuthToken, game.game_id, code);
		handleSubmitCodePost();
	}, 1000);

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
					const { state } = await apiGetGamePlayLatestState(
						apiAuthToken,
						game.game_id,
					);
					setLatestGameState(state);
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
		setLatestGameState,
	]);

	if (gameStateKind === "waiting") {
		return (
			<GolfPlayAppWaiting
				gameDisplayName={game.display_name}
				playerProfile={playerProfile}
			/>
		);
	} else if (gameStateKind === "starting") {
		return <GolfPlayAppStarting gameDisplayName={game.display_name} />;
	} else if (gameStateKind === "gaming") {
		return (
			<GolfPlayAppGaming
				gameDisplayName={game.display_name}
				playerProfile={playerProfile}
				problemTitle={game.problem.title}
				problemDescription={game.problem.description}
				sampleCode={game.problem.sample_code}
				initialCode={initialCode}
				onCodeChange={onCodeChange}
				onCodeSubmit={onCodeSubmit}
			/>
		);
	} else if (gameStateKind === "finished") {
		return <GolfPlayAppFinished />;
	} else {
		return null;
	}
}
