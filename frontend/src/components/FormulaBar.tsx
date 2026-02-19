import React, { useRef } from 'react';
import Editor, { Monaco } from '@monaco-editor/react';

interface FormulaBarProps {
    value: string;
    onChange: (value: string) => void;
}

export const FormulaBar: React.FC<FormulaBarProps> = ({ value, onChange }) => {
    const monacoRef = useRef<Monaco | null>(null);

    const handleEditorDidMount = (editor: any, monaco: Monaco) => {
        monacoRef.current = monaco;

        // Register AtomScript Language (Stub)
        // In Layer 7.5, we will add full syntax highlighting here.
        monaco.languages.register({ id: 'atomscript' });
        monaco.languages.setMonarchTokensProvider('atomscript', {
            tokenizer: {
                root: [
                    [/\b(Sum|Avg|Min|Max|IF)\b/, "keyword"],
                    [/\[.*?\]/, "type.identifier"], // [Region] dimensions
                    [/\d+/, "number"],
                ]
            }
        });

        // Define Quantatom Theme
        monaco.editor.defineTheme('quantatom-dark', {
            base: 'vs-dark',
            inherit: true,
            rules: [
                { token: 'keyword', foreground: 'FF0055', fontStyle: 'bold' }, // Neon Red
                { token: 'type.identifier', foreground: '00FF99' },            // Neon Green
            ],
            colors: {
                'editor.background': '#1e1e1e',
            }
        });

        monaco.editor.setTheme('quantatom-dark');
    };

    return (
        <div style={{ height: '40px', borderBottom: '1px solid #333', display: 'flex', alignItems: 'center' }}>
            <div style={{ width: '40px', color: '#888', textAlign: 'center', fontSize: '14px' }}>fx</div>
            <Editor
                height="100%"
                defaultLanguage="atomscript"
                value={value}
                onChange={(val) => onChange(val || '')}
                options={{
                    minimap: { enabled: false },
                    lineNumbers: 'off',
                    glyphMargin: false,
                    folding: false,
                    lineDecorationsWidth: 0,
                    lineNumbersMinChars: 0,
                    renderLineHighlight: 'none',
                    scrollbar: { vertical: 'hidden', horizontal: 'hidden' },
                    overviewRulerLanes: 0,
                    hideCursorInOverviewRuler: true,
                    scrollBeyondLastLine: false,
                    automaticLayout: true,
                    fontFamily: 'Inter, monospace',
                    fontSize: 14,
                }}
                onMount={handleEditorDidMount}
            />
        </div>
    );
};
