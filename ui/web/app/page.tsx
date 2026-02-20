"use client";

import React, { useState, useMemo } from 'react';
import { GridCanvas } from '../components/GridCanvas';
import { ChartCanvas } from '../components/ChartCanvas';
import { tableFromArrays } from 'apache-arrow';

export default function Home() {
    const [view, setView] = useState<'grid' | 'chart'>('grid');
    const [formula, setFormula] = useState("Sum([Revenue])");

    // Mock Data for "Visual Intelligence" proof
    // This replicates the Arrow stream schema we expect from Layer 6
    const mockData = useMemo(() => {
        return tableFromArrays({
            region: ["North America", "EMEA", "APAC", "LATAM", "India", "Oceania"],
            revenue: [4500000, 3200000, 5100000, 1800000, 2900000, 1100000],
            profit: [1200000, 800000, 1500000, 400000, 750000, 220000]
        });
    }, []);

    return (
        <div style={{
            flex: 1,
            width: '100%',
            height: '100%',
            display: 'flex',
            flexDirection: 'column',
            backgroundColor: 'transparent',
            color: 'var(--foreground-rgb)',
        }}>
            {/* Header: Platform Controls */}
            {/* Sub-Header: Content Controls (Glassmorphic) */}
            <div style={{
                padding: '16px 24px',
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                backgroundColor: 'rgba(100, 116, 139, 0.05)',
                borderBottom: '1px solid var(--glass-border)',
                backdropFilter: 'blur(12px)',
                zIndex: 20
            }}>
                <div style={{ display: 'flex', backgroundColor: 'rgba(100, 116, 139, 0.1)', borderRadius: '12px', padding: '6px', border: '1px solid var(--glass-border)', boxShadow: 'inset 0 2px 4px rgba(0,0,0,0.05)' }}>
                    <button
                        onClick={() => setView('grid')}
                        style={{
                            padding: '8px 24px',
                            borderRadius: '8px',
                            fontSize: '13px',
                            fontWeight: '600',
                            border: 'none',
                            cursor: 'pointer',
                            transition: 'all 0.2s',
                            backgroundColor: view === 'grid' ? '#3b82f6' : 'transparent',
                            color: view === 'grid' ? 'white' : 'inherit',
                            opacity: view === 'grid' ? 1 : 0.7,
                            boxShadow: view === 'grid' ? '0 4px 12px rgba(59, 130, 246, 0.3)' : 'none'
                        }}
                    >
                        Projector Grid
                    </button>
                    <button
                        onClick={() => setView('chart')}
                        style={{
                            padding: '8px 24px',
                            borderRadius: '8px',
                            fontSize: '13px',
                            fontWeight: '600',
                            border: 'none',
                            cursor: 'pointer',
                            transition: 'all 0.2s',
                            backgroundColor: view === 'chart' ? '#3b82f6' : 'transparent',
                            color: view === 'chart' ? 'white' : 'inherit',
                            opacity: view === 'chart' ? 1 : 0.7,
                            boxShadow: view === 'chart' ? '0 4px 12px rgba(59, 130, 246, 0.3)' : 'none'
                        }}
                    >
                        Visualizer
                    </button>
                </div>
            </div>

            {/* Formula Bar: Layer 7.3 Integration */}
            <section style={{ padding: '16px 24px', zIndex: 10 }}>
                <div style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '12px',
                    backgroundColor: 'var(--glass-bg)',
                    backdropFilter: 'blur(16px)',
                    border: '1px solid var(--glass-border)',
                    boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
                    borderRadius: '12px',
                    padding: '10px 16px'
                }}>
                    <span style={{ color: '#3b82f6', fontWeight: 'bold', fontSize: '16px' }}>ƒx</span>
                    <input
                        value={formula}
                        onChange={(e) => setFormula(e.target.value)}
                        style={{
                            width: '100%',
                            backgroundColor: 'transparent',
                            border: 'none',
                            outline: 'none',
                            color: 'inherit',
                            fontSize: '14px',
                            fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace'
                        }}
                        placeholder="Enter formula..."
                    />
                </div>
            </section>

            {/* Content Area: Holographic Projection + AI Analyst */}
            <div style={{ display: 'flex', flex: 1, overflow: 'hidden' }}>
                <section style={{ flex: 1, position: 'relative', overflow: 'hidden', padding: '32px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                    {view === 'grid' ? (
                        <div style={{
                            width: '100%',
                            height: '100%',
                            borderRadius: '16px',
                            border: '1px solid var(--glass-border)',
                            backgroundColor: 'rgba(15, 23, 42, 0.4)',
                            backdropFilter: 'blur(16px)',
                            overflow: 'hidden',
                            position: 'relative',
                            boxShadow: 'var(--glass-shadow)',
                            display: 'flex',
                            flexDirection: 'column'
                        }}>
                            <GridCanvas data={mockData} />
                            <div style={{ position: 'absolute', left: 0, right: 0, bottom: '24px', display: 'flex', justifyContent: 'center', pointerEvents: 'none' }}>
                                <div style={{
                                    backgroundColor: 'rgba(0, 0, 0, 0.6)',
                                    backdropFilter: 'blur(12px)',
                                    padding: '8px 16px',
                                    borderRadius: '9999px',
                                    border: '1px solid rgba(255, 255, 255, 0.1)',
                                    fontSize: '11px',
                                    textTransform: 'uppercase',
                                    letterSpacing: '0.1em',
                                    color: '#e2e8f0',
                                    fontWeight: '700',
                                    boxShadow: '0 4px 12px rgba(0,0,0,0.2)'
                                }}>
                                    <span style={{ color: '#10b981', marginRight: '8px', textShadow: '0 0 8px #10b981' }}>●</span>
                                    WebGPU Projection Active: 120 FPS
                                </div>
                            </div>
                        </div>
                    ) : (
                        <div style={{ width: '100%', maxWidth: '1200px' }}>
                            <ChartCanvas data={mockData} type="bar" />
                        </div>
                    )}
                </section>
            </div>
        </div>
    );
}
