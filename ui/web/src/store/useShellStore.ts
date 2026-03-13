"use client";

import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { devtools } from 'zustand/middleware'; // For Redux DevTools
import { immer } from 'zustand/middleware/immer'; // Immutable updates

// Enums for type safety (Rule 13)
export enum DensityType {
    COMPACT = 'compact',
    OPTIMAL = 'optimal',
    COMFORTABLE = 'comfortable',
}

export enum ThemeType {
    MIDNIGHT_NEBULA = 'midnight-nebula',
    ARCTIC_DAWN = 'arctic-dawn',
    EMERALD_VOID = 'emerald-void',
    SOLAR_FLARE = 'solar-flare',
    COSMIC_AMETHYST = 'cosmic-amethyst',
}

export enum ModuleId {
    HOME = 'home',
    GRID = 'grid',
    STUDIO = 'studio',
    WIZARD = 'wizard',
    INTEL = 'intel',
    REPORTS = 'reports',
}

export interface CalendarConfig {
    calendarType: 'Yearly' | 'Quarterly' | 'Monthly' | 'Weekly';
    fiscalStartMonth: string;
    currentFiscalYear: number;
    pastYearsCount: number;
    futureYearsCount: number;
}

export interface ApplicationSchema {
    id: string;
    name: string;
    description: string;
    dimensions: string[];
    calendar: CalendarConfig;
}

// Core state shape
interface ShellState {
    // Density (Law 14 – Adaptive)
    density: DensityType;
    setDensity: (density: DensityType) => void;

    // Theme (Law 19 – Multi-Theme Harmony)
    theme: ThemeType;
    setTheme: (theme: ThemeType) => void;

    // Sidebar (Law 16)
    isSidebarCollapsed: boolean;
    toggleSidebar: () => void;
    setSidebarCollapsed: (collapsed: boolean) => void;

    // Global Context / POV
    activeTenant: string;
    setActiveTenant: (tenant: string) => void;

    globalPOV: Record<string, string>;
    setGlobalPOV: (dimension: string, value: string) => void;
    resetPOV: () => void;

    activeBranch: string;
    setActiveBranch: (branch: string) => void;

    // AI Co-Pilot (Law 10 – Predictability)
    isFluxionOpen: boolean;
    toggleFluxion: () => void;

    // Network State (Law 5 – Forgiveness / Resilience)
    isOffline: boolean;
    setOffline: (status: boolean) => void;

    // Module Router
    activeModule: ModuleId;
    setActiveModule: (moduleId: ModuleId) => void;

    // Active Application Context
    activeApplication: ApplicationSchema | null;
    setActiveApplication: (app: ApplicationSchema | null) => void;
    clearActiveApplication: () => void;
}

// Default values
const defaultState: Partial<ShellState> = {
    density: DensityType.OPTIMAL,
    theme: ThemeType.MIDNIGHT_NEBULA,
    isSidebarCollapsed: true,
    activeTenant: 'QuantAtom Corp',
    globalPOV: {
        Time: 'FY26',
        Scenario: 'Actual',
        Version: 'Working',
        Currency: 'USD',
    },
    activeBranch: 'main',
    isFluxionOpen: false,
    isOffline: false,
    activeModule: ModuleId.HOME,
    activeApplication: null,
};

export const useShellStore = create<ShellState>()(
    devtools( // Law 11 – Metrics & Debugging
        persist(
            (set) => ({
                ...(defaultState as ShellState),

                setDensity: (density) => {
                    set({ density });
                    if (typeof document !== 'undefined') {
                        document.documentElement.style.setProperty(
                            '--density-multiplier',
                            density === DensityType.COMPACT ? '0.75' :
                                density === DensityType.COMFORTABLE ? '1.25' : '1.0'
                        );
                    }
                },

                setTheme: (theme) => {
                    set({ theme });
                    if (typeof document !== 'undefined') {
                        document.documentElement.setAttribute('data-theme', theme);
                    }
                },

                toggleSidebar: () => set(state => ({ isSidebarCollapsed: !state.isSidebarCollapsed })),
                setSidebarCollapsed: (collapsed) => set({ isSidebarCollapsed: collapsed }),

                setActiveTenant: (tenant) => set({ activeTenant: tenant }),

                setGlobalPOV: (dimension, value) =>
                    set(state => ({
                        globalPOV: { ...state.globalPOV, [dimension]: value }
                    })),

                resetPOV: () => set({ globalPOV: defaultState.globalPOV as Record<string, string> }),

                setActiveBranch: (branch) => set({ activeBranch: branch }),

                toggleFluxion: () => set(state => ({ isFluxionOpen: !state.isFluxionOpen })),

                setOffline: (status) => set({ isOffline: status }),

                setActiveModule: (moduleId) => set({ activeModule: moduleId }),

                setActiveApplication: (app) => set({ activeApplication: app }),

                clearActiveApplication: () => set({ activeApplication: null }),
            }),
            {
                name: 'quantatom-shell-storage-v2', // Versioned storage key
                storage: createJSONStorage(() =>
                    typeof window !== 'undefined'
                        ? window.localStorage
                        : {
                            getItem: () => null,
                            setItem: () => { },
                            removeItem: () => { },
                        } as any
                ),
                partialize: (state) => ({
                    density: state.density,
                    isSidebarCollapsed: state.isSidebarCollapsed,
                    activeTenant: state.activeTenant,
                    activeBranch: state.activeBranch,
                    activeModule: state.activeModule,
                    activeApplication: state.activeApplication,
                    // Never persist theme or POV (Law 19 & 10 – Predictability)
                }),
                version: 2, // For migration
                migrate: (persistedState, version) => {
                    if (version === 1) {
                        // Migration logic from v1 if needed
                        return persistedState as ShellState;
                    }
                    return persistedState as ShellState;
                },
                onRehydrateStorage: () => {
                    console.log('[ShellStore] Rehydrated from storage');
                    return (state, error) => {
                        if (error) console.error('[ShellStore] Rehydration failed:', error);
                    };
                },
            }
        )
    )
);

// Optional: Devtools middleware for Redux DevTools
if (process.env.NODE_ENV === 'development') {
    const devtoolsMiddleware = devtools(useShellStore);
    // Replace in production if needed
}