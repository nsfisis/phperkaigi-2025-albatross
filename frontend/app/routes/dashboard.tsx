import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { Form, useLoaderData } from "@remix-run/react";
import { ensureUserLoggedIn } from "../.server/auth";
import { apiGetGames } from "../api/client";
import BorderedContainer from "../components/BorderedContainer";
import NavigateLink from "../components/NavigateLink";
import UserIcon from "../components/UserIcon";

export const meta: MetaFunction = () => [
	{ title: "Dashboard | PHPerKaigi 2025 Albatross" },
];

export async function loader({ request }: LoaderFunctionArgs) {
	const { user, token } = await ensureUserLoggedIn(request);
	const { games } = await apiGetGames(token);
	return {
		user,
		games,
	};
}

export default function Dashboard() {
	const { user, games } = useLoaderData<typeof loader>()!;

	return (
		<div className="p-6 bg-gray-100 min-h-screen flex flex-col items-center gap-4">
			{user.icon_path && (
				<UserIcon
					iconPath={user.icon_path}
					displayName={user.display_name}
					className="w-24 h-24"
				/>
			)}
			<h1 className="text-2xl font-bold">
				<span className="text-gray-800">{user.display_name}</span>
				<span className="text-gray-500 ml-2">@{user.username}</span>
			</h1>
			<h2 className="text-xl font-semibold text-gray-700">試合</h2>
			<BorderedContainer>
				<div className="px-4">
					{games.length === 0 ? (
						<p>エントリーできる試合はありません</p>
					) : (
						<ul className="divide-y">
							{games.map((game) => (
								<li
									key={game.game_id}
									className="flex justify-between items-center py-3 gap-3"
								>
									<div>
										<span className="font-medium text-gray-800">
											{game.display_name}
										</span>
										<span className="text-sm text-gray-500 ml-2">
											{game.game_type === "multiplayer"
												? " (マルチ)"
												: " (1v1)"}
										</span>
									</div>
									<span>
										<NavigateLink to={`/golf/${game.game_id}/play`}>
											対戦
										</NavigateLink>
										<NavigateLink to={`/golf/${game.game_id}/watch`}>
											観戦
										</NavigateLink>
									</span>
								</li>
							))}
						</ul>
					)}
				</div>
			</BorderedContainer>
			<Form method="post" action="/logout">
				<button
					type="submit"
					className="px-4 py-2 bg-red-500 text-white rounded-sm transition duration-300 hover:bg-red-700 focus:ring-3 focus:ring-red-400 focus:outline-hidden"
				>
					ログアウト
				</button>
			</Form>
			{user.is_admin && (
				<a
					href={
						process.env.NODE_ENV === "development"
							? "http://localhost:8003/phperkaigi/2025/code-battle/admin/dashboard"
							: "/phperkaigi/2025/code-battle/admin/dashboard"
					}
					className="text-lg text-white bg-sky-600 px-4 py-2 rounded-sm transition duration-300 hover:bg-sky-500 focus:ring-3 focus:ring-sky-400 focus:outline-hidden"
				>
					Admin Dashboard
				</a>
			)}
		</div>
	);
}
