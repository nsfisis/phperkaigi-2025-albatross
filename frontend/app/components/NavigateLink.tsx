import { Link, LinkProps } from "@remix-run/react";

export default function NavigateLink(props: LinkProps) {
	return (
		<Link
			{...props}
			className="text-lg text-white bg-sky-600 px-4 py-2 border-2 border-sky-50 rounded transition duration-300 hover:bg-sky-500 focus:ring focus:ring-sky-400 focus:outline-none"
		/>
	);
}
