import BorderedContainerWithCaption from "../BorderedContainerWithCaption";
import CodeBlock from "./CodeBlock";

type Props = {
	title: string;
	description: string;
	sampleCode: string;
};

export default function Problem({ title, description, sampleCode }: Props) {
	return (
		<div className="p-4 flex flex-col gap-4">
			<div className="text-center text-xl font-bold">{title}</div>
			<BorderedContainerWithCaption caption="問題">
				<pre className="text-gray-700 whitespace-pre-wrap break-words">
					{description}
				</pre>
			</BorderedContainerWithCaption>
			<BorderedContainerWithCaption caption="サンプルコード">
				<CodeBlock code={sampleCode} language="php" />
			</BorderedContainerWithCaption>
		</div>
	);
}
