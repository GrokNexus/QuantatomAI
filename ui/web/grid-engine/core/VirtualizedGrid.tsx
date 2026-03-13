"use client";

import React from 'react';

interface VirtualizedGridProps {
    rows: string[];
    columns: string[];
    cells: Record<string, number>;
}

export function VirtualizedGrid({ rows, columns, cells }: VirtualizedGridProps) {
    return (
        <div className="grid-renderer overflow-auto border border-slate-800 rounded-lg">
            <table className="w-full border-collapse min-w-[480px]">
                <thead>
                    <tr>
                        <th className="text-left px-3 py-2 bg-slate-900 text-slate-300 sticky left-0">Row</th>
                        {columns.map((col, i) => (
                            <th key={i} className="text-right px-3 py-2 bg-slate-900 text-slate-300">
                                {col}
                            </th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {rows.map((rowLabel, rIdx) => {
                        const rowBg = rIdx % 2 === 0 ? 'bg-slate-950' : 'bg-slate-900/80';
                        return (
                            <tr key={rIdx} className={rowBg}>
                                <td className={`px-3 py-2 text-slate-100 font-semibold sticky left-0 ${rowBg}`}>
                                    {rowLabel}
                                </td>
                                {columns.map((_, cIdx) => {
                                    const key = `${rIdx}-${cIdx}`;
                                    const value = cells[key];
                                    return (
                                        <td key={cIdx} className="px-3 py-2 text-right text-slate-300">
                                            {value === undefined ? '' : value.toLocaleString(undefined, { maximumFractionDigits: 2 })}
                                        </td>
                                    );
                                })}
                            </tr>
                        );
                    })}
                </tbody>
            </table>
        </div>
    );
}
