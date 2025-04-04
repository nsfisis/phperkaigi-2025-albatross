import type { LoaderFunctionArgs, MetaFunction } from "react-router";
import { ensureUserNotLoggedIn } from "../.server/auth";
import BorderedContainerWithCaption from "../components/BorderedContainerWithCaption";
import NavigateLink from "../components/NavigateLink";

export const meta: MetaFunction = () => [
	{ title: "PHPerKaigi 2025 Albatross" },
];

export async function loader({ request }: LoaderFunctionArgs) {
	await ensureUserNotLoggedIn(request);
	return null;
}

export default function Index() {
	return (
		<div className="min-h-screen bg-sky-600 flex flex-col items-center justify-center gap-y-6">
			<img
				src="/phperkaigi/2025/code-battle/logo.svg"
				alt="PHPerKaigi 2025"
				className="w-64 h-64"
			/>
			<div className="text-center">
				<div className="font-bold text-sky-50 flex flex-col gap-y-2">
					<div className="text-5xl">PHPER CODE BATTLE</div>
				</div>
			</div>
			<div className="mx-2">
				<BorderedContainerWithCaption caption="オンライン予選開催中 (3/21 決勝当日まで)">
					<p className="text-gray-900 max-w-prose">
						PHPer コードバトルは指示された動作をする PHP
						コードをより短く書けた方が勝ち、という 1 対 1
						の対戦コンテンツです。3/21（金）day0 前夜祭では 3/8
						に実施されたオフライン予選と、当日まで開催しているオンライン予選を勝ち抜いたプレイヤーによるトーナメント形式での
						PHPer コードバトルを実施します。ここでは短いコードが正義です！
						可読性も保守性も放り投げた、イベントならではのコードをお楽しみください！
					</p>
				</BorderedContainerWithCaption>
			</div>
			<div>
				<NavigateLink to="/login">ログイン</NavigateLink>
			</div>
		</div>
	);
}
