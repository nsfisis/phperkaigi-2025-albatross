import React from "react";

type Props = {
	children: React.ReactNode;
};

export default function TwoColumnLayout({ children }: Props) {
	return (
		<div className="grow grid grid-cols-2 divide-x divide-gray-300">
			{children}
		</div>
	);
}
