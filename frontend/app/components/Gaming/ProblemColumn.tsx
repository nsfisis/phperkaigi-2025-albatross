import TitledColumn from "../TitledColumn";
import ProblemColumnContent from "./ProblemColumnContent";

type Props = {
	title: string;
	description: string;
	sampleCode: string;
};

export default function ProblemColumn({
	title,
	description,
	sampleCode,
}: Props) {
	return (
		<TitledColumn title={title}>
			<ProblemColumnContent description={description} sampleCode={sampleCode} />
		</TitledColumn>
	);
}
