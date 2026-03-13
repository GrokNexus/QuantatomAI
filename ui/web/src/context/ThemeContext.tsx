"use client";

import React, { createContext, useContext } from 'react';
import { useShellStore, ThemeType } from '@/store/useShellStore';

const ThemeProviderContext = createContext<{
    theme: string;
    toggleTheme: () => void;
}>({ theme: 'dark', toggleTheme: () => { } });

export function ThemeProvider({ children }: { children: React.ReactNode }) {
    const { theme, setTheme } = useShellStore();
    return (
        <ThemeProviderContext.Provider value={{
            theme,
            toggleTheme: () => setTheme(theme === ThemeType.MIDNIGHT_NEBULA ? ThemeType.ARCTIC_DAWN : ThemeType.MIDNIGHT_NEBULA)
        }}>
            {children}
        </ThemeProviderContext.Provider>
    );
}

export function useTheme() {
    return useContext(ThemeProviderContext);
}
