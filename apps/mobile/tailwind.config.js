/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./App.tsx', './src/**/*.{js,jsx,ts,tsx}'],
  presets: [require('nativewind/preset')],
  theme: {
    extend: {
      colors: {
        // Mirrors the hand-authored brand/text tokens in apps/web/src/index.css.
        // Web's shadcn semantic layer (primary/secondary/card/sidebar/chart, oklch-based)
        // is shadcn/ui-component-specific and doesn't apply here — mobile isn't using shadcn.
        brand: '#4A9FD8',
        'brand-accent': '#2563EB',
        'brand-navy': '#1E3A5F',
        'surface-2': '#FAFAFA',
        'surface-bg': '#F8F9FA',
        'text-primary': '#0A0A0A',
        'text-secondary': '#555555',
        'text-muted': '#888888',
        'text-subtle': '#AAAAAA',
      },
    },
  },
  plugins: [],
};
