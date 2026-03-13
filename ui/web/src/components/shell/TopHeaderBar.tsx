"use client";

import React, { useMemo } from 'react';
import { useShellStore } from '@/store/useShellStore';
import {
    GitBranch,
    FlaskConical,
    BrainCircuit,
    ChevronDown,
    WifiOff,
    Sun,
    Moon,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useTheme } from '@/context/ThemeContext';

export const TopHeaderBar: React.FC<{ onSidebarToggle?: () => void }> = ({ onSidebarToggle }) => {
    const { activeBranch, isOffline, activeApplication } = useShellStore();
    const { theme } = useTheme();

    const isSandbox = activeBranch !== 'main';

    // Memoized class computation (Law 11 – Performance)
    const headerClasses = useMemo(() => cn(
        "glass sticky top-0 z-[var(--z-shell)] backdrop-blur-lg border-b",
        "border-[var(--color-border)] px-[var(--space-6)] py-[var(--space-3)]",
        "flex items-center justify-between transition-all duration-[var(--motion-fluid)] ease-spring",
        isOffline && "border-[var(--color-error)] bg-[var(--color-error)]/5",
        isSandbox && "border-[var(--color-warning)] bg-[var(--color-warning)]/5"
    ), [isOffline, isSandbox]);

    const branchClasses = useMemo(() => cn(
        "flex items-center gap-[var(--space-2)] px-[var(--space-3)] py-[var(--space-1.5)]",
        "rounded-full text-sm font-medium transition-all duration-[var(--motion-swift)]",
        "hover:shadow-[var(--glass-glow)] hover:-translate-y-[1px]",
        isSandbox
            ? "bg-[var(--color-warning)]/10 text-[var(--color-warning)] border border-[var(--color-warning)]/30"
            : "glass-button px-[var(--space-3)] py-[var(--space-1.5)] border-0"
    ), [isSandbox]);

    return (
        <header
            className={headerClasses}
            role="banner"
            aria-label="Application header and global controls"
        >
            {/* Left Section: Logo + Branch Visualizer */}
            <div className="flex items-center gap-[var(--space-4)]">
                {/* Logo – Locked size/position (Rule 13) */}
                <div className="flex items-center gap-[var(--space-2)]">
                    <div
                        className="w-8 h-8 rounded-lg bg-gradient-to-br from-[var(--color-accent)] to-[var(--color-ai-brain)] flex items-center justify-center text-white font-bold text-xl"
                        aria-hidden="true"
                    >
                        Q
                    </div>
                    <span className="font-semibold text-lg text-[var(--color-text-primary)]">
                        QuantAtomAI
                    </span>
                </div>

                {/* Branch Switcher (Law 7 – Contextual) */}
                <button
                    className={branchClasses}
                    aria-label={`Current branch: ${activeBranch}${isSandbox ? ' (Sandbox)' : ''}`}
                    aria-describedby="branch-status"
                >
                    {isSandbox ? <FlaskConical size={16} /> : <GitBranch size={16} />}
                    <span className="font-medium">{activeBranch.toUpperCase()}</span>
                    <ChevronDown size={14} className="opacity-70" />
                </button>

                {/* Offline Indicator */}
                {isOffline && (
                    <div
                        className="flex items-center gap-[var(--space-2)] px-[var(--space-3)] py-[var(--space-1.5)] rounded-full bg-[var(--color-error)]/10 text-[var(--color-error)] text-sm font-medium border border-[var(--color-error)]/30"
                        id="branch-status"
                        role="status"
                        aria-live="polite"
                    >
                        <WifiOff size={14} />
                        Offline Mode
                    </div>
                )}
            </div>

            {/* Right Section: AI Co-Pilot + Profile */}
            <div className="flex items-center gap-[var(--space-3)]">
                {/* Theme Switcher */}
                <button
                    className="glass-button p-2"
                    aria-label="Toggle theme"
                >
                    <div className="relative w-4 h-4 overflow-hidden rounded-full">
                        <Sun className="absolute w-full h-full text-[var(--color-warning)] transition-transform duration-500 hover:rotate-90" />
                    </div>
                </button>

                {/* Fluxion AI Co-Pilot Orb (Law 8 – Motion With Meaning) */}
                <button
                    className="relative flex items-center justify-center w-10 h-10 rounded-full glass hover:scale-110 shadow-[var(--glass-glow)] transition-transform duration-[var(--motion-swift)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--color-accent)] focus-visible:ring-offset-2"
                    aria-label="Open Fluxion AI Co-Pilot"
                >
                    <BrainCircuit size={20} className="text-[var(--color-ai-brain)]" />
                    <span className="absolute -top-1 -right-1 w-4 h-4 bg-[var(--color-ai-brain)] rounded-full animate-pulse shadow-[var(--glass-glow)]" aria-hidden="true" />
                </button>

                {/* Profile Avatar */}
                <button
                    className="relative w-10 h-10 rounded-full overflow-visible glass hover:scale-110 hover:shadow-[var(--glass-glow)] ring-2 ring-white/10 transition-all duration-[var(--motion-swift)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--color-accent)] focus-visible:ring-offset-2"
                    aria-label="User profile menu"
                >
                    {/* Placeholder avatar – replace with real image */}
                    <div className="w-full h-full bg-gradient-to-br from-[var(--color-accent)] to-[var(--color-ai-brain)] flex items-center justify-center text-white font-bold rounded-full overflow-hidden">
                        S
                    </div>
                    {/* Status Dot */}
                    <span className="absolute bottom-0 right-0 w-3 h-3 bg-[var(--color-success)] border-2 border-[var(--color-bg-primary)] rounded-full" />
                </button>
            </div>
        </header>
    );
};