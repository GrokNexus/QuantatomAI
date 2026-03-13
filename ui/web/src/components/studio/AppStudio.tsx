"use client";

import React from 'react';
import { Panel, Group as PanelGroup, Separator as PanelResizeHandle } from 'react-resizable-panels';
import { InteractiveTree } from './InteractiveTree';
import { MonacoFormulaBar } from './MonacoFormulaBar';
import { DimensionInspector } from './DimensionInspector';

export const AppStudio: React.FC = () => {
    const [selectedNodeId, setSelectedNodeId] = React.useState<string | undefined>(undefined);

    return (
        <div style={{ height: '100%', width: '100%', display: 'flex' }}>
            {/* Law 4.1: Split Pane Governance */}
            <PanelGroup orientation="horizontal" id="quantatom-app-studio-layout">
                {/* Hierarchy Tree Virtualized Engine */}
                <Panel defaultSize={20} minSize={15} style={{ backgroundColor: 'var(--color-surface-elevate)', display: 'flex', flexDirection: 'column' }}>
                    <div style={{ padding: 'var(--space-4)', borderBottom: '1px solid var(--glass-border-color)' }}>
                        <input
                            placeholder="Filter nodes..."
                            style={{
                                width: '100%', marginTop: 'var(--space-2)',
                                padding: 'var(--space-2)', borderRadius: 'var(--radius-sm)',
                                border: 'none',
                                backgroundColor: 'rgba(255,255,255,0.03)', color: 'white', outline: 'none',
                                transition: 'background-color var(--duration-swift) var(--easing-gravity)'
                            }}
                        />
                    </div>
                    <div style={{ flex: 1, position: 'relative' }}>
                        <InteractiveTree selectedNodeId={selectedNodeId} onNodeSelect={setSelectedNodeId} />
                    </div>
                </Panel>

                <PanelResizeHandle style={{
                    width: '6px',
                    backgroundColor: 'var(--glass-border-color)',
                    cursor: 'col-resize',
                    transition: 'background-color 0.2s'
                }} />

                {/* Intelligence Interface: Monaco Editor */}
                <Panel minSize={40} style={{ display: 'flex', flexDirection: 'column', backgroundColor: 'var(--color-surface-base)' }}>
                    <div style={{ flex: 1, padding: 'var(--space-4)', display: 'flex', flexDirection: 'column' }}>
                        <h2 style={{ font: 'var(--text-h4)', color: 'var(--color-text-main)', margin: '0 0 var(--space-4) 0' }}>AtomScript Properties</h2>

                        <div style={{
                            backgroundColor: 'var(--color-surface-elevate)',
                            borderRadius: 'var(--radius-md)',
                            boxShadow: 'var(--glass-shadow-soft)',
                            overflow: 'hidden',
                            flexShrink: 0,
                            height: '300px' // Increased height to match screenshot spacing
                        }}>
                            <MonacoFormulaBar selectedNodeId={selectedNodeId} />
                        </div>

                        {/* Contextual Inspector Panel */}
                        <div style={{
                            flex: 1,
                            marginTop: 'var(--space-6)',
                            backgroundColor: 'var(--color-surface-elevate)',
                            borderRadius: 'var(--radius-md)',
                            boxShadow: 'var(--glass-shadow-soft)',
                            padding: 'var(--space-4)',
                            display: 'flex',
                            flexDirection: 'column'
                        }}>
                            <DimensionInspector selectedNodeId={selectedNodeId} />
                        </div>
                    </div>
                </Panel>
            </PanelGroup>
        </div>
    );
};
