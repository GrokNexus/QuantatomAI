"use client";

import React, { useState } from 'react';
import { useShellStore, CalendarConfig, ModuleId } from '../../store/useShellStore';

export const ApplicationWizard: React.FC = () => {
    const { setActiveApplication, setActiveModule } = useShellStore();

    const [appName, setAppName] = useState('');
    const [appDesc, setAppDesc] = useState('');

    const [dimensions, setDimensions] = useState<string[]>(['Time', 'Scenario', 'Version', 'Currency', 'Account']);
    const availableDimensions = ['Region', 'Product', 'Department', 'Project', 'Channel'];

    const [calendar, setCalendar] = useState<CalendarConfig>({
        calendarType: 'Monthly',
        fiscalStartMonth: 'January',
        currentFiscalYear: new Date().getFullYear(),
        pastYearsCount: 2,
        futureYearsCount: 5
    });

    const toggleDimension = (dim: string) => {
        setDimensions(prev => prev.includes(dim) ? prev.filter(d => d !== dim) : [...prev, dim]);
    };

    const handleCreate = () => {
        if (!appName) {
            alert("Application Name is required.");
            return;
        }

        setActiveApplication({
            id: `app-${Date.now()}`,
            name: appName,
            description: appDesc,
            dimensions,
            calendar
        });

        // Route to the grid automatically upon creation
        setActiveModule(ModuleId.GRID);
    };

    return (
        <div style={{
            display: 'flex', flexDirection: 'column',
            width: '100%', height: '100%',
            overflowY: 'auto', padding: 'var(--space-8)',
            alignItems: 'center', backgroundColor: 'transparent'
        }}>
            <div style={{
                maxWidth: '800px', width: '100%',
                backgroundColor: 'var(--color-surface-elevate)',
                border: '1px solid var(--glass-border-color)',
                borderRadius: 'var(--radius-lg)',
                boxShadow: 'var(--shadow-elevation-3)',
                padding: 'var(--space-8)',
                display: 'flex', flexDirection: 'column', gap: 'var(--space-6)'
            }}>
                <div>
                    <h2 style={{ margin: 0, font: 'var(--text-h2)', color: 'var(--color-text-main)' }}>Create New Application</h2>
                    <p style={{ margin: '8px 0 0', font: 'var(--text-body)', color: 'var(--color-text-dim)' }}>
                        Define the dimensional bounds and time horizons for your financial matrix.
                    </p>
                </div>

                <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
                    {/* Basic Info */}
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
                        <label style={{ font: 'var(--text-micro)', color: 'rgba(255,255,255,0.6)', textTransform: 'uppercase' }}>Application Name</label>
                        <input
                            value={appName} onChange={e => setAppName(e.target.value)}
                            placeholder="e.g. Corporate FP&A 2026"
                            style={{
                                padding: 'var(--space-3)', backgroundColor: 'rgba(0,0,0,0.3)',
                                border: '1px solid var(--glass-border-color)', borderRadius: 'var(--radius-sm)',
                                color: 'white', font: 'var(--text-body)', outline: 'none'
                            }}
                        />
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
                        <label style={{ font: 'var(--text-micro)', color: 'rgba(255,255,255,0.6)', textTransform: 'uppercase' }}>Description (Optional)</label>
                        <input
                            value={appDesc} onChange={e => setAppDesc(e.target.value)}
                            placeholder="Master planning model for Q3 projections..."
                            style={{
                                padding: 'var(--space-3)', backgroundColor: 'rgba(0,0,0,0.3)',
                                border: '1px solid var(--glass-border-color)', borderRadius: 'var(--radius-sm)',
                                color: 'white', font: 'var(--text-body)', outline: 'none'
                            }}
                        />
                    </div>

                    <hr style={{ borderColor: 'var(--glass-border-color)', margin: 'var(--space-4) 0', borderBottom: 'none' }} />

                    {/* Dimensionality Mapping */}
                    <div>
                        <h4 style={{ margin: '0 0 16px', color: 'var(--color-text-main)', font: 'var(--text-h4)' }}>Dimensionality Sync (DB Schema)</h4>
                        <p style={{ margin: '0 0 16px', color: 'var(--color-text-dim)', fontSize: '12px' }}>
                            Mandatory dimensions (Time, Scenario, Version, Currency) are locked. Select structural blocks:
                        </p>
                        <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
                            {['Time', 'Scenario', 'Version', 'Currency', 'Account'].map(d => (
                                <div key={d} style={{ padding: '8px 16px', backgroundColor: 'rgba(59, 130, 246, 0.2)', border: '1px solid var(--color-primary)', borderRadius: '999px', color: 'var(--color-primary)', fontSize: '13px' }}>
                                    🔒 {d}
                                </div>
                            ))}
                            {availableDimensions.map(d => {
                                const active = dimensions.includes(d);
                                return (
                                    <button
                                        key={d}
                                        onClick={() => toggleDimension(d)}
                                        style={{
                                            padding: '8px 16px',
                                            backgroundColor: active ? 'rgba(16, 185, 129, 0.2)' : 'transparent',
                                            border: active ? '1px solid var(--color-positive)' : '1px solid var(--glass-border-color)',
                                            borderRadius: '999px',
                                            color: active ? 'var(--color-positive)' : 'var(--color-text-dim)',
                                            fontSize: '13px', cursor: 'pointer', transition: 'all 0.1s'
                                        }}>
                                        {active ? '✓ ' : '+ '} {d}
                                    </button>
                                );
                            })}
                        </div>
                    </div>

                    <hr style={{ borderColor: 'var(--glass-border-color)', margin: 'var(--space-4) 0', borderBottom: 'none' }} />

                    {/* Time Calendar Config */}
                    <div>
                        <h4 style={{ margin: '0 0 16px', color: 'var(--color-text-main)', font: 'var(--text-h4)' }}>Time Calendar Horizons</h4>
                        <div style={{ display: 'flex', gap: 'var(--space-4)', flexWrap: 'wrap' }}>
                            <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)', flex: 1, minWidth: '200px' }}>
                                <label style={{ font: 'var(--text-micro)', color: 'rgba(255,255,255,0.6)', textTransform: 'uppercase' }}>Calendar Resolution</label>
                                <select
                                    value={calendar.calendarType} onChange={e => setCalendar({ ...calendar, calendarType: e.target.value as any })}
                                    style={{ padding: 'var(--space-3)', backgroundColor: 'rgba(0,0,0,0.3)', border: '1px solid var(--glass-border-color)', borderRadius: 'var(--radius-sm)', color: 'white' }}
                                >
                                    <option>Yearly</option><option>Quarterly</option><option>Monthly</option><option>Weekly</option>
                                </select>
                            </div>
                            <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)', flex: 1, minWidth: '200px' }}>
                                <label style={{ font: 'var(--text-micro)', color: 'rgba(255,255,255,0.6)', textTransform: 'uppercase' }}>Fiscal Start Month</label>
                                <select
                                    value={calendar.fiscalStartMonth} onChange={e => setCalendar({ ...calendar, fiscalStartMonth: e.target.value })}
                                    style={{ padding: 'var(--space-3)', backgroundColor: 'rgba(0,0,0,0.3)', border: '1px solid var(--glass-border-color)', borderRadius: 'var(--radius-sm)', color: 'white' }}
                                >
                                    {['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'].map(m => <option key={m}>{m}</option>)}
                                </select>
                            </div>
                            <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)', flex: 1, minWidth: '100px' }}>
                                <label style={{ font: 'var(--text-micro)', color: 'rgba(255,255,255,0.6)', textTransform: 'uppercase' }}>Current FY</label>
                                <input type="number" value={calendar.currentFiscalYear} onChange={e => setCalendar({ ...calendar, currentFiscalYear: parseInt(e.target.value) })} style={{ padding: 'var(--space-3)', backgroundColor: 'rgba(0,0,0,0.3)', border: '1px solid var(--glass-border-color)', borderRadius: 'var(--radius-sm)', color: 'white' }} />
                            </div>
                            <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)', flex: 1, minWidth: '100px' }}>
                                <label style={{ font: 'var(--text-micro)', color: 'rgba(255,255,255,0.6)', textTransform: 'uppercase' }}>Past Years</label>
                                <input type="number" value={calendar.pastYearsCount} onChange={e => setCalendar({ ...calendar, pastYearsCount: parseInt(e.target.value) })} style={{ padding: 'var(--space-3)', backgroundColor: 'rgba(0,0,0,0.3)', border: '1px solid var(--glass-border-color)', borderRadius: 'var(--radius-sm)', color: 'white' }} />
                            </div>
                            <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)', flex: 1, minWidth: '100px' }}>
                                <label style={{ font: 'var(--text-micro)', color: 'rgba(255,255,255,0.6)', textTransform: 'uppercase' }}>Future Years</label>
                                <input type="number" value={calendar.futureYearsCount} onChange={e => setCalendar({ ...calendar, futureYearsCount: parseInt(e.target.value) })} style={{ padding: 'var(--space-3)', backgroundColor: 'rgba(0,0,0,0.3)', border: '1px solid var(--glass-border-color)', borderRadius: 'var(--radius-sm)', color: 'white' }} />
                            </div>
                        </div>
                    </div>

                    <div style={{ marginTop: 'var(--space-6)', display: 'flex', justifyContent: 'flex-end' }}>
                        <button
                            onClick={handleCreate}
                            style={{
                                padding: '12px 32px',
                                backgroundColor: 'var(--color-primary)',
                                color: 'white',
                                border: 'none',
                                borderRadius: 'var(--radius-md)',
                                font: 'var(--text-body)',
                                fontWeight: 'bold',
                                cursor: 'pointer',
                                boxShadow: '0 4px 12px rgba(59, 130, 246, 0.4)',
                                transition: 'all 0.1s'
                            }}
                        >
                            Initialize Application & Matrix →
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};
