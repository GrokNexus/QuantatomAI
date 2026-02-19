import type { Metadata } from 'next';

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
            <body style={{ margin: 0, padding: 0, backgroundColor: '#1e1e1e', overflow: 'hidden' }}>{children}</body>
        </html>
    );
}
