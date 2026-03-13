import type { Metadata } from 'next';
import { GlobalShell } from '@/layouts/GlobalShell';
import './globals.css';
import { ThemeProvider } from '@/context/ThemeContext';

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
        <html lang="en" suppressHydrationWarning>
            <head>
                <link rel="preconnect" href="https://fonts.googleapis.com" />
                <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
                <link rel="preload" as="style" href="https://fonts.googleapis.com/css2?family=Google+Sans+Flex:opsz,slnt,wdth,wght,ROND@8..144,-10..0,25..150,400..500,0..100&display=swap" />
                <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Google+Sans+Flex:opsz,slnt,wdth,wght,ROND@8..144,-10..0,25..150,400..500,0..100&display=swap" />
                <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Google+Sans+Code:ital,wght@0,400..700;1,400..700&display=swap" />
            </head>
            <body>
                <ThemeProvider>
                    <GlobalShell>
                        {children}
                    </GlobalShell>
                </ThemeProvider>
            </body>
        </html>
    );
}
