import React from "react";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement>;

export default function SubmitButton(props: ButtonProps) {
	return (
		<button
			{...props}
			className="text-lg text-white bg-sky-600 px-4 py-2 rounded transition duration-300 hover:bg-sky-500 focus:ring focus:ring-sky-400 focus:outline-none"
		/>
	);
}
