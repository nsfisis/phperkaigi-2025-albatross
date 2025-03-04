import { atom } from "jotai";
import type { components } from "../api/schema";

const gameStartedAtAtom = atom<number | null>(null);
export const setGameStartedAtAtom = atom(null, (_, set, value: number | null) =>
	set(gameStartedAtAtom, value),
);

export type GameStateKind = "waiting" | "starting" | "gaming" | "finished";
type LatestGameState = components["schemas"]["LatestGameState"];
type RankingEntry = components["schemas"]["RankingEntry"];

export const gameStateKindAtom = atom<GameStateKind>((get) => {
	const startedAt = get(gameStartedAtAtom);
	if (!startedAt) {
		return "waiting";
	}

	const durationSeconds = get(durationSecondsAtom);
	const finishedAt = startedAt + durationSeconds;
	const now = get(currentTimestampAtom);
	if (now < startedAt) {
		return "starting";
	} else if (now < finishedAt) {
		return "gaming";
	} else {
		return "finished";
	}
});

const currentTimestampAtom = atom(0);
export const setCurrentTimestampAtom = atom(null, (_, set) =>
	set(currentTimestampAtom, Math.floor(Date.now() / 1000)),
);

const durationSecondsAtom = atom<number>(0);
export const setDurationSecondsAtom = atom(null, (_, set, value: number) =>
	set(durationSecondsAtom, value),
);

export const startingLeftTimeSecondsAtom = atom<number | null>((get) => {
	const startedAt = get(gameStartedAtAtom);
	if (startedAt === null) {
		return null;
	}
	const currentTimestamp = get(currentTimestampAtom);
	return Math.max(0, startedAt - currentTimestamp);
});

export const gamingLeftTimeSecondsAtom = atom<number | null>((get) => {
	const startedAt = get(gameStartedAtAtom);
	if (startedAt === null) {
		return null;
	}
	const durationSeconds = get(durationSecondsAtom);
	const finishedAt = startedAt + durationSeconds;
	const currentTimestamp = get(currentTimestampAtom);
	return Math.min(durationSeconds, Math.max(0, finishedAt - currentTimestamp));
});

const rankingAtom = atom<RankingEntry[]>([]);
export const setRankingAtom = atom(null, (_, set, value: RankingEntry[]) => {
	set(rankingAtom, value);
});

const rawLatestGameStatesAtom = atom<{
	[key: string]: LatestGameState | undefined;
}>({});
export const latestGameStatesAtom = atom((get) => get(rawLatestGameStatesAtom));
export const setLatestGameStatesAtom = atom(
	null,
	(_, set, value: { [key: string]: LatestGameState | undefined }) => {
		set(rawLatestGameStatesAtom, value);
	},
);
