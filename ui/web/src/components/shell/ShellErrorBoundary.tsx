"use client";

import React, { Component, ErrorInfo, ReactNode } from 'react';

interface Props {
    children?: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

export class ShellErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false,
        error: null
    };

    public static getDerivedStateFromError(error: Error): State {
        // Update state so the next render will show the fallback UI.
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error('Uncaught error in Shell Payload:', error, errorInfo);
    }

    public render() {
        if (this.state.hasError) {
            return (
                <div style={{
                    display: 'flex', flexDirection: 'column',
                    alignItems: 'center', justifyContent: 'center', flex: 1,
                    backgroundColor: 'var(--color-surface-base)',
                    color: 'var(--color-text-main)',
                    padding: 'var(--space-8)'
                }}>
                    <div style={{
                        padding: 'var(--space-6)',
                        backgroundColor: 'rgba(239, 68, 68, 0.1)',
                        border: '1px solid var(--color-variance-negative)',
                        borderRadius: 'var(--radius-lg)',
                        display: 'flex', flexDirection: 'column', gap: 'var(--space-4)',
                        maxWidth: '600px', width: '100%',
                        boxShadow: '0 8px 32px rgba(239, 68, 68, 0.2)'
                    }}>
                        <h2 style={{ color: 'var(--color-variance-negative)', margin: 0, font: 'var(--text-h3)' }}>
                            Component Crash Detected (Law 5.2)
                        </h2>
                        <p style={{ font: 'var(--text-body)', color: 'var(--color-text-main)', margin: 0 }}>
                            The payload area encountered an unexpected error. The Global Shell remains active.
                        </p>
                        <pre style={{
                            backgroundColor: 'rgba(0,0,0,0.5)',
                            padding: 'var(--space-3)',
                            borderRadius: 'var(--radius-md)',
                            font: 'var(--text-mono)',
                            color: 'var(--color-text-dim)',
                            overflowX: 'auto'
                        }}>
                            {this.state.error?.toString()}
                        </pre>
                        <button
                            onClick={() => this.setState({ hasError: false })}
                            style={{
                                alignSelf: 'flex-start',
                                padding: 'var(--space-2) var(--space-4)',
                                backgroundColor: 'var(--color-primary)',
                                color: '#fff',
                                border: 'none',
                                borderRadius: 'var(--radius-md)',
                                cursor: 'pointer',
                                font: 'var(--text-body)',
                                fontWeight: 600,
                                transition: 'all var(--duration-snap) var(--easing-linear)'
                            }}
                            onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--color-accent-hover)'}
                            onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'var(--color-primary)'}
                        >
                            Engine Reboot (Retry Render)
                        </button>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}
