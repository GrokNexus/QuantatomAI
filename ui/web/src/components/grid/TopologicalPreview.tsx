"use client";

import React, { useEffect, useState } from 'react';

interface TopologicalPreviewProps {
    cellId: string | null;
    formula?: string;
    lineage?: string[];
    position: { x: number; y: number } | null;
}

// Law 11 & 24: Intelligent Hover Previews
export const TopologicalPreview: React.FC<TopologicalPreviewProps> = ({ cellId, formula, lineage, position }) => {
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        let timer: NodeJS.Timeout;
        if (cellId && position) {
            // Law 16: Intent Delay. Don't flash instantly while scrubbing mouse. Wait 1 full second.
            timer = setTimeout(() => setIsVisible(true), 1000);
        } else {
            setIsVisible(false);
        }
        return () => clearTimeout(timer);
    }, [cellId, position]);

    if (!isVisible || !position || !cellId) return null;

    return (
        <div style={{
            position: 'absolute',
            top: position.y + 20, // Offset below cursor
            left: position.x + 20,
            backgroundColor: 'var(--color-surface-elevate)',
            border: '1px solid var(--glass-border-color)',
            boxShadow: 'var(--glass-shadow-soft)',
            borderRadius: 'var(--radius-md)',
            padding: 'var(--space-3)',
            zIndex: 'var(--z-tooltip)', // 80 layer mapping
            minWidth: '250px',
            maxWidth: '350px',
            pointerEvents: 'none', // Do not block canvas interactions
            color: 'white',
            animation: 'fadeIn var(--duration-fluid) var(--easing-spring)'
        }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 'var(--space-2)' }}>
                <span style={{ font: 'var(--text-micro)', color: 'var(--color-text-dim)', letterSpacing: '0.05em' }}>TARGET NODE</span>
                <span className="google-symbols" style={{ fontSize: '16px', color: 'var(--color-primary)' }}>account_tree</span>
            </div>

            <div style={{ font: 'var(--text-body)', color: 'var(--color-text-main)', marginBottom: 'var(--space-2)' }}>
                <strong>{cellId}</strong>
            </div>

            {formula && (
                <div style={{
                    backgroundColor: 'rgba(0,0,0,0.2)',
                    padding: 'var(--space-2)',
                    borderRadius: 'var(--radius-sm)',
                    font: 'var(--text-mono)',
                    color: 'var(--color-positive)',
                    marginBottom: 'var(--space-2)'
                }}>
                    {formula}
                </div>
            )}

            {lineage && lineage.length > 0 && (
                <div style={{ borderTop: '1px solid var(--glass-border-color)', paddingTop: 'var(--space-2)', marginTop: 'var(--space-2)' }}>
                    <div style={{ font: 'var(--text-micro)', color: 'var(--color-text-dim)', marginBottom: '8px' }}>LINEAGE TRACE</div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                        {lineage.map((step, idx) => (
                            <div key={idx} style={{ display: 'flex', alignItems: 'center', gap: '8px', font: 'var(--text-micro)' }}>
                                <span className="google-symbols" style={{ fontSize: '12px', color: 'var(--color-text-dim)' }}>arrow_forward</span>
                                <span>{step}</span>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
};
