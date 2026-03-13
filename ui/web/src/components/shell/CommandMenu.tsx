"use client";

import React, { useState, useEffect, useRef, useMemo } from 'react';
import Fuse from 'fuse.js';

// Mock routing & action data for Fuse
const ACTIONS = [
    { id: '1', title: 'Create New Grid', category: 'Grid Operations', action: 'CREATE_GRID' },
    { id: '2', title: 'Open Revenue Model', category: 'Navigation', action: 'NAV_REVENUE' },
    { id: '3', title: 'Open Headcount Planning', category: 'Navigation', action: 'NAV_HC' },
    { id: '4', title: 'Toggle Light Theme', category: 'Preferences', action: 'THEME_LIGHT' },
    { id: '5', title: 'Toggle Dark Theme', category: 'Preferences', action: 'THEME_DARK' },
    { id: '6', title: 'Drop Q2 Forecast', category: 'Destructive', action: 'DROP_FORECAST' },
    { id: '7', title: 'Sync with Cloud', category: 'System', action: 'SYNC' },
];

export const CommandMenu: React.FC = () => {
    const [isOpen, setIsOpen] = useState(false);
    const [query, setQuery] = useState('');
    const inputRef = useRef<HTMLInputElement>(null);

    // Law 4.1: Omni-Summon Logic
    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
                e.preventDefault();
                setIsOpen((open) => !open);
            }
            if (e.key === 'Escape') {
                setIsOpen(false);
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, []);

    useEffect(() => {
        if (isOpen && inputRef.current) {
            inputRef.current.focus();
            setQuery(''); // Reset on open
        }
    }, [isOpen]);

    // Law 4.2: Instant Search & Action via Fuse.js
    const fuse = useMemo(() => new Fuse(ACTIONS, {
        keys: ['title', 'category'],
        threshold: 0.3,
    }), []);

    const results = query ? fuse.search(query).map(r => r.item) : ACTIONS.slice(0, 5); // Default to recent mostly

    if (!isOpen) return null;

    return (
        // Overlay Backdrop: z-modal-bg (60)
        <div style={{
            position: 'fixed',
            top: 0, left: 0, right: 0, bottom: 0,
            backgroundColor: 'rgba(0, 0, 0, 0.4)',
            backdropFilter: 'blur(8px)',
            zIndex: 'var(--z-modal-bg)',
            display: 'flex',
            alignItems: 'flex-start',
            justifyContent: 'center',
            paddingTop: '15vh'
        }}>
            {/* Modal Payload: z-modal (70) */}
            <div
                style={{
                    width: '100%',
                    maxWidth: '600px',
                    backgroundColor: 'var(--color-surface-elevate)',
                    borderRadius: 'var(--radius-lg)',
                    boxShadow: 'var(--shadow-floating)',
                    border: '1px solid var(--glass-border-color)',
                    overflow: 'hidden',
                    display: 'flex',
                    flexDirection: 'column',
                    zIndex: 'var(--z-modal)',
                    animation: 'cmdk-slide-down var(--duration-fluid) var(--easing-spring)'
                }}
                onClick={(e) => e.stopPropagation()} // Prevent closing when clicking modal
            >
                <div style={{
                    padding: 'var(--space-3) var(--space-4)',
                    borderBottom: '1px solid var(--glass-border-color)'
                }}>
                    <input
                        ref={inputRef}
                        type="text"
                        placeholder="Type a command or search... (e.g. 'Revenue')"
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                        style={{
                            width: '100%',
                            backgroundColor: 'transparent',
                            border: 'none',
                            outline: 'none',
                            color: 'var(--color-text-main)',
                            font: 'var(--text-h4)'
                        }}
                    />
                </div>

                <div style={{
                    maxHeight: '400px',
                    overflowY: 'auto',
                    padding: 'var(--space-2)'
                }}>
                    {results.length === 0 ? (
                        <div style={{ padding: 'var(--space-4)', color: 'var(--color-text-dim)', textAlign: 'center', font: 'var(--text-body)' }}>
                            No results found.
                        </div>
                    ) : (
                        results.map((item, index) => (
                            <div
                                key={item.id}
                                style={{
                                    padding: 'var(--space-3) var(--space-4)',
                                    display: 'flex',
                                    justifyContent: 'space-between',
                                    alignItems: 'center',
                                    cursor: 'pointer',
                                    borderRadius: 'var(--radius-md)',
                                    color: item.category === 'Destructive' ? 'var(--color-variance-negative)' : 'var(--color-text-main)',
                                }}
                                onMouseEnter={(e) => {
                                    e.currentTarget.style.backgroundColor = 'var(--glass-border-color)';
                                }}
                                onMouseLeave={(e) => {
                                    e.currentTarget.style.backgroundColor = 'transparent';
                                }}
                                onClick={() => {
                                    // Simulated execution
                                    if (item.category === 'Destructive') {
                                        // Law 5: Undo-able Actions
                                        alert(`[Undo Toast]: Action "${item.title}" executed. [Undo]`);
                                    } else {
                                        alert(`Executed: ${item.action}`);
                                    }
                                    setIsOpen(false);
                                }}
                            >
                                <span style={{ font: 'var(--text-body)' }}>{item.title}</span>
                                <span style={{ font: 'var(--text-micro)', color: 'var(--color-text-dim)' }}>{item.category}</span>
                            </div>
                        ))
                    )}
                </div>

                <style dangerouslySetInnerHTML={{
                    __html: `
                    @keyframes cmdk-slide-down {
                        from { opacity: 0; transform: translateY(-20px) scale(0.98); }
                        to { opacity: 1; transform: translateY(0) scale(1); }
                    }
                `}} />
            </div>

            {/* Close layer explicitly on background click */}
            <div
                style={{ position: 'absolute', inset: 0, zIndex: -1 }}
                onClick={() => setIsOpen(false)}
            />
        </div>
    );
};
