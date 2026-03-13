"use client";

import React, { useState, useEffect } from 'react';

// Law 3: Dimension Meta-Data Panel
export const DimensionInspector: React.FC<{ selectedNodeId?: string }> = ({ selectedNodeId }) => {
    const [isDeleting, setIsDeleting] = useState(false);
    const [countdown, setCountdown] = useState(3);
    const [inheritOverride, setInheritOverride] = useState(false);

    // Reset state on node change (Law 16: cross-fade instead of jump cut)
    const [opacity, setOpacity] = useState(0);

    useEffect(() => {
        setOpacity(0);
        const timer = setTimeout(() => setOpacity(1), 50);
        return () => clearTimeout(timer);
    }, [selectedNodeId]);

    // Law 9: Destructive Gate Sequence
    useEffect(() => {
        let timer: NodeJS.Timeout;
        if (isDeleting && countdown > 0) {
            timer = setTimeout(() => setCountdown(prev => prev - 1), 1000);
        } else if (countdown === 0) {
            // Failsafe auto-reset if they wait too long to confirm
            timer = setTimeout(() => { setIsDeleting(false); setCountdown(3); }, 2000);
        }
        return () => clearTimeout(timer);
    }, [isDeleting, countdown]);

    if (!selectedNodeId) {
        return (
            <div style={{ flex: 1, padding: 'var(--space-4)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--color-text-dim)' }}>
                Select a node in the LTree to inspect properties.
            </div>
        );
    }

    return (
        <div style={{
            flex: 1,
            display: 'flex',
            flexDirection: 'column',
            gap: 'var(--space-4)',
            opacity,
            transition: 'opacity var(--duration-swift) linear',
        }}>
            <div>
                <label style={{ font: 'var(--text-micro)', color: 'var(--color-text-dim)', textTransform: 'uppercase' }}>Node ID</label>
                <div style={{ font: 'var(--text-body)', color: 'var(--color-text-main)', marginTop: '4px' }}>
                    {selectedNodeId}
                </div>
            </div>

            {/* Law 10: Value Inheritance Visuals */}
            <div>
                <label style={{ font: 'var(--text-micro)', color: 'var(--color-text-dim)', letterSpacing: '0.05em' }}>Exchange Rate Override</label>
                <div style={{ position: 'relative', marginTop: '6px' }}>
                    <input
                        type="text"
                        onFocus={() => setInheritOverride(true)}
                        onBlur={(e) => { if (!e.target.value) setInheritOverride(false); }}
                        placeholder={!inheritOverride ? "1.00 (Inherited from Global)" : "Enter custom rate..."}
                        style={{
                            width: '100%',
                            padding: 'var(--space-2) var(--space-3)',
                            backgroundColor: 'rgba(255,255,255,0.02)',
                            border: '1px solid',
                            borderColor: inheritOverride ? 'rgba(245, 158, 11, 0.3)' : 'var(--glass-border-color)',
                            color: inheritOverride ? 'var(--color-text-main)' : 'var(--color-text-dim)',
                            borderRadius: 'var(--radius-sm)',
                            outline: 'none',
                            font: inheritOverride ? 'var(--text-body)' : 'italic var(--text-body)',
                            transition: 'all var(--duration-fluid) var(--easing-spring)',
                        }}
                    />
                    {inheritOverride && (
                        <span style={{ position: 'absolute', right: '12px', top: '50%', transform: 'translateY(-50%)', font: 'var(--text-micro)', color: 'var(--color-warning)' }}>
                            <span className="google-symbols" style={{ fontSize: '14px', verticalAlign: 'middle', marginRight: '4px' }}>edit_note</span>
                            OVERRIDDEN
                        </span>
                    )}
                </div>
            </div>

            {/* Law 9: Destructive Gates */}
            <div style={{ marginTop: 'auto', paddingTop: 'var(--space-4)', borderTop: '1px solid var(--glass-border-color)' }}>
                {!isDeleting ? (
                    <button
                        onClick={() => setIsDeleting(true)}
                        style={{
                            width: '100%',
                            padding: 'var(--space-2)',
                            backgroundColor: 'transparent',
                            border: '1px solid var(--color-variance-negative)',
                            color: 'var(--color-variance-negative)',
                            borderRadius: 'var(--radius-sm)',
                            cursor: 'pointer',
                            transition: 'all var(--duration-snap)'
                        }}
                        onMouseEnter={(e) => { e.currentTarget.style.backgroundColor = 'rgba(239, 68, 68, 0.1)'; }}
                        onMouseLeave={(e) => { e.currentTarget.style.backgroundColor = 'transparent'; }}
                    >
                        Delete Dimension
                    </button>
                ) : (
                    <button
                        onClick={() => {
                            if (countdown === 0) {
                                console.log(`Deleted ${selectedNodeId}`);
                                setIsDeleting(false);
                                setCountdown(3);
                            }
                        }}
                        style={{
                            width: '100%',
                            padding: 'var(--space-2)',
                            backgroundColor: countdown === 0 ? 'var(--color-warning)' : 'transparent',
                            border: '1px solid var(--color-warning)',
                            color: countdown === 0 ? 'black' : 'var(--color-warning)',
                            borderRadius: 'var(--radius-sm)',
                            cursor: countdown === 0 ? 'pointer' : 'not-allowed',
                            fontWeight: countdown === 0 ? 'bold' : 'normal',
                            transition: 'all 0.2s'
                        }}
                    >
                        {countdown > 0 ? `Unlocking in ${countdown}...` : 'CONFIRM ANNIHILATION'}
                    </button>
                )}
            </div>
        </div>
    );
};
