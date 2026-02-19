import React from 'react';

export function VirtualizedGrid({ rows, columns, cells }) {
    return (
        <div className="grid-renderer">
            <table>
                <thead>
                    <tr>
                        {columns.map((col, i) => <th key={i}>{col.name}</th>)}
                    </tr>
                </thead>
                <tbody>
                    {rows.map((row, i) => (
                        <tr key={i}>
                            {columns.map((col, j) => (
                                <td key={j}>{cells[i * columns.length + j]?.value || ''}</td>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}
