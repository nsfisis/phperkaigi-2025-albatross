import { JSX, useLayoutEffect, useState } from "react";
import { type BundledLanguage, highlight } from "../../highlight";

type Props = {
	code: string;
	language: BundledLanguage;
};

export default function CodeBlock({ code, language }: Props) {
	const [nodes, setNodes] = useState<JSX.Element | null>(null);

	useLayoutEffect(() => {
		highlight(code, language).then(setNodes);
	}, [code, language]);

	return (
		<pre className="h-full w-full p-2 bg-gray-50 rounded-lg border border-gray-300 whitespace-pre-wrap break-words">
			{nodes === null ? <code>{code}</code> : nodes}
		</pre>
	);
}
