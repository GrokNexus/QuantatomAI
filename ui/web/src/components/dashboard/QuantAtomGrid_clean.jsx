"use client";

import React, { useState, useMemo, useEffect } from 'react';
import {
    TrendingUp, TrendingDown, ArrowUpRight, Filter, Download, MoreHorizontal, Plus, Eye, EyeOff, Maximize2, Minimize2, X, RotateCcw
} from 'lucide-react';
import { useTheme } from '@/context/ThemeContext';
import {
    AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer
} from 'recharts';

// Helper: format currency
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
// Simple sparkline component
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
const MOCK_CHART_DATA = [
    { name: 'Jan', actual: 4000, forecast: 4400 },
    { name: 'Feb', actual: 3000, forecast: 3200 },
    { name: 'Mar', actual: 2000, forecast: 2400 },
    { name: 'Apr', actual: 2780, forecast: 2900 },
    { name: 'May', actual: 1890, forecast: 2100 },
    { name: 'Jun', actual: 2390, forecast: 2500 },
    { name: 'Jul', actual: 3490, forecast: 3800 }
];

export function QuantAtomGrid() {
                // AI-Suggested Layouts & NLP Query
                const [aiLayout, setAiLayout] = useState(null);
                const [nlpQuery, setNlpQuery] = useState('');
                const [nlpResult, setNlpResult] = useState(null);
                const [aiLoading, setAiLoading] = useState(false);
                    {/* AI-Suggested Layouts & NLP Query Controls */}
                    <div className="glass-card p-4 rounded-xl mb-6 flex items-center gap-6">
                        <button
                            className="glass-button px-4 py-2 flex items-center gap-2"
                            onClick={async () => {
                                setAiLoading(true);
                                // Simulate AI layout suggestion
                                setTimeout(() => {
                                    setAiLayout({
                                        layout: 'AI-Optimized',
                                        description: 'Grid layout optimized for scenario comparison and anomaly detection.'
                                    });
                                    setAiLoading(false);
                                }, 1200);
                            }}
                            aria-busy={aiLoading}
                        >
                            <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <rect x="3" y="3" width="14" height="14" rx="3" fill="#6366f1" />
                                <circle cx="10" cy="10" r="4" fill="#a5b4fc" />
                            </svg>
                            {aiLoading ? 'Suggesting Layout...' : 'AI Layout Suggest'}
                        </button>
                        {aiLayout && (
                            <div className="ml-4 text-sm text-[var(--color-accent)]">
                                <b>{aiLayout.layout}</b>: {aiLayout.description}
                            </div>
                        )}
                        <div className="flex items-center gap-2">
                            <input
                                type="text"
                                value={nlpQuery}
                                onChange={e => setNlpQuery(e.target.value)}
                                placeholder="Ask grid: e.g. Show revenue by region"
                                className="glass-button px-2 py-1"
                                aria-label="NLP Query"
                            />
                            <button
                                className="glass-button px-3 py-2"
                                onClick={async () => {
                                    setAiLoading(true);
                                    // Simulate NLP query backend
                                    setTimeout(() => {
                                        setNlpResult(`Result for: "${nlpQuery}" (mocked)`);
                                        setAiLoading(false);
                                    }, 1000);
                                }}
                                aria-busy={aiLoading}
                            >
                                NLP Query
                            </button>
                        </div>
                        {nlpResult && (
                            <div className="ml-4 text-xs text-[var(--color-success)]">{nlpResult}</div>
                        )}
                    </div>
            // JIT Streaming & GPU Recalculation
            const [jitStreaming, setJitStreaming] = useState(false);
            const [gpuRecalc, setGpuRecalc] = useState(false);
            const [streamStatus, setStreamStatus] = useState('Idle');
            const [gpuStatus, setGpuStatus] = useState('Not Started');
                {/* JIT Streaming & GPU Recalculation Controls */}
                <div className="glass-card p-4 rounded-xl mb-6 flex items-center gap-6">
                    <div>
                        <label htmlFor="jitStreamingToggle" className="text-sm font-bold">JIT Streaming:</label>
                        <input
                            id="jitStreamingToggle"
                            type="checkbox"
                            checked={jitStreaming}
                            onChange={e => {
                                setJitStreaming(e.target.checked);
                                setStreamStatus(e.target.checked ? 'Streaming...' : 'Idle');
                                // Simulate backend streaming start/stop
                                if (e.target.checked) {
                                    setTimeout(() => setStreamStatus('Active'), 1000);
                                } else {
                                    setStreamStatus('Idle');
                                }
                            }}
                            className="ml-2"
                            aria-label="JIT Streaming"
                        />
                        <span className="ml-2 text-xs text-[var(--color-accent)]">{streamStatus}</span>
                    </div>
                    <div>
                        <label htmlFor="gpuRecalcToggle" className="text-sm font-bold">GPU Recalculation:</label>
                        <input
                            id="gpuRecalcToggle"
                            type="checkbox"
                            checked={gpuRecalc}
                            onChange={e => {
                                setGpuRecalc(e.target.checked);
                                setGpuStatus(e.target.checked ? 'Running...' : 'Not Started');
                                // Simulate backend GPU recalculation
                                if (e.target.checked) {
                                    setTimeout(() => setGpuStatus('Completed'), 1500);
                                } else {
                                    setGpuStatus('Not Started');
                                }
                            }}
                            className="ml-2"
                            aria-label="GPU Recalculation"
                        />
                        <span className="ml-2 text-xs text-[var(--color-success)]">{gpuStatus}</span>
                    </div>
                </div>
        // RBAC, Audit Logs, Tenant Controls
        const [userRole, setUserRole] = useState('viewer');
        const [tenantAIEnabled, setTenantAIEnabled] = useState(true);
        const [auditLogs, setAuditLogs] = useState([]);
        const [showAuditLogs, setShowAuditLogs] = useState(false);
            {/* RBAC & Tenant Controls UI */}
            <div className="glass-card p-4 rounded-xl mb-6 flex items-center gap-6">
                <div>
                    <label htmlFor="userRoleSelect" className="text-sm font-bold">User Role:</label>
                    <select
                        id="userRoleSelect"
                        value={userRole}
                        onChange={e => setUserRole(e.target.value)}
                        className="glass-button px-2 py-1 ml-2"
                        aria-label="User Role"
                    >
                        <option value="admin">Admin</option>
                        <option value="planner">Planner</option>
                        <option value="viewer">Viewer</option>
                    </select>
                </div>
                <div>
                    <label htmlFor="tenantAIEnabledToggle" className="text-sm font-bold">Tenant AI Enabled:</label>
                    <input
                        id="tenantAIEnabledToggle"
                        type="checkbox"
                        checked={tenantAIEnabled}
                        onChange={e => setTenantAIEnabled(e.target.checked)}
                        className="ml-2"
                        aria-label="Tenant AI Enabled"
                    />
                </div>
                <button
                    className="glass-button px-3 py-2"
                    onClick={() => setShowAuditLogs(!showAuditLogs)}
                    aria-pressed={showAuditLogs}
                >
                    {showAuditLogs ? 'Hide Audit Logs' : 'Show Audit Logs'}
                </button>
            </div>
            {/* Audit Logs Modal */}
            {showAuditLogs && (
                <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50" role="dialog" aria-modal="true" aria-labelledby="audit-logs-title">
                    <div className="bg-white dark:bg-[var(--color-bg-primary)] p-6 rounded-xl shadow-xl w-[600px] max-h-[80vh] overflow-auto">
                        <h2 id="audit-logs-title" className="text-lg font-bold mb-4">Audit Logs</h2>
                        <table className="w-full text-sm">
                            <thead>
                                <tr>
                                    <th>User</th>
                                    <th>Action</th>
                                    <th>Prompt</th>
                                    <th>Timestamp</th>
                                </tr>
                            </thead>
                            <tbody>
                                {auditLogs.length === 0 ? (
                                    <tr><td colSpan={4} className="text-center text-[var(--color-text-muted)]">No audit logs found.</td></tr>
                                ) : (
                                    auditLogs.map((log, i) => (
                                        <tr key={i}>
                                            <td>{log.user}</td>
                                            <td>{log.action}</td>
                                            <td>{log.prompt}</td>
                                            <td>{log.timestamp}</td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                        <button className="glass-button px-3 py-2 mt-4" onClick={() => setShowAuditLogs(false)}>Close</button>
                    </div>
                </div>
            )}
    const dimensionMembers = {
        entity: ["North America Ops", "EMEA Logistics"],
        status: ["On Track", "At Risk"],
    };
    const [selectedFilters, setSelectedFilters] = useState({ entity: [], status: [] });
    const [dateRange, setDateRange] = useState({ start: '', end: '' });
    const [searchText, setSearchText] = useState('');
    const [pivotMode, setPivotMode] = useState(false);
    const [displayProfile, setDisplayProfile] = useState('default');
    const gridDataWithDate = MOCK_GRID_DATA.map((row, i) => ({ ...row, date: `2025-0${i + 1}-01` }));
    const filteredGridData = useMemo(() => {
        return gridDataWithDate.filter(row => {
            if (selectedFilters.entity.length && !selectedFilters.entity.includes(row.entity)) return false;
            if (selectedFilters.status.length && !selectedFilters.status.includes(row.status)) return false;
            if (dateRange.start && row.date < dateRange.start) return false;
            if (dateRange.end && row.date > dateRange.end) return false;
            if (searchText) {
                const rowString = Object.values(row).join(' ').toLowerCase();
                if (!rowString.includes(searchText.toLowerCase())) return false;
            }
            return true;
        });
    }, [selectedFilters, dateRange, searchText]);
    const { theme, toggleTheme } = useTheme();
    const [showForecast, setShowForecast] = useState(true);
    const [isMaximized, setIsMaximized] = useState(false);
    const [isMinimized, setIsMinimized] = useState(false);
    const [gridView, setGridView] = useState('table');
    const [showFilterPanel, setShowFilterPanel] = useState(false);
    const [conflictModal, setConflictModal] = useState(null);
    const [deltaLog, setDeltaLog] = useState([]);
    const [vrMode, setVrMode] = useState(false);
    const [voiceQuery, setVoiceQuery] = useState('');
    const [isListening, setIsListening] = useState(false);
    useEffect(() => {
        if (!deltaLog.length) return;
        const grouped = deltaLog.reduce((acc, entry) => {
            const key = entry.rowId + '-' + entry.field;
            acc[key] = acc[key] || [];
            acc[key].push(entry.value);
            return acc;
        }, {});
        for (const key in grouped) {
            if (grouped[key].length > 1) {
                const [rowId, field] = key.split('-');
                setConflictModal({ rowId, field, values: grouped[key] });
                break;
            }
        }
    }, [deltaLog]);
    function resolveConflict(value) { setConflictModal(null); }
    const [offlineSnapshot, setOfflineSnapshot] = useState([]);
    function handleCellEdit(rowId, field, value) {
        setDeltaLog(log => [...log, { rowId, field, value, timestamp: Date.now() }]);
    }
    function takeSnapshot() { setOfflineSnapshot(formattedGridData); }
    function restoreSnapshot() { alert('Restored snapshot!'); }
    const entitySparklines = {
        'North America Ops': [4, 5, 6, 7, 8, 7, 6],
        'EMEA Logistics': [-2, -3, -4, -5, -4, -3, -2],
    };
    const formattedGridData = useMemo(() =>
        filteredGridData.map(row => ({
            ...row,
            budgetFormatted: formatCurrency(row.budget),
            actualFormatted: formatCurrency(row.actual),
            varianceFormatted: formatPercent(row.variance),
        })), [filteredGridData]);

    // Accessibility: keyboard shortcuts
    useEffect(() => {
        function handleKeyDown(e) {
            // Example shortcuts
            if (e.ctrlKey && e.key === 'f') {
                setShowFilterPanel(true);
                e.preventDefault();
            }
            if (e.ctrlKey && e.key === 'e') {
                alert('Export Excel');
                e.preventDefault();
            }
            if (e.ctrlKey && e.key === 'p') {
                alert('Export PDF');
                e.preventDefault();
            }
            if (e.ctrlKey && e.key === 'm') {
                setIsMinimized(true);
                e.preventDefault();
            }
            if (e.ctrlKey && e.key === 'x') {
                setIsMaximized(true);
                e.preventDefault();
            }
        }
        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, []);

    if (isMinimized) {
        return (
            <div>
                {conflictModal && (
                    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50" role="dialog" aria-modal="true" aria-labelledby="conflict-title">
                        <div className="bg-white dark:bg-[var(--color-bg-primary)] p-6 rounded-xl shadow-xl w-96">
                            <h2 id="conflict-title" className="text-lg font-bold mb-4">Resolve Conflict</h2>
                            <p className="mb-2">Conflicting edits detected for <b>{conflictModal.field}</b> in row <b>{conflictModal.rowId}</b>.</p>
                            <div className="mb-4">
                                {conflictModal.values.map((v, i) => (
                                    <button key={i} className="glass-button px-3 py-2 m-1" onClick={() => resolveConflict(v)}>
                                        Choose: {v}
                                    </button>
                                ))}
                            </div>
                            <button className="glass-button px-3 py-2 mt-2" onClick={() => setConflictModal(null)}>Cancel</button>
                        </div>
                    </div>
                )}
                <div className="fixed bottom-4 right-4 z-50">
                    <button
                        onClick={() => setIsMinimized(false)}
                        className="glass-button px-4 py-2 flex items-center gap-2"
                        aria-label="Restore minimized dashboard"
                        tabIndex={0}
                    >
                        <RotateCcw size={16} />
                        Restore Dashboard
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className={`space-y-8 ${isMaximized ? 'fixed inset-0 bg-[var(--color-bg-primary)] z-40 p-6 overflow-auto' : ''}`}>
            {/* VR Grid & Voice Query Controls */}
            <div className="flex items-center gap-4 mb-4">
                <button
                    className={`glass-button px-4 py-2 flex items-center gap-2 ${vrMode ? 'bg-[var(--color-accent)]/20' : ''}`}
                    onClick={() => setVrMode(!vrMode)}
                    aria-pressed={vrMode}
                >
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <rect x="2" y="6" width="16" height="8" rx="2" fill="#6366f1" />
                        <circle cx="6" cy="10" r="2" fill="#a5b4fc" />
                        <circle cx="14" cy="10" r="2" fill="#a5b4fc" />
                    </svg>
                    VR Grid
                </button>
                <button
                    className={`glass-button px-4 py-2 flex items-center gap-2 ${isListening ? 'bg-[var(--color-accent)]/20' : ''}`}
                    onClick={() => {
                        if (!isListening) {
                            setIsListening(true);
                            if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
                                const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
                                const recognition = new SpeechRecognition();
                                recognition.continuous = false;
                                recognition.interimResults = false;
                                recognition.lang = 'en-US';
                                recognition.onresult = event => {
                                    setVoiceQuery(event.results[0][0].transcript);
                                    setIsListening(false);
                                };
                                recognition.onerror = () => setIsListening(false);
                                recognition.onend = () => setIsListening(false);
                                recognition.start();
                            } else {
                                alert('Speech recognition not supported in this browser.');
                                setIsListening(false);
                            }
                        } else {
                            setIsListening(false);
                        }
                    }}
                    aria-pressed={isListening}
                    aria-label="Voice to Query"
                >
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <rect x="8" y="3" width="4" height="10" rx="2" fill="#10b981" />
                        <rect x="6" y="13" width="8" height="2" rx="1" fill="#a7f3d0" />
                    </svg>
                    {isListening ? 'Listening...' : 'Voice Query'}
                </button>
                {voiceQuery && (
                    <span className="ml-2 text-sm font-medium text-[var(--color-accent)]">Query: {voiceQuery}</span>
                )}
            </div>
            {/* VR Grid Mode Placeholder */}
            {vrMode && (
                <div className="glass-card p-6 rounded-2xl mb-6 flex flex-col items-center justify-center" style={{ minHeight: 300 }}>
                    <h2 className="text-xl font-bold mb-2">VR Grid Prototype</h2>
                    <p className="mb-4">Immersive 3D grid navigation coming soon. (WebXR/Three.js integration planned)</p>
                    <div className="w-full h-64 bg-gradient-to-br from-indigo-200 via-indigo-100 to-white rounded-xl flex items-center justify-center">
                        <span className="text-lg text-indigo-600">[VR Grid Canvas Placeholder]</span>
                    </div>
                </div>
            )}
            {/* Filter Panel Sidebar */}
            {showFilterPanel && (
                <div className="fixed top-0 right-0 w-80 h-full bg-white dark:bg-[var(--color-bg-primary)] shadow-xl z-50 p-6 flex flex-col gap-6 border-l border-[var(--color-border)]">
                    <div className="flex items-center justify-between mb-4">
                        <h2 className="text-lg font-semibold">Filters</h2>
                        <button className="glass-button p-2 rounded-lg" onClick={() => setShowFilterPanel(false)} aria-label="Close filter panel">
                            <X size={16} />
                        </button>
                    </div>
                    {/* Display Profile Selector */}
                    <div className="mb-4">
                        <label className="text-sm font-bold mb-2" htmlFor="displayProfileSelect">Display Profile</label>
                        <select id="displayProfileSelect" value={displayProfile} onChange={e => setDisplayProfile(e.target.value)} className="glass-button px-2 py-1 w-full" title="Display Profile" aria-label="Display Profile">
                            <option value="default">Default</option>
                            <option value="bold">Bold</option>
                        </select>
                    </div>
                    {/* Pivot Mode Toggle */}
                    <div className="mb-4">
                        <label className="flex items-center gap-2">
                            <input
                                type="checkbox"
                                checked={pivotMode}
                                onChange={e => setPivotMode(e.target.checked)}
                            />
                            <span>Pivot Mode</span>
                        </label>
                    </div>
                    {/* Instant Search */}
                    <div>
                        <input
                            type="text"
                            value={searchText}
                            onChange={e => setSearchText(e.target.value)}
                            placeholder="Search grid..."
                            className="glass-button px-2 py-1 w-full mb-4"
                            aria-label="Instant search"
                        />
                    </div>
                    {/* Entity Filter */}
                    <div>
                        <h3 className="text-sm font-bold mb-2">Entity</h3>
                        {dimensionMembers.entity.map(member => (
                            <label key={member} className="flex items-center gap-2 mb-2">
                                <input
                                    type="checkbox"
                                    checked={selectedFilters.entity.includes(member)}
                                    onChange={e => {
                                        setSelectedFilters(prev => {
                                            const checked = e.target.checked;
                                            const updated = checked
                                                ? [...prev.entity, member]
                                                : prev.entity.filter(m => m !== member);
                                            return { ...prev, entity: updated };
                                        });
                                    }}
                                />
                                <span>{member}</span>
                            </label>
                        ))}
                    </div>
                    {/* Status Filter */}
                    <div>
                        <h3 className="text-sm font-bold mb-2">Status</h3>
                        {dimensionMembers.status.map(member => (
                            <label key={member} className="flex items-center gap-2 mb-2">
                                <input
                                    type="checkbox"
                                    checked={selectedFilters.status.includes(member)}
                                    onChange={e => {
                                        setSelectedFilters(prev => {
                                            const checked = e.target.checked;
                                            const updated = checked
                                                ? [...prev.status, member]
                                                : prev.status.filter(m => m !== member);
                                            return { ...prev, status: updated };
                                        });
                                    }}
                                />
                                <span>{member}</span>
                            </label>
                        ))}
                    </div>
                    {/* Date Range Filter */}
                    <div>
                        <h3 className="text-sm font-bold mb-2">Date Range</h3>
                        <div className="flex gap-2 items-center">
                            <input
                                type="date"
                                value={dateRange.start}
                                onChange={e => setDateRange(r => ({ ...r, start: e.target.value }))}
                                className="glass-button px-2 py-1"
                                aria-label="Start date"
                            />
                            <span>to</span>
                            <input
                                type="date"
                                value={dateRange.end}
                                onChange={e => setDateRange(r => ({ ...r, end: e.target.value }))}
                                className="glass-button px-2 py-1"
                                aria-label="End date"
                            />
                        </div>
                    </div>
                    <button className="glass-button mt-4" onClick={() => setShowFilterPanel(false)}>Apply Filters</button>
                </div>
            )}

            {/* Header with Controls */}
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-semibold tracking-tight text-[var(--color-text-primary)]">Executive Dashboard</h1>
                    <p className="text-[var(--color-text-muted)] mt-1">Real-time planning and performance across 24 entities</p>
                </div>
                <div className="flex items-center gap-3 flex-wrap">
                    <button onClick={toggleTheme} className="glass-button p-2 rounded-lg" aria-label="Toggle theme">
                        {theme === 'dark' ? 'Light' : 'Dark'}
                    </button>
                    <button onClick={() => setGridView(gridView === 'table' ? 'cards' : 'table')} className="glass-button px-3 py-2 text-sm flex items-center gap-2" aria-pressed={gridView === 'cards'}>
                        {gridView === 'table' ? <Eye size={16} /> : <EyeOff size={16} />}
                        {gridView === 'table' ? 'Card View' : 'Table View'}
                        <span className="sr-only">Toggle grid view</span>
                    </button>
                    <div className="flex items-center gap-1">
                        <button onClick={() => setIsMaximized(!isMaximized)} className="glass-button p-2 rounded-lg" aria-label={isMaximized ? "Restore" : "Maximize"}>
                            {isMaximized ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
                        </button>
                        <button onClick={() => setIsMinimized(true)} className="glass-button p-2 rounded-lg text-[var(--color-warning)] hover:text-[var(--color-error)]" aria-label="Minimize dashboard">
                            <Minimize2 size={16} />
                        </button>
                        <button className="glass-button p-2 rounded-lg text-[var(--color-error)] hover:opacity-80" aria-label="Close dashboard">
                            <X size={16} />
                        </button>
                    </div>
                    <button className="glass-button px-4 py-2 flex items-center gap-2">
                        <Filter size={16} />
                        Filters
                    </button>
                    <button className="glass-button px-4 py-2 flex items-center gap-2" onClick={() => setShowFilterPanel(true)}>
                        <Filter size={16} />
                        Open Filter Panel
                    </button>
                    <button className="glass-button px-4 py-2 flex items-center gap-2" onClick={() => alert("Export Excel") } aria-label="Export to Excel">
                        <Download size={16} />
                        Export Excel
                    </button>
                    <button className="glass-button px-4 py-2 flex items-center gap-2" onClick={() => alert("Export PDF") } aria-label="Export to PDF">
                        <Download size={16} />
                        Export PDF
                    </button>
                    <button className="bg-[var(--color-accent)] text-white px-4 py-2 rounded-lg flex items-center gap-2 hover:opacity-90 transition-opacity shadow-md">
                        <Plus size={16} />
                        New Scenario
                    </button>
                </div>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6" role="region" aria-label="Key Performance Indicators">
                {[
                    { title: "Total Budget", value: "$8.05M", change: "+12.5%", positive: true },
                    { title: "Actual Spend", value: "$7.61M", change: "-2.4%", positive: false },
                    { title: "Forecast Accuracy", value: "94.2%", change: "+1.2%", positive: true },
                    { title: "Active Scenarios", value: "24", change: "+4", positive: true },
                ].map((stat, i) => (
                    <div key={i} className="glass-card p-6 rounded-2xl" role="status" aria-label={`${stat.title}: ${stat.value}, change ${stat.change}`}>
                        <p className="text-sm font-medium text-[var(--color-text-muted)]">{stat.title}</p>
                        <div className="mt-2 flex items-baseline justify-between">
                            <h2 className="text-2xl font-bold tracking-tight">{stat.value}</h2>
                            <div className={`flex items-center text-sm font-bold ${stat.positive ? 'text-[var(--color-success)]' : 'text-[var(--color-error)]'}`}>{stat.change}</div>
                        </div>
                    </div>
                ))}
            </div>

            {/* Charts & Grid Section */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Main Chart */}
                <section className="lg:col-span-2 glass-card p-6 rounded-2xl" aria-labelledby="chart-title">
                    <div className="flex flex-col sm:flex-row sm:items-center justify-between mb-6 gap-4">
                        <h3 id="chart-title" className="font-semibold text-lg">Performance vs Forecast</h3>
                        <div className="flex items-center gap-4">
                            <button onClick={() => setShowForecast(!showForecast)} className={`glass-button px-3 py-1.5 text-sm flex items-center gap-2 ${showForecast ? 'bg-[var(--color-accent)]/20' : ''}`} aria-pressed={showForecast}>
                                {showForecast ? <Eye size={16} /> : <EyeOff size={16} />}
                                Forecast
                            </button>
                            <select className="glass-button px-3 py-1.5 text-sm" aria-label="Time range">
                                <option>Last 7 Months</option>
                                <option>Last Year</option>
                            </select>
                        </div>
                    </div>
                    <div className="h-80 w-full">
                        <ResponsiveContainer>
                            <AreaChart data={MOCK_CHART_DATA} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                                <defs>
                                    <linearGradient id="actualGradient" x1="0" y1="0" x2="0" y2="1">
                                        <stop offset="5%" stopColor="var(--color-accent)" stopOpacity={0.3} />
                                        <stop offset="95%" stopColor="var(--color-accent)" stopOpacity={0} />
                                    </linearGradient>
                                </defs>
                                <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" vertical={false} />
                                <XAxis dataKey="name" stroke="var(--color-text-muted)" tick={{ fontSize: 12 }} />
                                <YAxis stroke="var(--color-text-muted)" tick={{ fontSize: 12 }} />
                                <Tooltip content={<div className="p-2 bg-white dark:bg-[var(--color-bg-primary)] rounded shadow text-xs">Custom Tooltip</div>} />
                                <Area type="monotone" dataKey="actual" stroke="var(--color-accent)" fill="url(#actualGradient)" strokeWidth={3} />
                                {showForecast && (
                                    <Area type="monotone" dataKey="forecast" stroke="#94a3b8" strokeDasharray="5 5" fill="transparent" strokeWidth={2} />
                                )}
                            </AreaChart>
                        </ResponsiveContainer>
                    </div>
                </section>

                {/* Entity Health */}
                <section className="glass-card p-6 rounded-2xl" aria-labelledby="health-title">
                    <h3 id="health-title" className="font-semibold text-lg mb-6">Entity Health</h3>
                    <div className="space-y-6">
                        {MOCK_GRID_DATA.slice(0, 4).map(item => (
                            <div key={item.id} className="space-y-2">
                                <div className="flex justify-between items-center">
                                    <span className="font-medium text-[var(--color-text-primary)]">{item.entity}</span>
                                    <span className="text-sm font-bold text-[var(--color-text-muted)]">{item.health}%</span>
                                </div>
                                <div className="h-3 glass rounded-full overflow-hidden shadow-inner">
                                    <div className={`h-full rounded-full transition-all duration-1000 ${getHealthColor(item.health, theme)} health-bar`} style={{ width: `${item.health}%` }}></div>
                                </div>
                            </div>
                        ))}
                    </div>
                    <button className="w-full mt-6 glass-button py-2 text-sm font-medium">View All Entities</button>
                </section>
            </div>

            {/* Entity Performance Grid */}
            <section className="glass-card rounded-2xl overflow-hidden" aria-labelledby="grid-title">
                <div className="p-6 border-b border-[var(--color-border)] flex items-center justify-between">
                    <h3 id="grid-title" className="font-semibold text-lg">Entity Performance Grid</h3>
                    <div className="flex items-center gap-2">
                        <button className="glass-button px-3 py-2" onClick={takeSnapshot} aria-label="Take offline snapshot">Snapshot</button>
                        <button className="glass-button px-3 py-2" onClick={restoreSnapshot} aria-label="Restore snapshot">Restore</button>
                        <button className="glass-button p-2" aria-label="Restore grid"><RotateCcw size={16} /></button>
                        <button className="glass-button p-2" aria-label={isMaximized ? "Restore" : "Maximize"}>{isMaximized ? <Minimize2 size={16} /> : <Maximize2 size={16} />}</button>
                        <button className="glass-button p-2 text-[var(--color-warning)]" aria-label="Minimize"><Minimize2 size={16} /></button>
                        <button className="glass-button p-2 text-[var(--color-error)]" aria-label="Close"><X size={16} /></button>
                    </div>
                </div>
                <div className="overflow-x-auto">
                    {pivotMode ? (
                        <table className="w-full text-left border-collapse relative">
                            <thead className="sticky top-0 z-10 glass shadow-sm">
                                <tr>
                                    <th className="px-6 py-4 text-xs font-bold uppercase text-[var(--color-text-muted)]">Entity</th>
                                    <th className="px-6 py-4 text-xs font-bold uppercase text-[var(--color-text-muted)]">Status</th>
                                    <th className="px-6 py-4 text-xs font-bold uppercase text-[var(--color-text-muted)]">Budget</th>
                                    <th className="px-6 py-4 text-xs font-bold uppercase text-[var(--color-text-muted)]">Actual</th>
                                    <th className="px-6 py-4 text-xs font-bold uppercase text-[var(--color-text-muted)]">Variance</th>
                                    <th className="px-6 py-4 text-xs font-bold uppercase text-[var(--color-text-muted)]">Trend</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-[var(--color-border)]">
                                {dimensionMembers.status.map(status => (
                                    <React.Fragment key={status}>
                                        <tr className="bg-[var(--color-surface-base)]/30">
                                            <td colSpan={6} className="px-6 py-2 text-xs font-bold text-[var(--color-accent)]">{status}</td>
                                        </tr>
                                        {formattedGridData.filter(row => row.status === status).map(row => (
                                            <tr key={row.id} className="hover:bg-white/5 transition-colors group">
                                                <td className="px-6 py-4 text-sm font-medium">{row.entity}</td>
                                                <td className="px-6 py-4 text-sm">
                                                    <input type="text" value={row.status} onChange={e => handleCellEdit(row.id, 'status', e.target.value)} className="w-24 px-1 py-0.5 border border-transparent group-hover:border-[var(--glass-border)] bg-transparent rounded focus:bg-[var(--glass-bg)]" title="Status" placeholder="Status" aria-label="Status" />
                                                </td>
                                                <td className="px-6 py-4 text-sm">
                                                    <input type="number" value={row.budget} onChange={e => handleCellEdit(row.id, 'budget', e.target.value)} className="w-24 px-1 py-0.5 border border-transparent group-hover:border-[var(--glass-border)] bg-transparent rounded focus:bg-[var(--glass-bg)]" title="Budget" placeholder="Budget" aria-label="Budget" />
                                                </td>
                                                <td className="px-6 py-4 text-sm">
                                                    <input type="number" value={row.actual} onChange={e => handleCellEdit(row.id, 'actual', e.target.value)} className="w-24 px-1 py-0.5 border border-transparent group-hover:border-[var(--glass-border)] bg-transparent rounded focus:bg-[var(--glass-bg)]" title="Actual" placeholder="Actual" aria-label="Actual" />
                                                </td>
                                                <td className={`px-6 py-4 text-sm font-medium ${getVarianceClass(row.variance, displayProfile)} ${row.variance > 0 ? 'variance-positive' : row.variance < 0 ? 'variance-negative' : ''}`}>{row.varianceFormatted}</td>
                                                <td className="px-6 py-4">
                                                    <Sparkline data={entitySparklines[row.entity] || []} />
                                                </td>
                                            </tr>
                                        ))}
                                    </React.Fragment>
                                ))}
                            </tbody>
                        </table>
                    ) : null}
                </div>
            </section>
        </div>
    );
}
