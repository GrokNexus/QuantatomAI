import React, { useState } from 'react';

// Law 5: Contextual Mutability & Intercept
export interface SpreadOptions {
    strategy: 'PRO_RATA' | 'EVEN' | 'HOLD_PUSH';
}

interface TopDownSpreadProps {
    nodeName: string;
    originalValue: number;
    newValue: number;
    onConfirm: (options: SpreadOptions) => void;
    onCancel: () => void;
}

export const TopDownSpread: React.FC<TopDownSpreadProps> = ({
    nodeName, originalValue, newValue, onConfirm, onCancel
}) => {
    const [strategy, setStrategy] = useState<SpreadOptions['strategy']>('PRO_RATA');

    const delta = newValue - originalValue;
    const isPositive = delta > 0;

    return (
        <div style={{
            position: 'absolute',
            top: 0, left: 0, right: 0, bottom: 0,
            backgroundColor: 'rgba(0,0,0,0.5)',
            backdropFilter: 'blur(4px)',
            zIndex: 'var(--z-modal)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
        }}>
            <div style={{
                backgroundColor: 'var(--color-surface-elevate)',
                border: '1px solid var(--glass-border-color)',
                borderRadius: 'var(--radius-lg)',
                padding: 'var(--space-6)',
                width: '400px',
                boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.05)',
                animation: 'slideUp var(--duration-fluid) var(--easing-gravity)'
            }}>
                <div style={{ marginBottom: 'var(--space-4)' }}>
                    <h3 style={{ margin: 0, color: 'var(--color-text-main)', font: 'var(--text-h3)' }}>Top-Down Spread</h3>
                    <p style={{ margin: '4px 0 0', color: 'var(--color-text-dim)', font: 'var(--text-body)' }}>
                        You are mutating a consolidated parent node: <strong>{nodeName}</strong>.
                    </p>
                </div>

                <div style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    padding: 'var(--space-3)',
                    backgroundColor: 'rgba(0,0,0,0.2)',
                    borderRadius: 'var(--radius-md)',
                    marginBottom: 'var(--space-4)'
                }}>
                    <div>
                        <div style={{ font: 'var(--text-micro)', color: 'var(--color-text-dim)', textTransform: 'uppercase' }}>Delta</div>
                        <div style={{
                            font: 'var(--text-mono)',
                            color: isPositive ? 'var(--color-variance-positive)' : 'var(--color-variance-negative)',
                            fontWeight: 'bold'
                        }}>
                            {isPositive ? '+' : ''}{delta.toLocaleString()}
                        </div>
                    </div>
                    <div style={{ fontSize: '24px', color: 'var(--glass-border-color)' }}>→</div>
                    <div style={{ textAlign: 'right' }}>
                        <div style={{ font: 'var(--text-micro)', color: 'var(--color-text-dim)', textTransform: 'uppercase' }}>New Value</div>
                        <div style={{ font: 'var(--text-mono)', color: 'var(--color-text-main)' }}>{newValue.toLocaleString()}</div>
                    </div>
                </div>

                {/* Law 9: Intelligent Distribution Selection */}
                <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
                    <label style={{ font: 'var(--text-micro)', color: 'var(--color-text-dim)', textTransform: 'uppercase' }}>Spread Strategy</label>
                    <button
                        onClick={() => setStrategy('PRO_RATA')}
                        style={{
                            padding: 'var(--space-3)',
                            backgroundColor: strategy === 'PRO_RATA' ? 'rgba(59, 130, 246, 0.05)' : 'rgba(255,255,255,0.01)',
                            border: '1px solid',
                            borderColor: strategy === 'PRO_RATA' ? 'var(--color-primary)' : 'var(--glass-border-color)',
                            color: strategy === 'PRO_RATA' ? 'var(--color-primary)' : 'var(--color-text-dim)',
                            borderRadius: 'var(--radius-sm)',
                            cursor: 'pointer',
                            textAlign: 'left',
                            display: 'flex', alignItems: 'center', gap: 'var(--space-3)',
                            transition: 'all var(--duration-swift) var(--easing-spring)'
                        }}
                        onMouseEnter={(e) => { if (strategy !== 'PRO_RATA') e.currentTarget.style.backgroundColor = 'rgba(255,255,255,0.03)' }}
                        onMouseLeave={(e) => { if (strategy !== 'PRO_RATA') e.currentTarget.style.backgroundColor = 'rgba(255,255,255,0.01)' }}
                    >
                        <span className="google-symbols" style={{ fontSize: '24px' }}>account_tree</span>
                        <div>
                            <strong style={{ display: 'block' }}>Pro-Rata to Base</strong>
                            <div style={{ fontSize: '11px', marginTop: '4px', opacity: 0.8 }}>Distribute proportional to current child values.</div>
                        </div>
                    </button>

                    <button
                        onClick={() => setStrategy('EVEN')}
                        style={{
                            padding: 'var(--space-3)',
                            backgroundColor: strategy === 'EVEN' ? 'rgba(59, 130, 246, 0.05)' : 'rgba(255,255,255,0.01)',
                            border: '1px solid',
                            borderColor: strategy === 'EVEN' ? 'var(--color-primary)' : 'var(--glass-border-color)',
                            color: strategy === 'EVEN' ? 'var(--color-primary)' : 'var(--color-text-dim)',
                            borderRadius: 'var(--radius-sm)',
                            cursor: 'pointer',
                            textAlign: 'left',
                            display: 'flex', alignItems: 'center', gap: 'var(--space-3)',
                            transition: 'all var(--duration-swift) var(--easing-spring)'
                        }}
                        onMouseEnter={(e) => { if (strategy !== 'EVEN') e.currentTarget.style.backgroundColor = 'rgba(255,255,255,0.03)' }}
                        onMouseLeave={(e) => { if (strategy !== 'EVEN') e.currentTarget.style.backgroundColor = 'rgba(255,255,255,0.01)' }}
                    >
                        <span className="google-symbols" style={{ fontSize: '24px' }}>call_split</span>
                        <div>
                            <strong style={{ display: 'block' }}>Evenly Split</strong>
                            <div style={{ fontSize: '11px', marginTop: '4px', opacity: 0.8 }}>Divide the delta evenly across all children.</div>
                        </div>
                    </button>

                    <button
                        onClick={() => setStrategy('HOLD_PUSH')}
                        style={{
                            padding: 'var(--space-3)',
                            backgroundColor: strategy === 'HOLD_PUSH' ? 'rgba(59, 130, 246, 0.05)' : 'rgba(255,255,255,0.01)',
                            border: '1px solid',
                            borderColor: strategy === 'HOLD_PUSH' ? 'var(--color-primary)' : 'var(--glass-border-color)',
                            color: strategy === 'HOLD_PUSH' ? 'var(--color-primary)' : 'var(--color-text-dim)',
                            borderRadius: 'var(--radius-sm)',
                            cursor: 'pointer',
                            textAlign: 'left',
                            display: 'flex', alignItems: 'center', gap: 'var(--space-3)',
                            transition: 'all var(--duration-swift) var(--easing-spring)'
                        }}
                        onMouseEnter={(e) => { if (strategy !== 'HOLD_PUSH') e.currentTarget.style.backgroundColor = 'rgba(255,255,255,0.03)' }}
                        onMouseLeave={(e) => { if (strategy !== 'HOLD_PUSH') e.currentTarget.style.backgroundColor = 'rgba(255,255,255,0.01)' }}
                    >
                        <span className="google-symbols" style={{ fontSize: '24px' }}>publish</span>
                        <div>
                            <strong style={{ display: 'block' }}>Hold Existing, Push Overage</strong>
                            <div style={{ fontSize: '11px', marginTop: '4px', opacity: 0.8 }}>Keep existing base, deposit the entire delta to a "Plug" account.</div>
                        </div>
                    </button>
                </div>

                <div style={{ display: 'flex', gap: 'var(--space-4)', marginTop: 'var(--space-6)' }}>
                    <button
                        onClick={onCancel}
                        style={{ flex: 1, padding: '10px', backgroundColor: 'transparent', border: '1px solid var(--glass-border-color)', color: 'var(--color-text-dim)', borderRadius: 'var(--radius-sm)', cursor: 'pointer' }}
                    >
                        Cancel
                    </button>
                    <button
                        onClick={() => onConfirm({ strategy })}
                        style={{ flex: 1, padding: '10px', backgroundColor: 'var(--color-primary)', border: 'none', color: '#FFF', borderRadius: 'var(--radius-sm)', cursor: 'pointer', fontWeight: 'bold' }}
                    >
                        Confirm Mutation
                    </button>
                </div>
            </div>
        </div>
    );
};
