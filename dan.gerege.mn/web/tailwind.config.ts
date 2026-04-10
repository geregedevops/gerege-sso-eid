import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: "#2563eb",
        "primary-light": "#3b82f6",
        surface: "#1e1e2e",
        bg: "#0f0f1a",
      },
    },
  },
  plugins: [],
};
export default config;
