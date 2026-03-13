"use client";

import React, { useEffect, useState, useMemo, useCallback, useRef } from 'react';
import { useShellStore } from '@/store/useShellStore';
import {
    LayoutDashboard,
    Grid3x3,
    Blocks,
    BrainCircuit,
    BarChart3,
    Database,
    Folder,
    Maximize2,
    ChevronRight,
    ChevronLeft,
} from 'lucide-react';
import { cn } from '@/lib/utils'; // Must exist!
import { ModuleId } from '@/store/useShellStore';

interface ModuleDef {
    id: string;
    label: string;
    icon: React.ElementType;
}

const BASE_MODULES: ModuleDef[] = [
    { id: 'home', label: 'Home Dashboard', icon: LayoutDashboard },
    { id: 'grid', label: 'Grid Matrix', icon: Grid3x3 },
    { id: 'studio', label: 'App Studio', icon: Blocks },
    { id: 'intel', label: 'Intelligence', icon: BrainCircuit },
    { id: 'reports', label: 'Reports', icon: BarChart3 },
];

export const OmniSidebar: React.FC<{ collapsed: boolean }> = ({ collapsed }) => {
    const {
        isSidebarCollapsed,
        toggleSidebar,
        activeModule,
        setActiveModule,
        activeApplication,
    } = useShellStore();

    const [isHovered, setIsHovered] = useState(false);
    const [clickCounts, setClickCounts] = useState<Record<string, number>>({});
    const hoverTimeout = useRef<NodeJS.Timeout | null>(null);

    // Load telemetry once
    useEffect(() => {
        if (typeof window === 'undefined') return;
        const stored = localStorage.getItem('quantatom-module-telemetry');
        if (stored) {
            try {
                const parsed = JSON.parse(stored) as Record<string, number>;
                setClickCounts(parsed);
            } catch { }
        }
    }, []);

    // Save on change
    useEffect(() => {
        if (typeof window === 'undefined') return;
        localStorage.setItem('quantatom-module-telemetry', JSON.stringify(clickCounts));
    }, [clickCounts]);

    const sortedModules = useMemo(() => {
        let modules = [...BASE_MODULES];
        // Always show all modules for enterprise navigation
        return modules.sort((a, b) => (clickCounts[b.id] || 0) - (clickCounts[a.id] || 0));
    }, [clickCounts]);

    const handleModuleClick = useCallback((id: string) => {
        setClickCounts(prev => ({
            ...prev,
            [id]: (prev[id] || 0) + 1,
        }));
        setActiveModule(id as ModuleId);
    }, [setActiveModule]);

    const handleMouseEnter = useCallback(() => {
        if (!isSidebarCollapsed) return;
        hoverTimeout.current = setTimeout(() => setIsHovered(true), 250);
    }, [isSidebarCollapsed]);

    const handleMouseLeave = useCallback(() => {
        if (hoverTimeout.current) clearTimeout(hoverTimeout.current);
        setIsHovered(false);
    }, []);

    const isExpanded = !isSidebarCollapsed || isHovered;

    return (
        <aside
            onMouseEnter={handleMouseEnter}
            onMouseLeave={handleMouseLeave}
            className={cn(
                "glass relative flex flex-col transition-all duration-[var(--motion-fluid)] ease-spring",
                isExpanded ? "w-[var(--sidebar-expanded)] backdrop-blur-2xl shadow-[var(--shadow-card)]" : "w-[var(--sidebar-collapsed)]",
                "group/sidebar hover:w-[var(--sidebar-expanded)]",
                "border-r border-[var(--color-border)] z-[var(--z-shell)] overflow-hidden",
                "bg-[var(--color-surface-base)]"
            )}
            role="navigation"
            aria-label="Main navigation sidebar"
            // aria-expanded removed
        >
            <div className="flex flex-col gap-[var(--space-1)] p-[var(--space-2)]">
                {sortedModules.map(item => (
                    <button
                        key={item.id}
                        onClick={() => handleModuleClick(item.id)}
                        className={cn(
                            "flex items-center gap-3 px-4 py-3 rounded-lg",
                            "text-[var(--color-text-muted)] hover:text-[var(--color-accent)] hover:bg-[var(--color-surface-base)]/60",
                            "transition-all duration-150 hover:scale-[1.02] hover:shadow-md",
                            "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--color-accent)] focus-visible:ring-offset-2",
                            activeModule === item.id && "bg-[var(--color-accent)]/20 text-[var(--color-accent)] font-semibold scale-105 shadow-[var(--glass-glow)]"
                        )}
                        aria-current={activeModule === item.id ? "page" : undefined}
                        aria-label={item.label}
                    >
                        {item.icon && React.createElement(item.icon, { size: 20, className: "shrink-0" })}
                        <span className="truncate">{item.label}</span>
                    </button>
                ))}
            </div>

            {activeApplication && (
                <div
                    className={cn(
                        "border-t border-[var(--color-border)] flex flex-col transition-opacity duration-[var(--motion-fluid)]",
                        isExpanded ? "opacity-100" : "opacity-0 pointer-events-none h-0"
                    )}
                >
                    <div className="flex items-center gap-[var(--space-4)] px-[var(--space-4)] py-[var(--space-3)] text-[var(--color-accent)] font-semibold">
                        <Database size={20} />
                        <span className="truncate">{activeApplication.name}</span>
                    </div>

                    <div className="flex items-center gap-[var(--space-3)] px-[var(--space-4)] py-[var(--space-2)] ml-[var(--space-8)] text-[var(--color-text-secondary)]">
                        <Folder size={18} className="text-amber-400" />
                        <span className="text-sm">Dimensions</span>
                    </div>

                    <div className="flex flex-col gap-[var(--space-1)] ml-[var(--space-12)] mt-[var(--space-1)] mb-[var(--space-2)] overflow-y-auto max-h-[40vh]">
                        {activeApplication.dimensions.map(dim => (
                            <button
                                key={dim}
                                className="flex items-center gap-[var(--space-2)] px-[var(--space-3)] py-[var(--space-2)] rounded-lg text-[var(--color-text-dim)] bg-transparent border border-transparent hover:border-[var(--glass-border)] hover:bg-[var(--glass-bg)] hover:shadow-[var(--shadow-card)] hover:text-[var(--color-accent)] transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--color-accent)] focus-visible:ring-offset-2"
                                aria-label={`View dimension ${dim}`}
                            >
                                <Maximize2 size={14} />
                                <span className="text-sm truncate">{dim}</span>
                            </button>
                        ))}
                    </div>
                </div>
            )}

            <div className="mt-auto border-t border-[var(--color-border)] p-[var(--space-2)]">
                <button
                    onClick={toggleSidebar}
                    className="flex items-center justify-center w-full py-[var(--space-3)] text-[var(--color-text-muted)] hover:text-[var(--color-accent)] transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--color-accent)] focus-visible:ring-offset-2"
                    aria-label={isSidebarCollapsed ? "Expand sidebar" : "Collapse sidebar"}
                    // aria-expanded removed
                >
                    {isSidebarCollapsed ? <ChevronRight size={18} /> : <ChevronLeft size={18} />}
                </button>
            </div>

            <div aria-live="polite" className="sr-only">
                Sidebar is now {isExpanded ? 'expanded' : 'collapsed'}
            </div>
        </aside>
    );
};