import type { Metadata } from 'next';
import '../styles/globals.css';

export const metadata: Metadata = {
    title: 'QuantatomAI Grid',
    description: 'Ultra-High Performance Grid with WebGPU',
};

export default function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <html lang="en">
            <body style={{ margin: 0, padding: 0, overflow: 'hidden', backgroundColor: '#0a0a0a', color: '#e5e5e5' }}>
                {children}
            </body>
        </html>
    );
}
