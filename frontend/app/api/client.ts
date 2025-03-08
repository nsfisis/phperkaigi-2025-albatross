import createClient from "openapi-fetch";
import { createContext } from "react";
import type { paths } from "./schema";

const apiClient = createClient<paths>({
	baseUrl:
		process.env.NODE_ENV === "development"
			? "http://localhost:8003/phperkaigi/2025/code-battle/api/"
			: "http://api-server/phperkaigi/2025/code-battle/api/",
});

export async function apiPostLogin(username: string, password: string) {
	const { data, error } = await apiClient.POST("/login", {
		body: {
			username,
			password,
		},
	});
	if (error) throw new Error(error.message);
	return data;
}

export async function apiGetGames(token: string) {
	const { data, error } = await apiClient.GET("/games", {
		params: {
			header: { Authorization: `Bearer ${token}` },
		},
	});
	if (error) throw new Error(error.message);
	return data;
}

export async function apiGetGame(token: string, gameId: number) {
	const { data, error } = await apiClient.GET("/games/{game_id}", {
		params: {
			header: { Authorization: `Bearer ${token}` },
			path: { game_id: gameId },
		},
	});
	if (error) throw new Error(error.message);
	return data;
}

export async function apiGetGamePlayLatestState(token: string, gameId: number) {
	const { data, error } = await apiClient.GET(
		"/games/{game_id}/play/latest_state",
		{
			params: {
				header: { Authorization: `Bearer ${token}` },
				path: { game_id: gameId },
			},
		},
	);
	if (error) throw new Error(error.message);
	return data;
}

export async function apiPostGamePlayCode(
	token: string,
	gameId: number,
	code: string,
) {
	const { error } = await apiClient.POST("/games/{game_id}/play/code", {
		params: {
			header: { Authorization: `Bearer ${token}` },
			path: { game_id: gameId },
		},
		body: { code },
	});
	if (error) throw new Error(error.message);
}

export async function apiPostGamePlaySubmit(
	token: string,
	gameId: number,
	code: string,
) {
	const { data, error } = await apiClient.POST("/games/{game_id}/play/submit", {
		params: {
			header: { Authorization: `Bearer ${token}` },
			path: { game_id: gameId },
		},
		body: { code },
	});
	if (error) throw new Error(error.message);
	return data;
}

export async function apiGetGameWatchRanking(token: string, gameId: number) {
	const { data, error } = await apiClient.GET(
		"/games/{game_id}/watch/ranking",
		{
			params: {
				header: { Authorization: `Bearer ${token}` },
				path: { game_id: gameId },
			},
		},
	);
	if (error) throw new Error(error.message);
	return data;
}

export async function apiGetGameWatchLatestStates(
	token: string,
	gameId: number,
) {
	const { data, error } = await apiClient.GET(
		"/games/{game_id}/watch/latest_states",
		{
			params: {
				header: { Authorization: `Bearer ${token}` },
				path: { game_id: gameId },
			},
		},
	);
	if (error) throw new Error(error.message);
	return data;
}

export const ApiAuthTokenContext = createContext<string>("");
