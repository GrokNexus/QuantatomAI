"use client";

import React, { useRef, useState, useEffect } from 'react';
import Editor, { useMonaco, Monaco } from '@monaco-editor/react';
import { editor } from 'monaco-editor';

interface MonacoFormulaBarProps {
    selectedNodeId?: string;
}

export const MonacoFormulaBar: React.FC<MonacoFormulaBarProps> = ({ selectedNodeId }) => {
    const monaco = useMonaco();
    const [isSaving, setIsSaving] = useState(false);
    const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);

    // Law 2.1: AtomScript Syntax & IntelliSense
    useEffect(() => {
        if (!monaco) return;

        // Register custom language
        monaco.languages.register({ id: 'atomscript' });

        // Define Monarch Tokenizer
        monaco.languages.setMonarchTokensProvider('atomscript', {
            tokenizer: {
                root: [
                    [/\b(LOOKUP|SPREAD|ALLOCATE|SUM|AVG|IF|THEN|ELSE|AND|OR)\b/, 'keyword'],
                    [/\b\d+(\.\d+)?\b/, 'number'],
                    [/\[.*?\]/, 'type.identifier'], // Dimension brackets e.g. [Region.NA]
                    [/[()]/, 'delimiter.parenthesis'],
                    [/,/, 'delimiter'],
                    [/".*?"/, 'string'],
                ]
            }
        });

        // Define UI Glass Theme (Law 1, 19)
        monaco.editor.defineTheme('quantatom-glass', {
            base: 'vs-dark',
            inherit: true,
            rules: [
                { token: 'keyword', foreground: '3B82F6', fontStyle: 'bold' },
                { token: 'number', foreground: '10B981' },
                { token: 'type.identifier', foreground: 'F59E0B' }, // Warning color for brackets
                { token: 'string', foreground: 'A78BFA' }
            ],
            colors: {
                'editor.background': '#0f172a00', // Transparent overlay
                'editor.foreground': '#F8FAFC',
                'editorLineNumber.foreground': '#475569',
                'editorCursor.foreground': '#3B82F6',
                'editor.selectionBackground': '#3B82F640',
                'editor.inactiveSelectionBackground': '#3B82F620',
            }
        });

        monaco.editor.setTheme('quantatom-glass');

    }, [monaco]);

    const handleEditorDidMount = (editorInstance: editor.IStandaloneCodeEditor, monacoInstance: Monaco) => {
        editorRef.current = editorInstance;

        // Law 2.4: Keyboard Escapism (Escape to blur)
        editorInstance.addCommand(monacoInstance.KeyCode.Escape, () => {
            const el = document.activeElement as HTMLElement;
            if (el) el.blur();
        });

        // Law 2.4: Ctrl+Enter to Dispatch
        editorInstance.addCommand(monacoInstance.KeyMod.CtrlCmd | monacoInstance.KeyCode.Enter, () => {
            setIsSaving(true);
            setTimeout(() => setIsSaving(false), 200); // UI Flash
            console.log("Formula Dispatched:", editorInstance.getValue());

            // Mock squiggly bounds check (Law 2.2)
            const value = editorInstance.getValue();
            if (value.includes('ERROR')) {
                const model = editorInstance.getModel();
                if (model) {
                    monacoInstance.editor.setModelMarkers(model, 'owner', [{
                        severity: monacoInstance.MarkerSeverity.Error,
                        message: 'Topological bounds exceeded: Dimension not found in Grid Context.',
                        startLineNumber: 1, endLineNumber: 1, startColumn: 1, endColumn: 6
                    }]);
                }
            } else {
                const model = editorInstance.getModel();
                if (model) monacoInstance.editor.setModelMarkers(model, 'owner', []);
            }
        });
    };

    return (
        <div style={{ height: '100%', width: '100%', position: 'relative', display: 'flex', flexDirection: 'column' }}>
            {/* Header / Fluxion Lens Trigger */}
            <div style={{
                padding: 'var(--space-2) var(--space-4)',
                backgroundColor: 'var(--color-surface-elevate)',
                borderBottom: '1px solid var(--glass-border-color)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                transition: 'background-color var(--duration-snap)'
            }}>
                <span style={{ font: 'var(--text-mono)', color: 'var(--color-primary)' }}>ƒx AtomScript</span>
                <button style={{
                    backgroundColor: 'rgba(16, 185, 129, 0.05)',
                    border: '1px solid rgba(16, 185, 129, 0.2)',
                    color: 'var(--color-ai-brain)',
                    borderRadius: 'var(--radius-sm)',
                    padding: 'var(--space-1) var(--space-2)',
                    cursor: 'pointer',
                    font: 'var(--text-micro)',
                    display: 'flex', alignItems: 'center', gap: '4px',
                    transition: 'all var(--duration-fluid) var(--easing-spring)'
                }}>
                    <span className="google-symbols" style={{ fontSize: '14px' }}>neurology</span> Fluxion Optimize
                </button>
            </div>
            {/* Editor Mount Point */}
            <div style={{
                flex: 1,
                backgroundColor: 'var(--color-surface-base)',
                position: 'relative',
                border: isSaving ? '1px solid var(--color-positive)' : '1px solid transparent', // Flash on save
                transition: 'border var(--duration-snap)'
            }}>
                <Editor
                    height="100%"
                    language="atomscript"
                    theme="quantatom-glass"
                    value={selectedNodeId ? `LOOKUP([${selectedNodeId}], FY26) * 1.05` : "// Select a dimension node to view properties"}
                    options={{
                        minimap: { enabled: false },
                        scrollBeyondLastLine: false,
                        fontSize: 14,
                        fontFamily: "'Google Sans Code', 'JetBrains Mono', 'Fira Code', monospace",
                        lineHeight: 24,
                        padding: { top: 16, bottom: 16 },
                        suggest: { showKeywords: false },
                    }}
                    onMount={handleEditorDidMount}
                />
            </div>
        </div>
    );
};
