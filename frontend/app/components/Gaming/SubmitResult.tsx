import React from "react";
import type { components } from "../../api/schema";
import SubmitStatusLabel from "../SubmitStatusLabel";

type Props = {
	status: components["schemas"]["ExecutionStatus"];
	submitButton?: React.ReactNode;
};

export default function SubmitResult({ status, submitButton }: Props) {
	return (
		<div className="flex flex-col gap-2">
			<div className="flex">
				{submitButton}
				<div className="grow font-bold text-xl text-center">
					<SubmitStatusLabel status={status} />
				</div>
			</div>
		</div>
	);
}
