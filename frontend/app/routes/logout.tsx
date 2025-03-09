import type { ActionFunctionArgs } from "react-router";
import { logout } from "../.server/auth";

export async function action({ request }: ActionFunctionArgs) {
	await logout(request);
	return null;
}
