import { useEffect, useState, useCallback } from 'react';

export interface CellAddress {
    row: number;
    col: number;
}

export interface GridSelection {
    start: CellAddress;
    end: CellAddress;
}

export const useGridHotkeys = (numRows: number, numCols: number) => {
    const [activeCell, setActiveCell] = useState<CellAddress>({ row: 0, col: 0 });
    const [selection, setSelection] = useState<GridSelection>({
        start: { row: 0, col: 0 },
        end: { row: 0, col: 0 }
    });
    const [isEditing, setIsEditing] = useState(false);
    const [editValue, setEditValue] = useState("");

    const moveCell = useCallback((rowDelta: number, colDelta: number, shift: boolean) => {
        setActiveCell(prev => {
            const nextRow = Math.max(0, Math.min(prev.row + rowDelta, numRows - 1));
            const nextCol = Math.max(0, Math.min(prev.col + colDelta, numCols - 1));
            const nextCell = { row: nextRow, col: nextCol };

            if (!shift) {
                setSelection({ start: nextCell, end: nextCell });
            } else {
                setSelection(prevSel => ({ ...prevSel, end: nextCell }));
            }
            return nextCell;
        });
    }, [numRows, numCols]);

    const jumpCell = useCallback((direction: 'up' | 'down' | 'left' | 'right', shift: boolean) => {
        setActiveCell(prev => {
            let nextRow = prev.row;
            let nextCol = prev.col;
            if (direction === 'up') nextRow = 0;
            if (direction === 'down') nextRow = Math.max(0, numRows - 1);
            if (direction === 'left') nextCol = 0;
            if (direction === 'right') nextCol = Math.max(0, numCols - 1);

            const nextCell = { row: nextRow, col: nextCol };
            if (!shift) {
                setSelection({ start: nextCell, end: nextCell });
            } else {
                setSelection(prevSel => ({ ...prevSel, end: nextCell }));
            }
            return nextCell;
        });
    }, [numRows, numCols]);

    const handleKeyDown = useCallback((e: globalThis.KeyboardEvent) => {
        if (isEditing) {
            if (e.key === 'Escape') {
                setIsEditing(false);
            } else if (e.key === 'Enter' && !e.shiftKey) {
                setIsEditing(false);
                moveCell(1, 0, false);
                e.preventDefault();
            }
            return;
        }

        let handled = true;
        const shift = e.shiftKey;
        const ctrl = e.ctrlKey || e.metaKey;

        switch (e.key) {
            case 'ArrowUp':
                if (ctrl) jumpCell('up', shift);
                else moveCell(-1, 0, shift);
                break;
            case 'ArrowDown':
                if (ctrl) jumpCell('down', shift);
                else moveCell(1, 0, shift);
                break;
            case 'ArrowLeft':
                if (ctrl) jumpCell('left', shift);
                else moveCell(0, -1, shift);
                break;
            case 'ArrowRight':
                if (ctrl) jumpCell('right', shift);
                else moveCell(0, 1, shift);
                break;
            case 'Tab':
                moveCell(0, shift ? -1 : 1, false);
                break;
            case 'Enter':
                if (e.shiftKey) moveCell(-1, 0, false);
                else moveCell(1, 0, false);
                break;
            case 'F2':
                setIsEditing(true);
                break;
            case 'Delete':
            case 'Backspace':
                console.log("Cleared values.");
                break;
            case 'd':
            case 'D':
                if (ctrl) {
                    console.log("Fill Down Triggered!");
                    e.preventDefault();
                } else handled = false;
                break;
            case 'r':
            case 'R':
                if (ctrl) {
                    console.log("Fill Right Triggered!");
                    e.preventDefault();
                } else handled = false;
                break;
            case ' ':
                if (shift) {
                    setSelection(prev => ({ start: { row: activeCell.row, col: 0 }, end: { row: activeCell.row, col: numCols - 1 } }));
                } else if (ctrl) {
                    setSelection(prev => ({ start: { row: 0, col: activeCell.col }, end: { row: numRows - 1, col: activeCell.col } }));
                } else {
                    handled = false;
                }
                break;
            default:
                if (e.key.length === 1 && !ctrl && !e.altKey) {
                    setIsEditing(true);
                    setEditValue(e.key);
                    handled = false;
                } else {
                    handled = false;
                }
        }

        if (handled) {
            e.preventDefault();
        }
    }, [activeCell, isEditing, moveCell, jumpCell, numRows, numCols, setSelection]);

    useEffect(() => {
        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [handleKeyDown]);

    return {
        activeCell, setActiveCell,
        selection, setSelection,
        isEditing, setIsEditing,
        editValue, setEditValue
    };
};
