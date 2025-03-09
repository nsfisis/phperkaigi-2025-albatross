import { useEffect, useState } from "react";
import { codeToHtml } from "../../shiki.bundle";

type Props = {
	code: string;
	language: string;
};

export default function CodeBlock({ code, language }: Props) {
	const [highlightedCode, setHighlightedCode] = useState<string | null>(null);

	useEffect(() => {
		let isMounted = true;

		(async () => {
			const highlighted = await codeToHtml(code, {
				lang: language,
				theme: "github-light",
			});
			if (isMounted) {
				setHighlightedCode(highlighted);
			}
		})();

		return () => {
			isMounted = false;
		};
	}, [code, language]);

	return (
		<pre className="h-full w-full p-2 bg-gray-50 rounded-lg border border-gray-300 whitespace-pre-wrap break-words">
			{highlightedCode === null ? null : (
				<code dangerouslySetInnerHTML={{ __html: highlightedCode }} />
			)}
		</pre>
	);
}
