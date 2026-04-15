import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: "#10b981",
        "primary-light": "#34d399",
        surface: "#1e1e2e",
        bg: "#0f0f1a",
      },
    },
  },
  plugins: [],
};
export default config;
