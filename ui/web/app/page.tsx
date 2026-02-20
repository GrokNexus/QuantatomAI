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
        <main style={{
            height: '100vh',
            width: '100vw',
            display: 'flex',
            flexDirection: 'column',
            backgroundColor: '#0a0a0a',
            color: '#e5e5e5',
            fontFamily: 'system-ui, -apple-system, sans-serif'
        }}>
            {/* Header: Platform Controls */}
            <header style={{
                height: '64px',
                borderBottom: '1px solid #262626',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '0 24px',
                backgroundColor: 'rgba(23, 23, 23, 0.8)',
                backdropFilter: 'blur(12px)',
                zIndex: 20
            }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
                    <div style={{
                        width: '32px',
                        height: '32px',
                        backgroundColor: '#2563eb',
                        borderRadius: '8px',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontWeight: 'bold',
                        fontSize: '12px'
                    }}>QA</div>
                    <h1 style={{ fontSize: '14px', fontWeight: '600', letterSpacing: '0.05em', textTransform: 'uppercase', color: '#a3a3a3', margin: 0 }}>
                        QuantatomAI Console
                    </h1>
                </div>

                <div style={{ display: 'flex', backgroundColor: 'rgba(38, 38, 38, 0.5)', borderRadius: '8px', padding: '4px', border: '1px solid #404040' }}>
                    <button
                        onClick={() => setView('grid')}
                        style={{
                            padding: '6px 16px',
                            borderRadius: '6px',
                            fontSize: '11px',
                            fontWeight: '600',
                            border: 'none',
                            cursor: 'pointer',
                            transition: 'all 0.2s',
                            backgroundColor: view === 'grid' ? '#3b82f6' : 'transparent',
                            color: view === 'grid' ? 'white' : '#737373'
                        }}
                    >
                        Projector Grid
                    </button>
                    <button
                        onClick={() => setView('chart')}
                        style={{
                            padding: '6px 16px',
                            borderRadius: '6px',
                            fontSize: '11px',
                            fontWeight: '600',
                            border: 'none',
                            cursor: 'pointer',
                            transition: 'all 0.2s',
                            backgroundColor: view === 'chart' ? '#3b82f6' : 'transparent',
                            color: view === 'chart' ? 'white' : '#737373'
                        }}
                    >
                        Visualizer
                    </button>
                </div>
            </header>

            {/* Formula Bar: Layer 7.3 Integration */}
            <section style={{ padding: '16px', borderBottom: '1px solid #262626', backgroundColor: '#171717' }}>
                <div style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '12px',
                    backgroundColor: '#000',
                    border: '1px solid #262626',
                    borderRadius: '8px',
                    padding: '8px'
                }}>
                    <span style={{ color: '#3b82f6', fontWeight: 'bold', fontSize: '14px', paddingLeft: '8px' }}>ƒx</span>
                    <input
                        value={formula}
                        onChange={(e) => setFormula(e.target.value)}
                        style={{
                            width: '100%',
                            backgroundColor: 'transparent',
                            border: 'none',
                            outline: 'none',
                            color: '#fff',
                            fontSize: '14px',
                            fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace'
                        }}
                        placeholder="Enter formula..."
                    />
                </div>
            </section>

            {/* Content Area: Holographic Projection */}
            <section style={{ flex: 1, position: 'relative', overflow: 'hidden', padding: '32px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                {view === 'grid' ? (
                    <div style={{
                        width: '100%',
                        height: '100%',
                        borderRadius: '24px',
                        border: '1px solid #262626',
                        backgroundColor: '#000',
                        overflow: 'hidden',
                        position: 'relative',
                        boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.5)'
                    }}>
                        <GridCanvas data={mockData} />
                        <div style={{ position: 'absolute', left: 0, right: 0, bottom: '24px', display: 'flex', justifyContent: 'center' }}>
                            <div style={{
                                backgroundColor: 'rgba(0, 0, 0, 0.8)',
                                backdropFilter: 'blur(8px)',
                                padding: '10px 20px',
                                borderRadius: '9999px',
                                border: '1px solid rgba(255, 255, 255, 0.1)',
                                fontSize: '10px',
                                textTransform: 'uppercase',
                                letterSpacing: '0.15em',
                                color: '#a3a3a3',
                                fontWeight: 'bold'
                            }}>
                                <span style={{ color: '#22c55e', marginRight: '8px' }}>●</span>
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
        </main>
    );
}
