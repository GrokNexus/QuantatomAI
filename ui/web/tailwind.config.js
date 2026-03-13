/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./app/**/*.{js,ts,jsx,tsx,mdx}",
        "./src/**/*.{js,ts,jsx,tsx,mdx}",
        "./components/**/*.{js,ts,jsx,tsx,mdx}",
    ],
    theme: {
        extend: {
            colors: {
                'bg-primary': 'var(--color-bg-primary)',
                'surface-base': 'var(--color-surface-base)',
                'accent': 'var(--color-accent)',
                'success': 'var(--color-success)',
                'warning': 'var(--color-warning)',
                'error': 'var(--color-error)',
                'text-primary': 'var(--color-text-primary)',
                'text-muted': 'var(--color-text-muted)',
                'border': 'var(--color-border)',
            },
        },
    },
    plugins: [],
}
