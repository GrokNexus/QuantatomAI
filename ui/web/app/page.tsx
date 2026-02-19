"use client";

import React, { useState } from 'react';
import { GridCanvas } from '../components/GridCanvas';
// import { useGridQuery } from '../grid-engine/query/useGridQuery'; // TODO: Enable when service is running

export default function Home() {
    // const { rowCount, isLoading, error } = useGridQuery("VIEW-001");
    const [formula, setFormula] = useState("Sum([Revenue])");

    return (
        <main style={{ height: '100vh', width: '100vw', display: 'flex', flexDirection: 'column' }}>
            <h1 style={{ color: 'white', position: 'absolute', top: '20px', left: '20px', zIndex: 10, margin: 0 }}>
                QuantatomAI Grid (WebGPU)
            </h1>

            {/* Layer 7.3: Formula Editor Stub */}
            <div style={{ padding: '20px', backgroundColor: '#252526', color: '#fff', marginTop: '60px' }}>
                <input
                    value={formula}
                    onChange={(e) => setFormula(e.target.value)}
                    style={{ width: '100%', padding: '8px', backgroundColor: '#3c3c3c', border: '1px solid #555', color: '#fff' }}
                />
            </div>

            <div style={{ flex: 1, position: 'relative' }}>
                <GridCanvas />
            </div>
        </main>
    );
}
