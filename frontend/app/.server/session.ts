import { createCookieSessionStorage } from "react-router";

export const cookieOptions = {
	sameSite: "lax" as const,
	path: "/",
	httpOnly: true,
	secure: process.env.NODE_ENV === "production",
	secrets: [process.env.ALBATROSS_COOKIE_SECRET ?? "local"],
};

const innerSessionStorage = createCookieSessionStorage({
	cookie: {
		name: "albatross_session",
		...cookieOptions,
	},
});
type InnerSessionStorage = typeof innerSessionStorage;

// This class is used to recover from invalid sessions.
// It may occur if the session had been created before the authentication library was updated.
class RecoverableSessionStorage {
	innerStorage: InnerSessionStorage;

	constructor(innerStorage: InnerSessionStorage) {
		this.innerStorage = innerStorage;
	}

	// If the session is invalid, return a new session.
	// It may occur if the session had been created before the authentication library was updated.
	getSession(...args: Parameters<InnerSessionStorage["getSession"]>) {
		try {
			return this.innerStorage.getSession(...args);
		} catch (e) {
			void e;
			return this.innerStorage.getSession();
		}
	}

	commitSession(...args: Parameters<InnerSessionStorage["commitSession"]>) {
		return this.innerStorage.commitSession(...args);
	}

	destroySession(...args: Parameters<InnerSessionStorage["destroySession"]>) {
		return this.innerStorage.destroySession(...args);
	}
}

export const sessionStorage = new RecoverableSessionStorage(
	innerSessionStorage,
);
