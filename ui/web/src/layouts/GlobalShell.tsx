"use client";

import React, { useEffect, useState, Suspense } from 'react';
import { useShellStore } from '@/store/useShellStore';
import { TopHeaderBar } from '@/components/shell/TopHeaderBar';
import { OmniSidebar } from '@/components/shell/OmniSidebar';
import { CommandMenu } from '@/components/shell/CommandMenu';
import { ShellErrorBoundary } from '@/components/shell/ShellErrorBoundary';
import { Loader2, ChevronLeft, ChevronRight } from 'lucide-react';
import { useTheme } from '@/context/ThemeContext';
import { cn } from '@/lib/utils'; // Your classnames helper

export const GlobalShell: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { theme, density, activeModule } = useShellStore();
    const [mounted, setMounted] = useState(false);
    const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

    // Mount & enforce persisted state (Law 13 + 19)
    useEffect(() => {
        setMounted(true);
        document.documentElement.setAttribute('data-theme', theme);
        document.documentElement.style.setProperty('--density-multiplier',
            density === 'compact' ? '0.75' : density === 'comfortable' ? '1.25' : '1.0'
        );
    }, [theme, density]);

    if (!mounted) {
        return null; // Prevent hydration mismatch
    }

    // Module-based content switching
    let mainContent = null;
    switch (activeModule) {
        case 'grid': {
            const QuantAtomGrid = require('@/components/dashboard/QuantAtomGrid').QuantAtomGrid;
            mainContent = <QuantAtomGrid />;
            break;
        }
        case 'home': {
            mainContent = <div className="p-8"><h1 className="text-3xl font-bold mb-4">Home Dashboard</h1><p>Welcome to QuantAtomAI Enterprise Dashboard.</p></div>;
            break;
        }
        case 'studio': {
            mainContent = <div className="p-8"><h1 className="text-3xl font-bold mb-4">App Studio</h1><p>Build and manage your enterprise apps here.</p></div>;
            break;
        }
        case 'intel': {
            mainContent = <div className="p-8"><h1 className="text-3xl font-bold mb-4">Intelligence</h1><p>AI-driven insights and analytics.</p></div>;
            break;
        }
        case 'reports': {
            mainContent = <div className="p-8"><h1 className="text-3xl font-bold mb-4">Reports</h1><p>Enterprise reporting and exports.</p></div>;
            break;
        }
        default: {
            mainContent = <div className="p-8"><h1 className="text-3xl font-bold mb-4">Module Not Found</h1></div>;
        }
    }

    const content = children ?? mainContent;

    return (
        <div
            className={cn(
                "flex flex-col h-dvh w-dvw bg-[var(--color-bg-primary)] text-[var(--color-text-primary)] overflow-hidden",
                "transition-colors duration-[var(--motion-fluid)]",
                "global-shell-font"
            )}
        >
            <header
                className="glass sticky top-0 z-[var(--z-shell)] backdrop-blur-lg border-b border-[var(--color-border)]"
                role="banner"
                aria-label="Application header"
            >
                <TopHeaderBar
                    onSidebarToggle={() => setSidebarCollapsed(!sidebarCollapsed)}
                />
            </header>
            <div className="flex flex-1 overflow-hidden relative">
                <aside
                    className={cn(
                        "glass fixed left-0 top-[var(--header-height)] bottom-0 h-[calc(100vh-var(--header-height))] flex flex-col transition-[width] duration-[var(--motion-fluid)] ease-spring",
                        sidebarCollapsed ? "w-[var(--sidebar-collapsed)]" : "w-[var(--sidebar-expanded)]",
                        "group/sidebar hover:w-[var(--sidebar-expanded)]",
                        "border-r border-[var(--color-border)] z-[var(--z-shell)] bg-[var(--color-surface-base)]"
                    )}
                    role="navigation"
                    aria-label="Main navigation"
                >
                    <OmniSidebar collapsed={sidebarCollapsed} />
                    <button
                        onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
                        className="mx-auto my-4 glass-button rounded-full transition-opacity"
                        aria-label={sidebarCollapsed ? "Expand sidebar" : "Collapse sidebar"}
                    >
                        {sidebarCollapsed ? <ChevronRight size={20} /> : <ChevronLeft size={20} />}
                    </button>
                </aside>
                <main
                    className={cn(
                        "flex-1 flex flex-col overflow-hidden relative",
                        sidebarCollapsed ? "ml-[var(--sidebar-collapsed)]" : "ml-[var(--sidebar-expanded)]"
                    )}
                    role="main"
                    aria-label="Main content"
                >
                    <ShellErrorBoundary>
                        <Suspense fallback={
                            <div className="flex-1 flex items-center justify-center">
                                <div className="flex flex-col items-center gap-4 text-[var(--color-text-muted)]">
                                    <Loader2 className="h-12 w-12 animate-spin text-[var(--color-accent)]" />
                                    <p className="text-lg">Loading matrix...</p>
                                </div>
                            </div>
                        }>
                            {content}
                        </Suspense>
                    </ShellErrorBoundary>
                </main>
            </div>
            <CommandMenu />
        </div>
    );
};