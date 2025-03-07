import PHPWasm from "./php-wasm.js";

process.once("message", async ({ code, input }) => {
	const PRELUDE = `
  define('STDIN', fopen('php://stdin', 'r'));
  define('STDOUT', fopen('php://stdout', 'r'));
  define('STDERR', fopen('php://stderr', 'r'));

  `;
	const BUFFER_MAX = 1024;

	let stdinPos = 0; // bytewise
	const stdinBuf = Buffer.from(input);
	let stdoutPos = 0; // bytewise
	const stdoutBuf = Buffer.alloc(BUFFER_MAX);
	let stderrPos = 0; // bytewise
	const stderrBuf = Buffer.alloc(BUFFER_MAX);

	const { ccall } = await PHPWasm({
		stdin: () => {
			if (stdinBuf.length <= stdinPos) {
				return null;
			}
			return stdinBuf.readUInt8(stdinPos++);
		},
		stdout: (asciiCode) => {
			if (asciiCode === null) {
				return; // flush
			}
			if (BUFFER_MAX <= stdoutPos) {
				return; // ignore
			}
			stdoutBuf.writeUInt8(
				asciiCode < 0 ? asciiCode + 256 : asciiCode,
				stdoutPos++,
			);
		},
		stderr: (asciiCode) => {
			if (asciiCode === null) {
				return; // flush
			}
			if (BUFFER_MAX <= stderrPos) {
				return; // ignore
			}
			stderrBuf.writeUInt8(
				asciiCode < 0 ? asciiCode + 256 : asciiCode,
				stderrPos++,
			);
		},
	});

	let err;
	let result;
	try {
		result = ccall("php_wasm_run", "number", ["string"], [PRELUDE + code]);
	} catch (e) {
		err = e;
	}
	if (err) {
		process.send({
			status: "runtime_error",
			stdout: stdoutBuf.subarray(0, stdoutPos).toString(),
			stderr: `${stderrBuf.subarray(0, stderrPos).toString()}\n${err.toString()}`,
		});
	} else {
		process.send({
			status: result === 0 ? "success" : "runtime_error",
			stdout: stdoutBuf.subarray(0, stdoutPos).toString(),
			stderr: stderrBuf.subarray(0, stderrPos).toString(),
		});
	}
});
