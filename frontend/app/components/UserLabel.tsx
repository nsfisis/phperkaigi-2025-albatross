type Props = {
	label: string;
};

export default function UserLabel({ label }: Props) {
	return (
		<span className="bg-sky-700 text-sky-50 rounded-lg p-3 text-sm">
			{label}
		</span>
	);
}
