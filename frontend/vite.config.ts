import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
	base: "/phperkaigi/2025/code-battle/",
	plugins: [tailwindcss(), reactRouter(), tsconfigPaths()],
});
