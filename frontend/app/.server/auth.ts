import { type JwtPayload, jwtDecode } from "jwt-decode";
import { redirect } from "react-router";
import { Authenticator } from "remix-auth";
import { FormStrategy } from "remix-auth-form";
import { apiLogin } from "../api/client";
import { components } from "../api/schema";
import { createUnstructuredCookie } from "./cookie";
import { cookieOptions, sessionStorage } from "./session";

const authenticator = new Authenticator<string>();

authenticator.use(
	new FormStrategy(async ({ form }) => {
		const username = String(form.get("username"));
		const password = String(form.get("password"));
		return (await apiLogin(username, password)).token;
	}),
	"default",
);

export type User = components["schemas"]["User"];

// This cookie is used to directly store the JWT for the API server.
// Remix's createCookie() returns "structured" cookies, which cannot be reused directly by non-Remix servers.
const tokenCookie = createUnstructuredCookie("albatross_token", cookieOptions);

/**
 * @throws Error on failure
 */
export async function login(request: Request): Promise<never> {
	const jwt = await authenticator.authenticate("default", request);

	const session = await sessionStorage.getSession(
		request.headers.get("cookie"),
	);
	session.set("user", jwt);

	throw redirect("/dashboard", {
		headers: [
			["Set-Cookie", await sessionStorage.commitSession(session)],
			["Set-Cookie", await tokenCookie.serialize(jwt)],
		],
	});
}

export async function logout(request: Request): Promise<never> {
	const session = await sessionStorage.getSession(
		request.headers.get("cookie"),
	);
	throw redirect("/", {
		headers: [
			["Set-Cookie", await sessionStorage.destroySession(session)],
			[
				"Set-Cookie",
				await tokenCookie.serialize("", { maxAge: 0, expires: new Date(0) }),
			],
		],
	});
}

async function getCurrentValidSession(
	request: Request,
): Promise<{ user: User; token: string } | null> {
	const session = await sessionStorage.getSession(
		request.headers.get("cookie"),
	);
	const token = session.get("user");
	if (!token) {
		return null;
	}
	const user = jwtDecode<User & JwtPayload>(token);
	const exp = user.exp;
	if (exp != null && new Date((exp - 3600) * 1000) < new Date()) {
		// If the token will expire in less than an hour, refresh it.
		return null;
	}
	return { user, token };
}

export async function ensureUserLoggedIn(
	request: Request,
): Promise<{ user: User; token: string }> {
	const session = await getCurrentValidSession(request);
	if (!session) {
		throw redirect("/login");
	}
	return session;
}

export async function ensureUserNotLoggedIn(request: Request): Promise<null> {
	const session = await getCurrentValidSession(request);
	if (session) {
		throw redirect("/dashboard");
	}
	return null;
}
