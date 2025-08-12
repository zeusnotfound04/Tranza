import type { Config } from "tailwindcss";

export default {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--background)",
        foreground: "var(--foreground)",
        // Dynamic TOC Colors (adapted for standard Tailwind)
        'd-bg': 'hsl(var(--d-bg))',
        'd-fg': 'hsl(var(--d-fg))',
        'd-border': 'hsl(var(--d-border))',
        'd-sheet': 'hsl(var(--d-sheet))',
        'd-muted': 'hsl(var(--d-muted))',
      },
    },
  },
  plugins: [],
} satisfies Config;
