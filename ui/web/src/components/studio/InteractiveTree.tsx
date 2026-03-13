"use client";

import React, { useState, useMemo } from 'react';
// @ts-ignore
import { List } from 'react-window';
import { useShellStore } from '../../store/useShellStore';

// Initial Mock Data (LTree)
interface TreeNode {
    id: string;
    path: string; // e.g. "Global.NA.US"
    label: string;
    level: number;
    hasChildren: boolean;
    isExpanded: boolean;
}

const buildDimensionTree = (dimensions: string[]): TreeNode[] => {
    return dimensions.map((dim, i) => ({
        id: `dim-${i}`,
        path: `Application.${dim}`,
        label: `${dim} Dimension`,
        level: 0,
        hasChildren: false, // For now, flat list until we build the editor
        isExpanded: false
    }));
};

interface InteractiveTreeProps {
    selectedNodeId?: string;
    onNodeSelect?: (id: string) => void;
}

export const InteractiveTree: React.FC<InteractiveTreeProps> = ({ selectedNodeId, onNodeSelect }) => {
    const { density, activeApplication } = useShellStore();

    // Dynamic binding to user's created application
    const treeData = useMemo(() => {
        if (!activeApplication) return [];
        return buildDimensionTree(activeApplication.dimensions);
    }, [activeApplication]);

    // Law 25: Node Density Toggle Hooks
    const rowHeight = useMemo(() => {
        if (density === 'compact') return 28;
        if (density === 'comfortable') return 44;
        return 36; // optimal 
    }, [density]);

    // Law 20: Virtualized Rendering Core
    const Row = ({ index, style, ariaAttributes }: { index: number, style: React.CSSProperties, ariaAttributes: any }) => {
        const node = treeData[index];
        const isSelected = node.id === selectedNodeId;
        return (
            <div
                className="interactive-node"
                onClick={() => onNodeSelect?.(node.id)}
                style={{
                    ...style,
                    display: 'flex',
                    alignItems: 'center',
                    paddingLeft: `calc(16px + ${node.level * 20}px)`,
                    cursor: 'pointer',
                    color: isSelected ? 'var(--color-primary)' : 'rgba(255,255,255,0.85)',
                    backgroundColor: isSelected ? 'rgba(59, 130, 246, 0.2)' : 'transparent',
                    font: 'var(--text-body)',
                    borderBottom: '1px solid rgba(255,255,255,0.02)',
                    transition: 'background-color 0.1s linear, transform 75ms cubic-bezier(0.175, 0.885, 0.32, 1.1)'
                }}
                onMouseEnter={(e) => {
                    if (!isSelected) e.currentTarget.style.backgroundColor = 'var(--glass-border-color)';
                }}
                onMouseLeave={(e) => {
                    if (!isSelected) e.currentTarget.style.backgroundColor = 'transparent';
                    e.currentTarget.style.transform = 'scale(1)';
                }}
                onMouseDown={(e) => {
                    e.currentTarget.style.transform = 'scale(0.98)';
                }}
                onMouseUp={(e) => {
                    e.currentTarget.style.transform = 'scale(1)';
                }}
            >
                {/* Expand / Collapse Icon */}
                < span className="google-symbols" style={{
                    width: '20px',
                    display: 'inline-flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    color: node.hasChildren ? 'var(--color-primary)' : 'transparent',
                    fontSize: '20px',
                    marginRight: '8px',
                    transition: 'transform var(--duration-swift) var(--easing-gravity)'
                }}>
                    {node.hasChildren ? (node.isExpanded ? 'expand_more' : 'chevron_right') : 'radio_button_unchecked'}
                </span >

                {/* Node Label */}
                < span style={{
                    whiteSpace: 'nowrap',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    marginRight: '8px',
                    flex: '1 1 auto'
                }}>
                    {node.label}
                </span >

                {/* LTree Path debugging ghost */}
                < span style={{
                    paddingRight: '16px',
                    font: 'var(--text-micro)',
                    color: 'rgba(255,255,255,0.4)',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                    flexShrink: 1
                }}>
                    {node.path}
                </span >
            </div >
        );
    };

    return (
        <div style={{ width: '100%', height: '100%', backgroundColor: 'transparent', transform: 'translateZ(0)' }}>
            <List
                style={{ overflowX: 'hidden', height: '100%', width: '100%' }}
                rowCount={treeData.length}
                rowHeight={rowHeight}
                rowComponent={Row as any}
                rowProps={{}}
            />
        </div>
    );
};
