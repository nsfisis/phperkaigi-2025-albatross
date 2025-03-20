import React from "react";

type Props = {
	children: React.ReactNode;
	title: React.ReactNode;
};

export default function TitledColumn({ children, title }: Props) {
	return (
		<div className="p-4 flex flex-col gap-4">
			<div className="text-center text-xl font-bold">{title}</div>
			{children}
		</div>
	);
}
