                                        // Statistical functions
                                        if (/^VAR\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const arr = JSON.parse(args[0]);
                                                const avg = arr.reduce((a, b) => a + Number(b), 0) / arr.length;
                                                return arr.reduce((a, b) => a + Math.pow(Number(b) - avg, 2), 0) / arr.length;
                                            } catch { return 'ERR'; }
                                        }
                                        if (/^STDEV\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const arr = JSON.parse(args[0]);
                                                const avg = arr.reduce((a, b) => a + Number(b), 0) / arr.length;
                                                return Math.sqrt(arr.reduce((a, b) => a + Math.pow(Number(b) - avg, 2), 0) / arr.length);
                                            } catch { return 'ERR'; }
                                        }
                                        if (/^MEDIAN\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const arr = JSON.parse(args[0]).sort((a,b)=>a-b);
                                                const mid = Math.floor(arr.length/2);
                                                return arr.length%2!==0 ? arr[mid] : (arr[mid-1]+arr[mid])/2;
                                            } catch { return 'ERR'; }
                                        }
                                        if (/^MODE\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const arr = JSON.parse(args[0]);
                                                const freq = {};
                                                arr.forEach(v => freq[v] = (freq[v]||0)+1);
                                                return Object.keys(freq).reduce((a,b)=>freq[a]>freq[b]?a:b);
                                            } catch { return 'ERR'; }
                                        }
                                        // Financial functions
                                        if (/^IRR\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const values = JSON.parse(args[0]);
                                                let irr = 0.1;
                                                for (let iter = 0; iter < 100; iter++) {
                                                    let npv = 0;
                                                    for (let i = 0; i < values.length; i++) {
                                                        npv += values[i] / Math.pow(1 + irr, i);
                                                    }
                                                    if (Math.abs(npv) < 1e-6) break;
                                                    irr -= npv / 10000;
                                                }
                                                return (irr * 100).toFixed(2) + '%';
                                            } catch { return 'ERR'; }
                                        }
                                        if (/^NPV\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const rate = Number(args[0]);
                                                const values = JSON.parse(args[1]);
                                                return values.reduce((acc, v, i) => acc + v / Math.pow(1 + rate, i + 1), 0);
                                            } catch { return 'ERR'; }
                                        }
                                        // Advanced planning functions (mocked for demo)
                                        if (/^GOALSEEK\(/i.test(expr)) {
                                            // GOALSEEK(target, formula)
                                            const args = parseArgs(expr);
                                            return 'GoalSeek: ' + args.join(', ');
                                        }
                                        if (/^SOLVER\(/i.test(expr)) {
                                            // SOLVER(objective, constraints)
                                            const args = parseArgs(expr);
                                            return 'Solver: ' + args.join(', ');
                                        }
                                        if (/^SPREAD\(/i.test(expr)) {
                                            // SPREAD(value, range)
                                            const args = parseArgs(expr);
                                            return 'Spread: ' + args.join(', ');
                                        }
                                        if (/^ALLOCATE\(/i.test(expr)) {
                                            // ALLOCATE(total, drivers)
                                            const args = parseArgs(expr);
                                            return 'Allocate: ' + args.join(', ');
                                        }
                                        if (/^DRIVER\(/i.test(expr)) {
                                            // DRIVER(base, factor)
                                            const args = parseArgs(expr);
                                            return 'Driver: ' + args.join(', ');
                                        }
                                        if (/^TOPDOWN\(/i.test(expr)) {
                                            // TOPDOWN(total, levels)
                                            const args = parseArgs(expr);
                                            return 'TopDown: ' + args.join(', ');
                                        }
                                        if (/^RULE\(/i.test(expr)) {
                                            // RULE(expression)
                                            const args = parseArgs(expr);
                                            return 'Rule: ' + args.join(', ');
                                        }
                                        // Aggregation and ranking
                                        if (/^AGGREGATE\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            return 'Aggregate: ' + args.join(', ');
                                        }
                                        if (/^PERCENTILE\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const arr = JSON.parse(args[0]).sort((a,b)=>a-b);
                                                const p = Number(args[1]);
                                                const idx = Math.ceil(p * arr.length) - 1;
                                                return arr[idx];
                                            } catch { return 'ERR'; }
                                        }
                                        if (/^RANK\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const arr = JSON.parse(args[0]);
                                                const val = Number(args[1]);
                                                return arr.sort((a,b)=>b-a).indexOf(val) + 1;
                                            } catch { return 'ERR'; }
                                        }
                                        if (/^LARGE\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const arr = JSON.parse(args[0]).sort((a,b)=>b-a);
                                                const n = Number(args[1]);
                                                return arr[n-1];
                                            } catch { return 'ERR'; }
                                        }
                                        if (/^SMALL\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const arr = JSON.parse(args[0]).sort((a,b)=>a-b);
                                                const n = Number(args[1]);
                                                return arr[n-1];
                                            } catch { return 'ERR'; }
                                        }
                                        if (/^SUMPRODUCT\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const arr1 = JSON.parse(args[0]);
                                                const arr2 = JSON.parse(args[1]);
                                                return arr1.reduce((acc, v, i) => acc + v * arr2[i], 0);
                                            } catch { return 'ERR'; }
                                        }
                                        if (/^PRODUCT\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            try {
                                                const arr = JSON.parse(args[0]);
                                                return arr.reduce((acc, v) => acc * Number(v), 1);
                                            } catch { return 'ERR'; }
                                        }
                                        if (/^SUBTOTAL\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            return 'Subtotal: ' + args.join(', ');
                                        }
                                        if (/^GROUPBY\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            return 'GroupBy: ' + args.join(', ');
                                        }
                                        // Forecasting and regression
                                        if (/^FORECAST\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            return 'Forecast: ' + args.join(', ');
                                        }
                                        if (/^TREND\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            return 'Trend: ' + args.join(', ');
                                        }
                                        if (/^SLOPE\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            return 'Slope: ' + args.join(', ');
                                        }
                                        if (/^INTERCEPT\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            return 'Intercept: ' + args.join(', ');
                                        }
                                        if (/^CORREL\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            return 'Correl: ' + args.join(', ');
                                        }
                                        if (/^COVAR\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            return 'Covar: ' + args.join(', ');
                                        }
                                        if (/^REGRESSION\(/i.test(expr)) {
                                            const args = parseArgs(expr);
                                            return 'Regression: ' + args.join(', ');
                                        }
                                // Text functions
                                if (/^CONCAT\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return args.join('');
                                }
                                if (/^LEFT\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return args[0].substring(0, Number(args[1]));
                                }
                                if (/^RIGHT\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return args[0].slice(-Number(args[1]));
                                }
                                if (/^MID\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return args[0].substring(Number(args[1]) - 1, Number(args[1]) - 1 + Number(args[2]));
                                }
                                if (/^LEN\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return args[0].length;
                                }
                                if (/^FIND\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return args[0].indexOf(args[1]) + 1;
                                }
                                if (/^SEARCH\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return args[0].toLowerCase().indexOf(args[1].toLowerCase()) + 1;
                                }
                                // Date functions
                                if (/^DATE\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return new Date(Number(args[0]), Number(args[1]) - 1, Number(args[2])).toISOString().slice(0, 10);
                                }
                                if (/^TODAY\(/i.test(expr)) {
                                    return new Date().toISOString().slice(0, 10);
                                }
                                if (/^NOW\(/i.test(expr)) {
                                    return new Date().toISOString();
                                }
                                if (/^YEAR\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return new Date(args[0]).getFullYear();
                                }
                                if (/^MONTH\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return new Date(args[0]).getMonth() + 1;
                                }
                                if (/^DAY\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return new Date(args[0]).getDate();
                                }
                                // Lookup functions
                                if (/^HLOOKUP\(/i.test(expr)) {
                                    // HLOOKUP(lookup_value, table_array, row_index)
                                    const args = parseArgs(expr);
                                    const [lookupValue, tableKey, rowIndex] = args;
                                    if (typeof window.gridData !== 'undefined') {
                                        const col = window.gridData.find(r => r[tableKey] == lookupValue);
                                        if (col) {
                                            const keys = Object.keys(col);
                                            return col[keys[Number(rowIndex) - 1]];
                                        }
                                    }
                                    return 'ERR';
                                }
                                if (/^LOOKUP\(/i.test(expr)) {
                                    // LOOKUP(lookup_value, array)
                                    const args = parseArgs(expr);
                                    try {
                                        const lookup = args[0];
                                        const arr = JSON.parse(args[1]);
                                        return arr.find(v => v == lookup) || 'ERR';
                                    } catch { return 'ERR'; }
                                }
                                // Array functions
                                if (/^TRANSPOSE\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    try {
                                        const arr = JSON.parse(args[0]);
                                        return arr.map((_, i) => arr.map(row => row[i]));
                                    } catch { return 'ERR'; }
                                }
                                if (/^UNIQUE\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    try {
                                        const arr = JSON.parse(args[0]);
                                        return Array.from(new Set(arr));
                                    } catch { return 'ERR'; }
                                }
                                if (/^SORT\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    try {
                                        const arr = JSON.parse(args[0]);
                                        return arr.sort();
                                    } catch { return 'ERR'; }
                                }
                                if (/^FILTER\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    try {
                                        const arr = JSON.parse(args[0]);
                                        const crit = args[1];
                                        return arr.filter(v => eval(crit.replace(/x/g, v)));
                                    } catch { return 'ERR'; }
                                }
                                // IS functions
                                if (/^ISBLANK\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return args[0] === '' || args[0] === null;
                                }
                                if (/^ISNUMBER\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return !isNaN(Number(args[0]));
                                }
                                if (/^ISTEXT\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return typeof args[0] === 'string';
                                }
                                if (/^ISERROR\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return args[0] === 'ERR';
                                }
                                // Math functions
                                if (/^ROUND\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return Math.round(Number(args[0]));
                                }
                                if (/^ROUNDUP\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return Math.ceil(Number(args[0]));
                                }
                                if (/^ROUNDDOWN\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return Math.floor(Number(args[0]));
                                }
                                if (/^FLOOR\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return Math.floor(Number(args[0]));
                                }
                                if (/^CEILING\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return Math.ceil(Number(args[0]));
                                }
                                if (/^ABS\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return Math.abs(Number(args[0]));
                                }
                                if (/^SQRT\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return Math.sqrt(Number(args[0]));
                                }
                                if (/^POWER\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return Math.pow(Number(args[0]), Number(args[1]));
                                }
                                if (/^EXP\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return Math.exp(Number(args[0]));
                                }
                                if (/^LOG\(/i.test(expr)) {
                                    const args = parseArgs(expr);
                                    return Math.log(Number(args[0]));
                                }
                        // INDEX
                        if (/^INDEX\(/i.test(expr)) {
                            const args = parseArgs(expr);
                            // INDEX(array, row_num)
                            try {
                                const arr = JSON.parse(args[0]);
                                const idx = Number(args[1]) - 1;
                                return arr[idx];
                            } catch { return 'ERR'; }
                        }
                        // MATCH
                        if (/^MATCH\(/i.test(expr)) {
                            const args = parseArgs(expr);
                            // MATCH(lookup_value, array)
                            try {
                                const lookup = args[0];
                                const arr = JSON.parse(args[1]);
                                return arr.indexOf(lookup) + 1;
                            } catch { return 'ERR'; }
                        }
                        // COUNT
                        if (/^COUNT\(/i.test(expr)) {
                            const args = parseArgs(expr);
                            try {
                                const arr = JSON.parse(args[0]);
                                return arr.filter(v => v !== null && v !== '').length;
                            } catch { return 'ERR'; }
                        }
                        // COUNTA
                        if (/^COUNTA\(/i.test(expr)) {
                            const args = parseArgs(expr);
                            try {
                                const arr = JSON.parse(args[0]);
                                return arr.length;
                            } catch { return 'ERR'; }
                        }
                        // COUNTIF
                        if (/^COUNTIF\(/i.test(expr)) {
                            const args = parseArgs(expr);
                            try {
                                const arr = JSON.parse(args[0]);
                                const crit = args[1];
                                return arr.filter(v => v == crit).length;
                            } catch { return 'ERR'; }
                        }
                        // SUM
                        if (/^SUM\(/i.test(expr)) {
                            const args = parseArgs(expr);
                            try {
                                const arr = JSON.parse(args[0]);
                                return arr.reduce((acc, v) => acc + Number(v), 0);
                            } catch { return 'ERR'; }
                        }
                        // AVERAGE
                        if (/^AVERAGE\(/i.test(expr)) {
                            const args = parseArgs(expr);
                            try {
                                const arr = JSON.parse(args[0]);
                                return arr.length ? arr.reduce((acc, v) => acc + Number(v), 0) / arr.length : 0;
                            } catch { return 'ERR'; }
                        }
                        // MIN
                        if (/^MIN\(/i.test(expr)) {
                            const args = parseArgs(expr);
                            try {
                                const arr = JSON.parse(args[0]);
                                return Math.min(...arr);
                            } catch { return 'ERR'; }
                        }
                        // MAX
                        if (/^MAX\(/i.test(expr)) {
                            const args = parseArgs(expr);
                            try {
                                const arr = JSON.parse(args[0]);
                                return Math.max(...arr);
                            } catch { return 'ERR'; }
                        }
                // Helper: parse arguments for functions
                function parseArgs(fnCall) {
                    const args = fnCall.match(/\((.*)\)/)[1].split(',').map(a => a.trim());
                    return args;
                }
                    // Logical functions
                    if (/^IF\(/i.test(expr)) {
                        const args = parseArgs(expr);
                        // IF(condition, true_val, false_val)
                        // Only supports simple conditions for now
                        const cond = args[0];
                        let condResult = false;
                        try { condResult = eval(cond); } catch {}
                        return condResult ? args[1] : args[2];
                    }
                    if (/^AND\(/i.test(expr)) {
                        const args = parseArgs(expr);
                        return args.every(arg => { try { return !!eval(arg); } catch { return false; } });
                    }
                    if (/^OR\(/i.test(expr)) {
                        const args = parseArgs(expr);
                        return args.some(arg => { try { return !!eval(arg); } catch { return false; } });
                    }
                    if (/^NOT\(/i.test(expr)) {
                        const args = parseArgs(expr);
                        return !(args.length && eval(args[0]));
                    }
            // Backend sync: Collaboration features
            async function syncCellActivity(rowId, colKey, userId) {
                try {
                    const response = await fetch('/api/grid/cell-activity', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ rowId, colKey, userId })
                    });
                    if (!response.ok) throw new Error('Failed to sync cell activity');
                    const result = await response.json();
                    setCellActivity(activity => ({ ...activity, [`${rowId}-${colKey}`]: result.userId }));
                } catch (err) {
                    console.error('Cell activity sync failed:', err);
                }
            }

            async function syncOnlineUsers() {
                try {
                    const response = await fetch('/api/grid/online-users', {
                        method: 'GET',
                        headers: { 'Content-Type': 'application/json' }
                    });
                    if (!response.ok) throw new Error('Failed to fetch online users');
                    const result = await response.json();
                    setOnlineUsers(result.users);
                } catch (err) {
                    console.error('Online users sync failed:', err);
                }
            }
            function handleCellClick(rowId, colKey, userId) {
                // Optimistic UI update
                setCellActivity(activity => ({ ...activity, [`${rowId}-${colKey}`]: userId }));
                // Backend sync
                syncCellActivity(rowId, colKey, userId);
            }
            useEffect(() => {
                if (collabPanelOpen) {
                    syncOnlineUsers();
                }
            }, [collabPanelOpen]);
        // Handler for adding comment
        function handleAddComment() {
            if (!commentCell || !newComment) return;
            const { rowId, colKey } = commentCell;
            // Optimistic UI update
            setCellComments(comments => {
                const key = `${rowId}-${colKey}`;
                const newComments = { ...comments };
                newComments[key] = [...(newComments[key] || []), { user: 'Me', text: newComment, time: new Date().toLocaleString() }];
                return newComments;
            });
            // Backend sync
            syncAddComment(rowId, colKey, { user: 'Me', text: newComment, time: new Date().toLocaleString() });
            setNewComment('');
        }
    // Backend sync: Comments and Audit
    async function syncAddComment(rowId, colKey, comment) {
        try {
            const response = await fetch('/api/grid/add-comment', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ rowId, colKey, comment })
            });
            if (!response.ok) throw new Error('Failed to sync comment');
            const result = await response.json();
            setCellComments(comments => {
                const key = `${rowId}-${colKey}`;
                const newComments = { ...comments };
                newComments[key] = result.comments;
                return newComments;
            });
        } catch (err) {
            console.error('Comment sync failed:', err);
        }
    }

    async function syncAuditTrail(rowId, colKey, user, oldValue, newValue) {
        try {
            const response = await fetch('/api/grid/audit', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ rowId, colKey, user, oldValue, newValue })
            });
            if (!response.ok) throw new Error('Failed to sync audit');
            const result = await response.json();
            setAuditTrail(trail => [...result.auditTrail]);
        } catch (err) {
            console.error('Audit sync failed:', err);
        }
    }
// Theme studio: custom themes
        const [theme, setTheme] = useState('default');
        const themeOptions = ['default', 'dark', 'glass', 'emerald', 'rose', 'amber'];
        function handleThemeChange(e) {
            setTheme(e.target.value);
        }
    // Height logic: multi-line, auto-height
    function getRowHeight(row) {
        // Example: calculate height based on max cell content length
        const maxLen = columnOrder.reduce((max, colKey) => {
            const val = row[colKey];
            return Math.max(max, val && val.toString().length ? val.toString().length : 0);
        }, 0);
        return Math.max(40, Math.min(120, maxLen * 2)); // min 40px, max 120px
    }
// ...existing imports...
import React, { useState, useEffect, useRef } from 'react';
import {
    TrendingUp, TrendingDown, ArrowUpRight, Filter, Download, MoreHorizontal, Plus, Eye, EyeOff, Maximize2, Minimize2, X, RotateCcw
} from 'lucide-react';
import { useTheme } from '@/context/ThemeContext';
import {
    AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer
} from 'recharts';

// Helper functions
function formatCurrency(value) {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(value);
}
function formatPercent(value) {
    return `${value > 0 ? '+' : ''}${value.toFixed(1)}%`;
}
function getHealthColor(health, theme) {
    if (health >= 90) return theme === 'dark' ? 'bg-emerald-600' : 'bg-emerald-500';
    if (health >= 70) return theme === 'dark' ? 'bg-amber-600' : 'bg-amber-500';
    return theme === 'dark' ? 'bg-rose-600' : 'bg-rose-500';
}
function getVarianceClass(variance, displayProfile) {
    if (displayProfile === 'default') {
        return variance > 10 ? 'bg-emerald-100 text-emerald-700' : variance < -10 ? 'bg-rose-100 text-rose-700' : '';
    } else if (displayProfile === 'bold') {
        return variance > 10 ? 'font-bold text-emerald-700' : variance < -10 ? 'font-bold text-rose-700' : '';
    }
    return '';
}
function getVarianceHeatmapStyle(variance) {
    const intensity = Math.min(Math.abs(variance) / 20, 1);
    if (variance > 0) {
        return { background: `rgba(16, 185, 129, ${intensity * 0.3})` };
    } else if (variance < 0) {
        return { background: `rgba(239, 68, 68, ${intensity * 0.3})` };
    }
    return {};
}
function Sparkline(props) {
    const data = props.data || [];
    if (!data.length) return null;
    const width = 60, height = 24;
    const min = Math.min(...data), max = Math.max(...data);
    const points = data.map((v, i) => {
        const x = (i / (data.length - 1)) * width;
        const y = height - ((v - min) / (max - min || 1)) * height;
        return `${x},${y}`;
    }).join(' ');
    return (
        <svg width={width} height={height} viewBox={`0 0 ${width} ${height}`}>
            <polyline points={points} fill="none" stroke="#10b981" strokeWidth="2" />
        </svg>
    );
}

// Mock data
const MOCK_GRID_DATA = [
    { id: 'QA-001', entity: 'North America Ops', status: 'On Track', budget: 1200000, actual: 1150000, variance: 4.2, health: 92 },
    { id: 'QA-002', entity: 'EMEA Logistics', status: 'At Risk', budget: 850000, actual: 890000, variance: -4.7, health: 68 }
];
const MOCK_HIERARCHY_DATA = [
    { id: 'QA-001', entity: 'North America Ops', status: 'On Track', budget: 1200000, actual: 1150000, variance: 4.2, health: 92, children: [
        { id: 'QA-001-1', entity: 'US', status: 'On Track', budget: 800000, actual: 780000, variance: 2.6, health: 95 },
        { id: 'QA-001-2', entity: 'Canada', status: 'At Risk', budget: 400000, actual: 370000, variance: 7.5, health: 88 }
    ] },
    { id: 'QA-002', entity: 'EMEA Logistics', status: 'At Risk', budget: 850000, actual: 890000, variance: -4.7, health: 68 }
];

function evaluateFormula(formula, row, namedRanges = {}) {
    if (!formula.startsWith('=')) return formula;
    try {
        let expr = formula.slice(1);
        // Replace named ranges first
        Object.keys(namedRanges).forEach(name => {
            expr = expr.replaceAll(name, namedRanges[name]);
        });
        // Replace column names with values
        Object.keys(row).forEach(key => {
            expr = expr.replaceAll(key, row[key]);
        });
        // Modular Excel-like formula engine
        // Supported: SUMIFS, VLOOKUP, XIRR, basic math
        if (/^SUMIFS\(/i.test(expr)) {
            // SUMIFS(range, criteria_range, criteria)
            // Example: SUMIFS(budget, status, "On Track")
            const args = expr.match(/SUMIFS\((.*)\)/i)[1].split(',').map(a => a.trim());
            const [sumRangeKey, criteriaRangeKey, criteria] = args;
            // Assume gridData is available in closure
            if (typeof window.gridData !== 'undefined') {
                return window.gridData.filter(r => r[criteriaRangeKey] == criteria).reduce((acc, r) => acc + Number(r[sumRangeKey]), 0);
            }
            return 'ERR';
        }
        if (/^VLOOKUP\(/i.test(expr)) {
            // VLOOKUP(lookup_value, table_array, col_index)
            const args = expr.match(/VLOOKUP\((.*)\)/i)[1].split(',').map(a => a.trim());
            const [lookupValue, tableKey, colIndex] = args;
            if (typeof window.gridData !== 'undefined') {
                const row = window.gridData.find(r => r[tableKey] == lookupValue);
                if (row) {
                    const keys = Object.keys(row);
                    return row[keys[Number(colIndex) - 1]];
                }
            }
            return 'ERR';
        }
        if (/^XIRR\(/i.test(expr)) {
            // XIRR(values, dates)
            // Example: XIRR([1000,-500,-500], ["2024-01-01","2024-06-01","2025-01-01"])
            const args = expr.match(/XIRR\((.*)\)/i)[1].split(',').map(a => a.trim());
            try {
                const values = JSON.parse(args[0]);
                const dates = JSON.parse(args[1]);
                // Simple IRR calculation (not full XIRR)
                const n = values.length;
                let irr = 0.1;
                for (let iter = 0; iter < 100; iter++) {
                    let npv = 0;
                    for (let i = 0; i < n; i++) {
                        const t = (new Date(dates[i]) - new Date(dates[0])) / (365 * 24 * 3600 * 1000);
                        npv += values[i] / Math.pow(1 + irr, t);
                    }
                    if (Math.abs(npv) < 1e-6) break;
                    irr -= npv / 10000;
                }
                return (irr * 100).toFixed(2) + '%';
            } catch {
                return 'ERR';
            }
        }
        // Fallback: basic math
        // eslint-disable-next-line no-eval
        return eval(expr);
    } catch (err) {
        return 'ERR: ' + err.message;
    }
}

function getPrecedentKeys(formula, columnOrder, namedRanges) {
    if (!formula.startsWith('=')) return [];
    const precedents = [];
    // Find column keys referenced in formula
    columnOrder.forEach(colKey => {
        if (formula.includes(colKey)) precedents.push(colKey);
    });
    // Find named ranges referenced in formula
    Object.keys(namedRanges).forEach(name => {
        if (formula.includes(name)) precedents.push(name);
    });
    return precedents;
}

function linearRegression(y) {
    // Simple linear regression for y = a + bx, x = [0,1,2,...]
    const n = y.length;
    if (n < 2) return { a: y[0] || 0, b: 0 };
    const x = Array.from({ length: n }, (_, i) => i);
    const sumX = x.reduce((a, b) => a + b, 0);
    const sumY = y.reduce((a, b) => a + b, 0);
    const sumXY = x.reduce((a, b, i) => a + b * y[i], 0);
    const sumXX = x.reduce((a, b) => a + b * b, 0);
    const b = (n * sumXY - sumX * sumY) / (n * sumXX - sumX * sumX || 1);
    const a = (sumY - b * sumX) / n;
    return { a, b };
}

function validateValue(value, rule) {
    if (!rule) return '';
    if (rule.required && (value === undefined || value === null || value === '')) return 'Required';
    if (rule.type === 'number' && isNaN(Number(value))) return 'Must be a number';
    if (rule.min !== undefined && Number(value) < rule.min) return `Min ${rule.min}`;
    if (rule.max !== undefined && Number(value) > rule.max) return `Max ${rule.max}`;
    return '';
}

function isDateCol(colKey, gridData) {
    // Heuristic: if most values in col are valid dates, treat as date col
    let valid = 0, total = 0;
    for (const row of gridData) {
        if (row[colKey] !== undefined && row[colKey] !== null && row[colKey] !== '') {
            total++;
            if (!isNaN(Date.parse(row[colKey]))) valid++;
        }
    }
    return total > 0 && valid / total > 0.7;
}

function isDateInRange(dateStr, from, to) {
    const d = new Date(dateStr);
    if (isNaN(d)) return false;
    if (from && d < new Date(from)) return false;
    if (to && d > new Date(to)) return false;
    return true;
}

function getDatePresetRange(preset) {
    const now = new Date();
    if (preset === 'today') {
        const from = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        const to = new Date(from); to.setDate(from.getDate() + 1); to.setMilliseconds(-1);
        return { from: from.toISOString().slice(0, 10), to: to.toISOString().slice(0, 10) };
    }
    if (preset === 'this_week') {
        const from = new Date(now); from.setDate(now.getDate() - now.getDay());
        const to = new Date(from); to.setDate(from.getDate() + 7); to.setMilliseconds(-1);
        return { from: from.toISOString().slice(0, 10), to: to.toISOString().slice(0, 10) };
    }
    if (preset === 'this_month') {
        const from = new Date(now.getFullYear(), now.getMonth(), 1);
        const to = new Date(now.getFullYear(), now.getMonth() + 1, 1); to.setMilliseconds(-1);
        return { from: from.toISOString().slice(0, 10), to: to.toISOString().slice(0, 10) };
    }
    return { from: '', to: '' };
}

function computePivot(data, rowField, colField, valueField, agg) {
    const rowKeys = Array.from(new Set(data.map(r => r[rowField])));
    const colKeys = Array.from(new Set(data.map(r => r[colField])));
    const table = {};
    rowKeys.forEach(rk => { table[rk] = {}; colKeys.forEach(ck => { table[rk][ck] = []; }); });
    data.forEach(row => {
        table[row[rowField]][row[colField]].push(row[valueField]);
    });
    const aggFn = (arr) => {
        if (agg === 'sum') return arr.reduce((a, b) => a + Number(b), 0);
        if (agg === 'avg') return arr.length ? arr.reduce((a, b) => a + Number(b), 0) / arr.length : 0;
        if (agg === 'count') return arr.length;
        return '';
    };
    const result = rowKeys.map(rk =>
        colKeys.map(ck => aggFn(table[rk][ck]))
    );
    return { rowKeys, colKeys, result };
}

export function QuantAtomGrid() {
    // --- AG Grid Parity & Beyond ---
    // Advanced Clipboard Operations
    function handleCopyRange(range, includeHeaders = false) {
        // TODO: Excel-compatible, tab-delimited, headers option
    }
    function handleCutRange(range) {
        // TODO: Cut logic
    }
    function handlePasteRange(range, data) {
        // TODO: Paste logic
    }
    function processClipboardData(data, type) {
        // TODO: Custom processing for clipboard data (cell/header/group)
    }
    const clipboardAPI = {
        copySelectedRows: () => {},
        copySelectedRange: () => {},
        cutSelectedRange: () => {},
        pasteToRange: () => {},
    };

    // Accessibility Features
    const accessibilityConfig = {
        aria: {
            role: 'grid',
            ariaRowCount: hierarchyData.length,
            ariaColCount: columnOrder.length,
            ariaMultiselectable: true,
            ariaLabel: 'QuantAtomGrid',
        },
        keyboardNavigation: true,
        screenReaderSupport: true,
        highContrastTheme: false,
        ensureDomOrder: true,
        suppressColumnVirtualisation: false,
        suppressRowVirtualisation: false,
    };

    // Column Moving
    function moveColumn(fromIndex, toIndex) {
        const newOrder = [...columnOrder];
        const [col] = newOrder.splice(fromIndex, 1);
        newOrder.splice(toIndex, 0, col);
        setColumnOrder(newOrder);
    }
    function lockColumnPosition(colId, position) {
        if (position === 'left') handlePinLeft(colId);
        else if (position === 'right') handlePinRight(colId);
    }
    const suppressColumnMoveAnimation = false;
    const suppressDragLeaveHidesColumns = true;
    const suppressMoveWhenColumnDragging = false;

    // Row Dragging
    const enableRowDragging = true;
    const rowDragManaged = true;
    function rowDragCallback(params) {
        // Custom row drag logic
    }

    // Pivoting
    const pivotConfig = {
        pivotMode: pivotMode,
        pivotColumns: [pivotCol],
        enablePivot: true,
        pivotPanelShow: 'onlyWhenPivoting',
    };
    function setPivotColumns(cols) { /* TODO */ }
    function getPivotResultColumns() { /* TODO */ }

    // Aggregation
    const aggregationConfig = {
        aggFuncs: {
            sum: (values) => values.reduce((a, b) => a + b, 0),
            min: (values) => Math.min(...values),
            max: (values) => Math.max(...values),
            avg: (values) => values.length ? values.reduce((a, b) => a + b, 0) / values.length : 0,
            count: (values) => values.length,
            first: (values) => values[0],
            last: (values) => values[values.length - 1],
            custom: (values) => {/* Custom aggregation */},
        },
        groupTotalRow: 'bottom',
        grandTotalRow: 'bottom',
    };

    // Master-Detail
    const masterDetailConfig = {
        masterDetail: true,
        isRowMaster: (row) => !!row.detail,
        detailCellRenderer: null,
        detailCellRendererParams: {},
        detailRowHeight: 200,
        keepDetailRows: true,
        keepDetailRowsCount: 10,
    };

    // Filtering
    const filterConfig = {
        filters: {
            text: (value, filter) => value.includes(filter),
            number: (value, filter) => value === filter,
            date: (value, filter) => value === filter,
            set: (value, filterSet) => filterSet.has(value),
            multi: (value, filters) => filters.some(f => f(value)),
        },
        quickFilter: '',
        externalFilter: null,
    };

    // Sidebar/Tool Panels
    const sideBarConfig = {
        toolPanels: [
            { id: 'columns', label: 'Columns', icon: 'columns', component: null },
            { id: 'filters', label: 'Filters', icon: 'filter', component: null },
        ],
        position: 'left',
        defaultToolPanel: 'columns',
        hiddenByDefault: false,
        parent: null,
    };
    function setSideBarVisible(visible) { setShowToolPanel(visible); }
    function openToolPanel(id) { setShowToolPanel(true); }
    function closeToolPanel() { setShowToolPanel(false); }

    // Extensibility Hooks
    const extensibilityHooks = {
        onCellRender: null,
        onRowRender: null,
        onColumnRender: null,
        onSidebarRender: null,
        onClipboard: null,
        onAggregation: null,
        onPivot: null,
        onMasterDetail: null,
        onFilter: null,
    };

    // UI Controls
    function renderGridControls() {
        return (
            <div className="quantatom-grid-controls mb-4 flex flex-wrap gap-2">
                {/* Clipboard controls */}
                <button className="glass-button" onClick={() => handleCopyRange(/*range*/)}>Copy Range</button>
                <button className="glass-button" onClick={() => handleCutRange(/*range*/)}>Cut Range</button>
                <button className="glass-button" onClick={() => handlePasteRange(/*range, data*/)}>Paste Range</button>
                {/* Visual Reporting Enhancements */}
                <button className="glass-button" onClick={() => setShowHeatmap(v => !v)}>Toggle Heatmap</button>
                <button className="glass-button" onClick={() => setShowSankey(v => !v)}>Show Sankey</button>
                <button className="glass-button" onClick={() => setShowExportPanel(v => !v)}>Export/Import</button>
                {/* Offline/Conflict Handling */}
                <button className="glass-button" onClick={() => setShowSnapshotPanel(v => !v)}>Snapshot/Delta Log</button>
                <button className="glass-button" onClick={() => setShowMergeUX(v => !v)}>Guided Merge UX</button>
                {/* Security/Governance UI */}
                <button className="glass-button" onClick={() => setShowRBACPanel(v => !v)}>RBAC/Privacy</button>
                <button className="glass-button" onClick={() => setShowAuditPanel(v => !v)}>Audit/ESG</button>
                {/* Accessibility/Keyboard UX */}
                <button className="glass-button" onClick={() => accessibilityConfig.highContrastTheme = !accessibilityConfig.highContrastTheme}>Toggle High Contrast</button>
                <button className="glass-button" onClick={() => setShowKeyboardShortcuts(v => !v)}>Keyboard Shortcuts</button>
                {/* Extensibility/Plugin Support */}
                <button className="glass-button" onClick={() => setShowPluginPanel(v => !v)}>Plugin APIs</button>
                {/* Column moving controls */}
                <button className="glass-button" onClick={() => moveColumn(0, 1)}>Move Column</button>
                {/* Row dragging controls */}
                <button className="glass-button" onClick={() => {/* toggle row dragging */}}>Toggle Row Dragging</button>
                {/* Pivot controls */}
                <button className="glass-button" onClick={() => setPivotMode(v => !v)}>Toggle Pivot Mode</button>
                {/* Aggregation controls */}
                <button className="glass-button" onClick={() => aggregationConfig.groupTotalRow = aggregationConfig.groupTotalRow === 'bottom' ? 'top' : 'bottom'}>Toggle Group Total Row</button>
                {/* Master-detail controls */}
                <button className="glass-button" onClick={() => masterDetailConfig.masterDetail = !masterDetailConfig.masterDetail}>Toggle Master Detail</button>
                {/* Filter controls */}
                <button className="glass-button" onClick={() => filterConfig.quickFilter = ''}>Clear Quick Filter</button>
                {/* Sidebar controls */}
                <button className="glass-button" onClick={() => setSideBarVisible(true)}>Show Sidebar</button>
                <button className="glass-button" onClick={() => setSideBarVisible(false)}>Hide Sidebar</button>
                <button className="glass-button" onClick={() => openToolPanel('columns')}>Open Columns Panel</button>
                <button className="glass-button" onClick={() => openToolPanel('filters')}>Open Filters Panel</button>
                <button className="glass-button" onClick={closeToolPanel}>Close Tool Panel</button>
            </div>
        );
        // --- Enterprise gap state hooks ---
        const [showHeatmap, setShowHeatmap] = useState(false);
        const [showSankey, setShowSankey] = useState(false);
        const [showExportPanel, setShowExportPanel] = useState(false);
        const [showSnapshotPanel, setShowSnapshotPanel] = useState(false);
        const [showMergeUX, setShowMergeUX] = useState(false);
        const [showRBACPanel, setShowRBACPanel] = useState(false);
        const [showAuditPanel, setShowAuditPanel] = useState(false);
        const [showKeyboardShortcuts, setShowKeyboardShortcuts] = useState(false);
        const [showPluginPanel, setShowPluginPanel] = useState(false);
    }
    const [draggedRowId, setDraggedRowId] = useState(null);
    const [hierarchyData, setHierarchyData] = useState(MOCK_HIERARCHY_DATA);
    const [expandedTreeRows, setExpandedTreeRows] = useState([]);
    const [showToolPanel, setShowToolPanel] = useState(false);
    const [columnOrder, setColumnOrder] = useState(['entity', 'status', 'budget', 'actual', 'variance', 'health']);
    const [hiddenColumns, setHiddenColumns] = useState([]);
    const [pinnedLeft, setPinnedLeft] = useState([]);
    const [pinnedRight, setPinnedRight] = useState([]);

    // Backend sync: REST API endpoint
    async function syncCellEdit(rowId, colKey, newValue) {
        try {
            const response = await fetch('/api/grid/cell-edit', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ rowId, colKey, newValue })
            });
            if (!response.ok) throw new Error('Failed to sync cell edit');
            const result = await response.json();
            // Optionally update state with backend-confirmed value
            setHierarchyData(data => {
                const newData = [...data];
                const idx = newData.findIndex(r => r.id === rowId);
                if (idx !== -1) {
                    newData[idx] = { ...newData[idx], [colKey]: result.value };
                }
                return newData;
            });
        } catch (err) {
            // Handle error, show notification, revert UI, etc.
            console.error('Cell edit sync failed:', err);
        }
    }
    const [columnWidths, setColumnWidths] = useState({});
    const [isMaximized, setIsMaximized] = useState(false);
    // Flashing cells on update
    const [flashingCells, setFlashingCells] = useState({});
    function triggerCellFlash(rowId, colKey) {
        setFlashingCells(prev => ({ ...prev, [`${rowId}-${colKey}`]: true }));
        setTimeout(() => {
            setFlashingCells(prev => ({ ...prev, [`${rowId}-${colKey}`]: false }));
        }, 600);
    }
    // Example: simulate cell update (for demo)
    useEffect(() => {
        const interval = setInterval(() => {
            const row = hierarchyData[Math.floor(Math.random() * hierarchyData.length)];
            if (!row) return;
            const colKey = columnOrder[Math.floor(Math.random() * columnOrder.length)];
            triggerCellFlash(row.id, colKey);
        }, 4000);
        return () => clearInterval(interval);
    }, [hierarchyData, columnOrder]);
    // Row Dragging
    function handleRowDragStart(rowId) { setDraggedRowId(rowId); }
    function handleRowDragOver(e) { e.preventDefault(); }
    function handleRowDrop(targetRowId) {
        if (!draggedRowId || draggedRowId === targetRowId) return;
        const newData = [...hierarchyData];
        const draggedIdx = newData.findIndex(r => r.id === draggedRowId);
        const targetIdx = newData.findIndex(r => r.id === targetRowId);
        if (draggedIdx < 0 || targetIdx < 0) return;
        const [draggedRow] = newData.splice(draggedIdx, 1);
        newData.splice(targetIdx, 0, draggedRow);
        setHierarchyData(newData);
        setDraggedRowId(null);
    }
    // Tree expand/collapse
    function handleExpandTreeRow(rowId) {
        setExpandedTreeRows(rows => rows.includes(rowId) ? rows.filter(r => r !== rowId) : [...rows, rowId]);
    }
    // Column/Tool Panel
    function handleToggleColumn(colKey) {
        setHiddenColumns(cols =>
            cols.includes(colKey) ? cols.filter(c => c !== colKey) : [...cols, colKey]
        );
    }
    function handleReorderColumn(colKey, direction) {
        const idx = columnOrder.indexOf(colKey);
        if (idx < 0) return;
        let newOrder = [...columnOrder];
        if (direction === 'left' && idx > 0) {
            [newOrder[idx - 1], newOrder[idx]] = [newOrder[idx], newOrder[idx - 1]];
        } else if (direction === 'right' && idx < newOrder.length - 1) {
            [newOrder[idx], newOrder[idx + 1]] = [newOrder[idx + 1], newOrder[idx]];
        }
        setColumnOrder(newOrder);
    }
    function handlePinLeft(colKey) {
        setPinnedLeft(cols => [...cols, colKey]);
        setPinnedRight(cols => cols.filter(c => c !== colKey));
    }
    function handlePinRight(colKey) {
        setPinnedRight(cols => [...cols, colKey]);
        setPinnedLeft(cols => cols.filter(c => c !== colKey));
    }
    function handleUnpin(colKey) {
        setPinnedLeft(cols => cols.filter(c => c !== colKey));
        setPinnedRight(cols => cols.filter(c => c !== colKey));
    }
    function autoSizeColumn(colKey) {
        // Example: auto-size based on mock data
        const maxLen = Math.max(...hierarchyData.map(row => (row[colKey] ? row[colKey].toString().length : 0)));
        setColumnWidths(w => ({ ...w, [colKey]: Math.max(60, maxLen * 10 + 32) }));
    }
    // Named range add handler
    function handleAddNamedRange() {
        if (!newRangeName || !newRangeValue) return;
        setNamedRanges(ranges => ({ ...ranges, [newRangeName]: newRangeValue }));
        setNewRangeName('');
        setNewRangeValue('');
    }
    // Render tree rows
    function renderTreeRows(rows, level = 0) {
        return rows.map(row => (
            <React.Fragment key={row.id}>
                <tr
                    className="hover:bg-white/5 transition-colors group"
                    draggable
                    onDragStart={() => handleRowDragStart(row.id)}
                    onDragOver={handleRowDragOver}
                    onDrop={() => handleRowDrop(row.id)}
                    style={{ height: getRowHeight(row) }}
                >
                    {columnOrder.map(colKey => {
                        const error = validationErrors[`${row.id}-${colKey}`];
                        const userId = cellActivity[`${row.id}-${colKey}`];
                        const user = onlineUsers.find(u => u.id === userId);
                        const conflict = conflicts[`${row.id}-${colKey}`];
                        const comments = cellComments[`${row.id}-${colKey}`];
                        // Heatmap overlay logic: highlight cells if heatmap is active
                        const heatmapClass = showHeatmap && typeof row[colKey] === 'number' && row[colKey] > 0 ? ' heatmap-cell' : '';
                        return hiddenColumns.includes(colKey) ? null : (
                            <td
                                key={colKey}
                                className={`px-6 py-4 text-sm font-medium${flashingCells && flashingCells[`${row.id}-${colKey}`] ? ' cell-flash' : ''}${error ? ' bg-red-100 border-2 border-red-400' : ''}${heatmapClass}`}
                                style={{ paddingLeft: colKey === 'entity' ? `${level * 24}px` : undefined, whiteSpace: 'pre-wrap', wordBreak: 'break-word', position: 'relative' }}
                                tabIndex={0}
                                onClick={() => { setSelectedCell({ rowId: row.id, colKey }); setFormulaBar(row[colKey]?.toString() || ''); }}
                            >
                                {selectedCell && selectedCell.rowId === row.id && selectedCell.colKey === colKey && formulaBar.startsWith('=')
                                    ? evaluateFormula(formulaBar, row, namedRanges)
                                    : row[colKey]}
                                {error && <span className="ml-2 text-xs text-red-600">⚠ {error}</span>}
                                {user && (
                                    <span className="absolute top-1 right-1 text-xs" style={{ color: user.color }} title={user.name}>{user.avatar}</span>
                                )}
                                {conflict && (
                                    <span className="absolute bottom-1 right-1 text-xs text-red-500 cursor-pointer" title="Resolve conflict" onClick={() => { setConflictCell({ rowId: row.id, colKey }); setConflictPanelOpen(true); }}>⚡</span>
                                )}
                                {comments && comments.length > 0 && (
                                    <span className="absolute bottom-1 left-1 text-xs text-blue-500 cursor-pointer" title="View comments" onClick={() => { setCommentCell({ rowId: row.id, colKey }); setCommentsPanelOpen(true); }}>💬</span>
                                )}
                                {colKey === 'entity' && row.children && row.children.length > 0 && (
                                    <button className="glass-button p-1 mr-2" onClick={() => handleExpandTreeRow(row.id)} aria-label={expandedTreeRows.includes(row.id) ? 'Collapse row' : 'Expand row'}>
                                        {expandedTreeRows.includes(row.id) ? '▼' : '▶'}
                                    </button>
                                )}
                                {colKey === 'entity' && draggedRowId === row.id && <span className="ml-2 text-xs text-[var(--color-accent)]">🟦 Dragging</span>}
                            </td>
                        );
                    })}
                </tr>
                {expandedTreeRows.includes(row.id) && row.children && row.children.length > 0 && renderTreeRows(row.children, level + 1)}
            </React.Fragment>
        ));
    }

    // Compute unique values for each column
    const uniqueValues = {};
    columnOrder.forEach(colKey => {
        uniqueValues[colKey] = Array.from(new Set(gridData.map(row => row[colKey])));
    });

    // Filtered data
    const filteredGridData = gridData.filter(row => {
        // ...existing filter logic...
        const passesFilters = columnOrder.every(colKey => {
            if (isDateCol(colKey, gridData) && dateFilters[colKey]) {
                const { from, to } = dateFilters[colKey];
                if (!isDateInRange(row[colKey], from, to)) return false;
            }
            if (!activeFilters[colKey] || activeFilters[colKey].size === 0) return true;
            return activeFilters[colKey].has(row[colKey]);
        });
        if (!passesFilters) return false;
        if (!searchText) return true;
        return columnOrder.some(colKey => (row[colKey] + '').toLowerCase().includes(searchText.toLowerCase()));
    });

    // Clipboard copy
    function handleCopy(e) {
        if (!selectedCell) return;
        const { rowId, colKey } = selectedCell;
        const row = gridData.find(r => r.id === rowId);
        if (!row) return;
        e.clipboardData.setData('text/plain', row[colKey]);
        e.preventDefault();
    }

    // Clipboard paste
    function handlePaste(e) {
        if (!selectedCell) return;
        const { rowId, colKey } = selectedCell;
        const rowIdx = gridData.findIndex(r => r.id === rowId);
        if (rowIdx === -1) return;
        const text = e.clipboardData.getData('text/plain');
        // Optimistic UI update
        setGridData(data => {
            const newData = [...data];
            newData[rowIdx] = { ...newData[rowIdx], [colKey]: text };
            return newData;
        });
        // Backend sync
        syncCellEdit(rowId, colKey, text);
        e.preventDefault();
    }

    // Keyboard shortcuts
    useEffect(() => {
        const el = gridRef.current;
        if (!el) return;
        el.addEventListener('copy', handleCopy);
        el.addEventListener('paste', handlePaste);
        return () => {
            el.removeEventListener('copy', handleCopy);
            el.removeEventListener('paste', handlePaste);
        };
    });

    // Main return
    return (
        <div ref={gridRef} tabIndex={0} style={{ outline: 'none' }} aria-label={accessibilityConfig.aria.ariaLabel} role={accessibilityConfig.aria.role} aria-rowcount={accessibilityConfig.aria.ariaRowCount} aria-colcount={accessibilityConfig.aria.ariaColCount} aria-multiselectable={accessibilityConfig.aria.ariaMultiselectable}>
            {renderGridControls()}
            {/* Flashing cell animation styles */}
            <style>{`
                .cell-flash {
                    animation: flash-cell 0.6s;
                }
                @keyframes flash-cell {
                    0% { background: #fef08a; }
                    60% { background: #fef08a; }
                    100% { background: inherit; }
                }
                .heatmap-cell {
                    background: linear-gradient(90deg, #fef08a 0%, #10b981 100%) !important;
                }
            `}</style>
            {/* Glassmorphism styles */}
            <style>{`
                .glass-card {
                    background: rgba(255,255,255,0.18);
                    backdrop-filter: blur(12px) saturate(1.2);
                    border-radius: 1.5rem;
                    box-shadow: 0 4px 32px rgba(0,0,0,0.08);
                    border: 1px solid rgba(255,255,255,0.22);
                }
                .glass-button {
                    background: rgba(255,255,255,0.22);
                    backdrop-filter: blur(6px);
                    border-radius: 0.5rem;
                    border: 1px solid rgba(255,255,255,0.18);
                }
                .glass-panel {
                    background: rgba(255,255,255,0.22);
                    backdrop-filter: blur(16px);
                    border-radius: 1rem;
                    border: 1px solid rgba(255,255,255,0.18);
                }
            `}</style>
            {/* Heatmap overlay */}
            {showHeatmap && (
                <div className="fixed top-0 left-0 w-full h-full pointer-events-none z-40">
                    <div className="absolute inset-0 flex items-center justify-center text-2xl font-bold text-green-700 opacity-80">Heatmap Overlay Active</div>
                </div>
            )}
            {/* Sankey fragment placeholder */}
            {showSankey && (
                <div className="fixed top-24 left-1/4 w-1/2 h-96 bg-white/90 rounded-lg shadow-lg p-6 z-50 border border-green-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Sankey Fragment (Demo)</h3>
                    <div className="flex items-center justify-center h-full">[Sankey visualization placeholder]</div>
                </div>
            )}
            {/* Export/Import panel UI */}
            {showExportPanel && (
                <div className="fixed top-32 left-1/3 w-1/3 h-80 bg-white/90 rounded-lg shadow-lg p-6 z-50 border border-blue-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Export / Import Data</h3>
                    <div className="flex flex-col gap-4">
                        <button className="glass-button" onClick={() => alert('Exported to CSV!')}>Export CSV</button>
                        <button className="glass-button" onClick={() => alert('Exported to Excel!')}>Export Excel</button>
                        <input type="file" className="glass-button" accept=".csv,.xlsx" onChange={() => alert('Import triggered!')} />
                        <div className="text-xs text-gray-500">Supports CSV/XLSX export/import. (Demo only)</div>
                    </div>
                    <button className="glass-button mt-4" onClick={() => setShowExportPanel(false)}>Close</button>
                </div>
            )}
            {/* Snapshot/Delta Log panel UI */}
            {showSnapshotPanel && (
                <div className="fixed top-40 left-1/3 w-1/3 h-80 bg-white/90 rounded-lg shadow-lg p-6 z-50 border border-yellow-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Snapshot / Delta Log</h3>
                    <div className="flex flex-col gap-4">
                        <button className="glass-button" onClick={() => alert('Snapshot taken!')}>Take Snapshot</button>
                        <button className="glass-button" onClick={() => alert('Delta log exported!')}>Export Delta Log</button>
                        <button className="glass-button" onClick={() => alert('Undo triggered!')}>Undo</button>
                        <button className="glass-button" onClick={() => alert('Redo triggered!')}>Redo</button>
                        <div className="text-xs text-gray-500">Snapshots and delta logs help with offline/conflict recovery. (Demo only)</div>
                    </div>
                    <button className="glass-button mt-4" onClick={() => setShowSnapshotPanel(false)}>Close</button>
                </div>
            )}
            {/* Guided Merge UX panel UI */}
            {showMergeUX && (
                <div className="fixed top-48 left-1/3 w-1/3 h-80 bg-white/90 rounded-lg shadow-lg p-6 z-50 border border-orange-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Guided Merge UX</h3>
                    <div className="flex flex-col gap-4">
                        <button className="glass-button" onClick={() => alert('Merge started!')}>Start Merge</button>
                        <button className="glass-button" onClick={() => alert('Conflict resolved!')}>Resolve Conflict</button>
                        <button className="glass-button" onClick={() => alert('Delta applied!')}>Apply Delta</button>
                        <div className="text-xs text-gray-500">Guided merge helps resolve offline/conflict edits. (Demo only)</div>
                    </div>
                    <button className="glass-button mt-4" onClick={() => setShowMergeUX(false)}>Close</button>
                </div>
            )}
            {/* RBAC/Privacy panel UI */}
            {showRBACPanel && (
                <div className="fixed top-56 left-1/3 w-1/3 h-80 bg-white/90 rounded-lg shadow-lg p-6 z-50 border border-purple-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">RBAC / Privacy Controls</h3>
                    <div className="flex flex-col gap-4">
                        <button className="glass-button" onClick={() => alert('Role assigned!')}>Assign Role</button>
                        <button className="glass-button" onClick={() => alert('Access revoked!')}>Revoke Access</button>
                        <button className="glass-button" onClick={() => alert('Privacy policy updated!')}>Update Privacy Policy</button>
                        <div className="text-xs text-gray-500">Role-based access and privacy controls. (Demo only)</div>
                    </div>
                    <button className="glass-button mt-4" onClick={() => setShowRBACPanel(false)}>Close</button>
                </div>
            )}
            {/* Audit/ESG panel UI */}
            {showAuditPanel && (
                <div className="fixed top-64 left-1/3 w-1/3 h-80 bg-white/90 rounded-lg shadow-lg p-6 z-50 border border-green-300" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Audit / ESG Panel</h3>
                    <div className="flex flex-col gap-4">
                        <button className="glass-button" onClick={() => alert('Audit log exported!')}>Export Audit Log</button>
                        <button className="glass-button" onClick={() => alert('ESG report generated!')}>Generate ESG Report</button>
                        <button className="glass-button" onClick={() => alert('Compliance checked!')}>Check Compliance</button>
                        <div className="text-xs text-gray-500">Audit and ESG controls for governance. (Demo only)</div>
                    </div>
                    <button className="glass-button mt-4" onClick={() => setShowAuditPanel(false)}>Close</button>
                </div>
            )}
            {/* Keyboard Shortcuts panel UI */}
            {showKeyboardShortcuts && (
                <div className="fixed top-72 left-1/3 w-1/3 h-80 bg-white/90 rounded-lg shadow-lg p-6 z-50 border border-blue-300" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Keyboard Shortcuts</h3>
                    <ul className="mb-4 text-sm">
                        <li><b>Arrow Keys:</b> Move cell selection</li>
                        <li><b>Enter:</b> Edit cell</li>
                        <li><b>Ctrl+C:</b> Copy cell</li>
                        <li><b>Ctrl+V:</b> Paste cell</li>
                        <li><b>Ctrl+Z:</b> Undo</li>
                        <li><b>Ctrl+Y:</b> Redo</li>
                        <li><b>Tab:</b> Next cell</li>
                        <li><b>Shift+Tab:</b> Previous cell</li>
                        <li><b>Esc:</b> Cancel edit</li>
                    </ul>
                    <div className="text-xs text-gray-500">Accessibility and keyboard navigation. (Demo only)</div>
                    <button className="glass-button mt-4" onClick={() => setShowKeyboardShortcuts(false)}>Close</button>
                </div>
            )}
            {/* High Contrast overlay UI */}
            {accessibilityConfig.highContrastTheme && (
                <div className="fixed top-0 left-0 w-full h-full pointer-events-none z-50">
                    <div className="absolute inset-0 flex items-center justify-center text-2xl font-bold text-black bg-yellow-200 opacity-90">High Contrast Theme Active</div>
                </div>
            )}
            {/* Plugin APIs panel UI */}
            {showPluginPanel && (
                <div className="fixed top-80 left-1/3 w-1/3 h-80 bg-white/90 rounded-lg shadow-lg p-6 z-50 border border-pink-300" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Plugin APIs & Extensibility</h3>
                    <div className="flex flex-col gap-4">
                        <button className="glass-button" onClick={() => alert('Demo plugin loaded!')}>Load Demo Plugin</button>
                        <button className="glass-button" onClick={() => alert('Plugin API docs opened!')}>View API Docs</button>
                        <button className="glass-button" onClick={() => alert('Plugin uninstalled!')}>Uninstall Plugin</button>
                        <div className="text-xs text-gray-500">Grid supports plugin APIs for custom logic, UI, and integrations. (Demo only)</div>
                    </div>
                    <button className="glass-button mt-4" onClick={() => setShowPluginPanel(false)}>Close</button>
                </div>
            )}
            {/* Theme Studio Controls */}
            <div className="mb-4 flex items-center gap-4">
                <label htmlFor="theme-select" className="font-semibold text-sm">Theme:</label>
                <select id="theme-select" value={theme} onChange={handleThemeChange} className="glass-button px-2 py-1 rounded">
                    {themeOptions.map(opt => (
                        <option key={opt} value={opt}>{opt.charAt(0).toUpperCase() + opt.slice(1)}</option>
                    ))}
                </select>
                {/* Future: Add custom theme editor UI here */}
            </div>
            {/* Tool Panel Sidebar */}
            {showToolPanel && (
                <div className="fixed top-0 left-0 w-80 h-full glass-panel shadow-xl z-50 p-6 flex flex-col gap-6 border-r border-[var(--color-border)]">
                    <div className="flex items-center justify-between mb-4">
                        <h2 className="text-lg font-semibold">Column & Filter Management</h2>
                        <button className="glass-button p-2 rounded-lg" onClick={() => setShowToolPanel(false)} aria-label="Close tool panel">
                            <X size={16} />
                        </button>
                    </div>
                    <div>
                        <h3 className="text-sm font-bold mb-2">Columns</h3>
                        {columnOrder.map(colKey => (
                            <div key={colKey} className="flex items-center gap-2 mb-2">
                                <input
                                    type="checkbox"
                                    checked={!hiddenColumns.includes(colKey)}
                                    onChange={() => handleToggleColumn(colKey)}
                                    aria-label={`Show/hide ${colKey}`}
                                />
                                <span>{colKey.charAt(0).toUpperCase() + colKey.slice(1)}</span>
                                <button className="glass-button p-1 text-xs" onClick={() => handleReorderColumn(colKey, 'left')} disabled={columnOrder.indexOf(colKey) === 0} title="Move Left">⬅️</button>
                                <button className="glass-button p-1 text-xs" onClick={() => handleReorderColumn(colKey, 'right')} disabled={columnOrder.indexOf(colKey) === columnOrder.length - 1} title="Move Right">➡️</button>
                            </div>
                        ))}
                    </div>
                    {/* Add filter controls here as needed */}
                </div>
            )}
            <button className="glass-button px-4 py-2 flex items-center gap-2" onClick={() => setShowToolPanel(true)}>
                <Filter size={16} />
                Tool Panel
            </button>
            <section className="glass-card rounded-2xl overflow-hidden mt-6" aria-labelledby="grid-title">
                <div className={`p-6 border-b border-[var(--color-border)] flex items-center justify-between theme-${theme}`}>
                    <h3 id="grid-title" className="font-semibold text-lg">Entity Performance Grid</h3>
                    <button className="glass-button p-2" onClick={() => setIsMaximized(!isMaximized)} aria-label={isMaximized ? "Restore" : "Maximize"}>{isMaximized ? <Minimize2 size={16} /> : <Maximize2 size={16} />}</button>
                </div>
                <div className="overflow-x-auto">
                    <table className="w-full text-left border-collapse relative">
                        <thead className="sticky top-0 z-10 glass shadow-sm">
                            <tr>
                                {columnOrder.map(colKey => (
                                    hiddenColumns.includes(colKey) ? null : (
                                        <th
                                            key={colKey}
                                            className={`px-6 py-4 text-xs font-bold uppercase text-[var(--color-text-muted)] cursor-move ${pinnedLeft.includes(colKey) ? 'bg-[var(--color-accent)]/10 sticky left-0 z-20' : ''} ${pinnedRight.includes(colKey) ? 'bg-[var(--color-accent)]/10 sticky right-0 z-20' : ''}`}
                                            style={{ width: columnWidths[colKey] }}
                                            onDoubleClick={() => autoSizeColumn(colKey)}
                                            title="Double-click to auto-size; Drag to reorder; Pin left/right"
                                            draggable
                                            onDragStart={() => handleDragStart(colKey)}
                                            onDragOver={e => handleDragOver(e, colKey)}
                                            onDrop={() => handleDrop(colKey)}
                                        >
                                            {colKey.charAt(0).toUpperCase() + colKey.slice(1)}
                                            <span className="ml-2">
                                                {!pinnedLeft.includes(colKey) && !pinnedRight.includes(colKey) && (
                                                    <>
                                                        <button className="glass-button p-1 text-xs" onClick={() => handlePinLeft(colKey)} title="Pin Left">⏪</button>
                                                        <button className="glass-button p-1 text-xs" onClick={() => handlePinRight(colKey)} title="Pin Right">⏩</button>
                                                    </>
                                                )}
                                                {pinnedLeft.includes(colKey) && (
                                                    <button className="glass-button p-1 text-xs" onClick={() => handleUnpin(colKey)} title="Unpin">❌</button>
                                                )}
                                                {pinnedRight.includes(colKey) && (
                                                    <button className="glass-button p-1 text-xs" onClick={() => handleUnpin(colKey)} title="Unpin">❌</button>
                                                )}
                                            </span>
                                        </th>
                                    )
                                ))}
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-[var(--color-border)]">
                            {renderTreeRows(filteredGridData, 0)}
                        </tbody>
                    </table>
                </div>
            </section>
            <button className="glass-button ml-2" onClick={() => setScriptingPanelOpen(v => !v)}>
                {scriptingPanelOpen ? 'Close Scripting' : 'Open Scripting'}
            </button>
            {scriptingPanelOpen && (
                <div className="fixed top-20 right-8 z-50 bg-white/90 rounded-lg shadow-lg p-6 w-[480px] border border-gray-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Grid Scripting (JS)</h3>
                    <textarea
                        className="w-full h-40 p-2 font-mono text-xs border rounded mb-2"
                        value={scriptCode}
                        onChange={e => setScriptCode(e.target.value)}
                        spellCheck={false}
                        aria-label="Script editor"
                    />
                    <div className="flex items-center gap-2">
                        <button className="glass-button" onClick={runScript}>Run Script</button>
                        <button className="glass-button" onClick={() => setScriptCode('// Example: gridData[0].amount += 100;\nreturn gridData;')}>Reset Example</button>
                        <span className="text-xs text-red-500 ml-2">{scriptError}</span>
                    </div>
                    <div className="text-xs text-gray-500 mt-2">Script runs in browser context. <b>gridData</b> is an array of row objects. <b>Return the new array.</b></div>
                </div>
            )}
            <button className="glass-button ml-2" onClick={() => setValidationPanelOpen(v => !v)}>
                {validationPanelOpen ? 'Close Validation' : 'Open Validation'}
            </button>
            {validationPanelOpen && (
                <div className="fixed top-40 right-8 z-50 bg-white/90 rounded-lg shadow-lg p-6 w-[420px] border border-gray-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Data Validation Rules</h3>
                    <div className="space-y-4">
                        {columnOrder.map(colKey => (
                            <div key={colKey} className="border-b pb-2 mb-2">
                                <div className="font-semibold text-xs mb-1">{colKey}</div>
                                <div className="flex gap-2 items-center text-xs">
                                    <label>Type:</label>
                                    <select value={validationRules[colKey]?.type || ''} onChange={e => handleRuleChange(colKey, 'type', e.target.value)}>
                                        <option value="">Any</option>
                                        <option value="number">Number</option>
                                    </select>
                                    <label>Required:</label>
                                    <input type="checkbox" checked={!!validationRules[colKey]?.required} onChange={e => handleRuleChange(colKey, 'required', e.target.checked)} />
                                    <label>Min:</label>
                                    <input type="number" className="w-16" value={validationRules[colKey]?.min ?? ''} onChange={e => handleRuleChange(colKey, 'min', e.target.value === '' ? undefined : Number(e.target.value))} />
                                    <label>Max:</label>
                                    <input type="number" className="w-16" value={validationRules[colKey]?.max ?? ''} onChange={e => handleRuleChange(colKey, 'max', e.target.value === '' ? undefined : Number(e.target.value))} />
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
            <button className="glass-button ml-2" onClick={() => setFilterPanelOpen(v => !v)}>
                {filterPanelOpen ? 'Close Filters' : 'Open Filters'}
            </button>
            {filterPanelOpen && (
                <div className="fixed top-48 right-8 z-50 bg-white/90 rounded-lg shadow-lg p-6 w-[420px] border border-gray-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Set Filters</h3>
                    <div className="space-y-4">
                        {columnOrder.map(colKey => (
                            <div key={colKey} className="border-b pb-2 mb-2">
                                <div className="font-semibold text-xs mb-1">{colKey}</div>
                                <div className="flex flex-wrap gap-2 text-xs">
                                    {uniqueValues[colKey].map(val => (
                                        <label key={val} className="flex items-center gap-1">
                                            <input
                                                type="checkbox"
                                                checked={activeFilters[colKey]?.has(val) || false}
                                                onChange={e => {
                                                    setActiveFilters(filters => {
                                                        const set = new Set(filters[colKey] || []);
                                                        if (e.target.checked) set.add(val); else set.delete(val);
                                                        return { ...filters, [colKey]: set };
                                                    });
                                                }}
                                            />
                                            {val === '' ? <span className="italic text-gray-400">(empty)</span> : val}
                                        </label>
                                    ))}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
            <button className="glass-button ml-2" onClick={() => setDateFilterPanelOpen(v => !v)}>
                {dateFilterPanelOpen ? 'Close Date Filters' : 'Date Filters'}
            </button>
            {dateFilterPanelOpen && (
                <div className="fixed top-56 right-8 z-50 bg-white/90 rounded-lg shadow-lg p-6 w-[420px] border border-gray-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Advanced Date Filters</h3>
                    <div className="space-y-4">
                        {columnOrder.filter(colKey => isDateCol(colKey, gridData)).map(colKey => (
                            <div key={colKey} className="border-b pb-2 mb-2">
                                <div className="font-semibold text-xs mb-1">{colKey}</div>
                                <div className="flex gap-2 items-center text-xs mb-1">
                                    <label>From:</label>
                                    <input type="date" value={dateFilters[colKey]?.from || ''} onChange={e => setDateFilters(f => ({ ...f, [colKey]: { ...f[colKey], from: e.target.value } }))} />
                                    <label>To:</label>
                                    <input type="date" value={dateFilters[colKey]?.to || ''} onChange={e => setDateFilters(f => ({ ...f, [colKey]: { ...f[colKey], to: e.target.value } }))} />
                                    <button className="glass-button px-2 py-1" onClick={() => setDateFilters(f => ({ ...f, [colKey]: { from: '', to: '' } }))}>Clear</button>
                                </div>
                                <div className="flex gap-2 text-xs">
                                    <button className="glass-button px-2 py-1" onClick={() => setDateFilters(f => ({ ...f, [colKey]: getDatePresetRange('today') }))}>Today</button>
                                    <button className="glass-button px-2 py-1" onClick={() => setDateFilters(f => ({ ...f, [colKey]: getDatePresetRange('this_week') }))}>This Week</button>
                                    <button className="glass-button px-2 py-1" onClick={() => setDateFilters(f => ({ ...f, [colKey]: getDatePresetRange('this_month') }))}>This Month</button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
            <div className="flex items-center gap-2 mb-2">
                <input
                    className="border rounded px-2 py-1 text-sm w-64"
                    type="text"
                    placeholder="Instant search..."
                    value={searchText}
                    onChange={e => setSearchText(e.target.value)}
                    aria-label="Instant search"
                />
                {searchText && <button className="glass-button px-2 py-1" onClick={() => setSearchText('')}>Clear</button>}
            </div>
            <button className="glass-button ml-2" onClick={() => setPivotPanelOpen(v => !v)}>
                {pivotPanelOpen ? 'Close Pivot' : 'Pivot Mode'}
            </button>
            {pivotPanelOpen && (
                <div className="fixed top-64 right-8 z-50 bg-white/90 rounded-lg shadow-lg p-6 w-[480px] border border-gray-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Pivot Table</h3>
                    <div className="flex flex-wrap gap-2 mb-2 text-xs">
                        <label>Row:</label>
                        <select value={pivotRow} onChange={e => setPivotRow(e.target.value)}>
                            <option value="">-- Select --</option>
                            {columnOrder.map(col => <option key={col} value={col}>{col}</option>)}
                        </select>
                        <label>Column:</label>
                        <select value={pivotCol} onChange={e => setPivotCol(e.target.value)}>
                            <option value="">-- Select --</option>
                            {columnOrder.map(col => <option key={col} value={col}>{col}</option>)}
                        </select>
                        <label>Value:</label>
                        <select value={pivotValue} onChange={e => setPivotValue(e.target.value)}>
                            <option value="">-- Select --</option>
                            {columnOrder.map(col => <option key={col} value={col}>{col}</option>)}
                        </select>
                        <label>Aggregation:</label>
                        <select value={pivotAgg} onChange={e => setPivotAgg(e.target.value)}>
                            <option value="sum">Sum</option>
                            <option value="avg">Average</option>
                            <option value="count">Count</option>
                        </select>
                        <button className="glass-button px-2 py-1" onClick={() => setPivotMode(v => !v)}>{pivotMode ? 'Exit Pivot' : 'Show Pivot'}</button>
                    </div>
                    {pivotData && (
                        <div className="overflow-x-auto mt-2">
                            <table className="min-w-full border text-xs">
                                <thead>
                                    <tr>
                                        <th className="border px-2 py-1 bg-gray-100">{pivotRow}</th>
                                        {pivotData.colKeys.map(ck => <th key={ck} className="border px-2 py-1 bg-gray-100">{ck}</th>)}
                                    </tr>
                                </thead>
                                <tbody>
                                    {pivotData.rowKeys.map((rk, i) => (
                                        <tr key={rk}>
                                            <td className="border px-2 py-1 font-semibold bg-gray-50">{rk}</td>
                                            {pivotData.result[i].map((v, j) => <td key={j} className="border px-2 py-1">{v}</td>)}
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                </div>
            )}
            <button className="glass-button ml-2" onClick={() => setCollabPanelOpen(v => !v)}>
                {collabPanelOpen ? 'Close Collaboration' : 'Collaboration'}
            </button>
            {collabPanelOpen && (
                <div className="fixed top-72 right-8 z-50 bg-white/90 rounded-lg shadow-lg p-6 w-[320px] border border-gray-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Online Users</h3>
                    <ul className="mb-2">
                        {onlineUsers.map(u => (
                            <li key={u.id} className="flex items-center gap-2 text-xs mb-1">
                                <span style={{ color: u.color }}>{u.avatar}</span>
                                <span>{u.name}</span>
                            </li>
                        ))}
                    </ul>
                    <div className="text-xs text-gray-500">Cells being edited are marked with user avatars.</div>
                </div>
            )}
            <button className="glass-button ml-2" onClick={() => setConflictPanelOpen(v => !v)}>
                {conflictPanelOpen ? 'Close Conflicts' : 'Conflict Resolution'}
            </button>
            {conflictPanelOpen && conflictCell && (
                <div className="fixed top-80 right-8 z-50 bg-white/90 rounded-lg shadow-lg p-6 w-[340px] border border-gray-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Resolve Conflict</h3>
                    <div className="mb-2 text-xs">Cell: {conflictCell.rowId} / {conflictCell.colKey}</div>
                    <div className="mb-2 text-xs">Users: {conflicts[`${conflictCell.rowId}-${conflictCell.colKey}`]?.map(uid => onlineUsers.find(u => u.id === uid)?.name).join(', ')}</div>
                    <div className="mb-2">
                        <label className="text-xs">Resolution:</label>
                        <select className="w-full p-1 border rounded" value={conflictResolution} onChange={e => setConflictResolution(e.target.value)}>
                            <option value="">-- Choose --</option>
                            <option value="keep">Keep current value</option>
                            <option value="merge">Merge values</option>
                        </select>
                    </div>
                    <button className="glass-button mt-2" onClick={() => { setConflictPanelOpen(false); setConflictResolution(''); }}>Resolve</button>
                </div>
            )}
            <button className="glass-button ml-2" onClick={() => setCommentsPanelOpen(v => !v)}>
                {commentsPanelOpen ? 'Close Comments' : 'Comments & Audit Trail'}
            </button>
            {commentsPanelOpen && commentCell && (
                <div className="fixed top-96 right-8 z-50 bg-white/90 rounded-lg shadow-lg p-6 w-[360px] border border-gray-200" style={{backdropFilter:'blur(8px)'}}>
                    <h3 className="font-bold mb-2">Cell Comments</h3>
                    <div className="mb-2 text-xs">Cell: {commentCell.rowId} / {commentCell.colKey}</div>
                    <ul className="mb-2">
                        {(cellComments[`${commentCell.rowId}-${commentCell.colKey}`] || []).map((c, i) => (
                            <li key={i} className="text-xs mb-1"><b>{c.user}</b>: {c.text} <span className="text-gray-400">({c.time})</span></li>
                        ))}
                    </ul>
                    <textarea className="w-full h-16 p-2 border rounded mb-2 text-xs" value={newComment} onChange={e => setNewComment(e.target.value)} placeholder="Add a comment..." />
                    <button className="glass-button" onClick={() => handleAddComment()}>Add Comment</button>
                    <h3 className="font-bold mt-4 mb-2">Audit Trail</h3>
                    <ul className="mb-2">
                        {auditTrail.filter(a => a.rowId === commentCell.rowId && a.colKey === commentCell.colKey).map((a, i) => (
                            <li key={i} className="text-xs mb-1">{a.time}: <b>{a.user}</b> changed from <span className="text-gray-400">{a.oldValue}</span> to <span className="text-gray-600">{a.newValue}</span></li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
}
