"use client";

import React, { useState } from 'react';

interface HistorySliderProps {
    minDate: Date;
    maxDate: Date;
    currentDate: Date;
    onDateChange: (date: Date) => void;
}

// Law 24 & CRDT Time-Machine Hook
export const HistorySlider: React.FC<HistorySliderProps> = ({ minDate, maxDate, currentDate, onDateChange }) => {
    const minTime = minDate.getTime();
    const maxTime = maxDate.getTime();
    const [sliderValue, setSliderValue] = useState(currentDate.getTime());

    const handleSliderChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const newTime = parseInt(e.target.value, 10);
        setSliderValue(newTime);
        onDateChange(new Date(newTime));
    };

    const percentage = ((sliderValue - minTime) / (maxTime - minTime)) * 100;

    return (
        <div style={{
            position: 'absolute',
            bottom: 'var(--space-6)',
            left: '50%',
            transform: 'translateX(-50%)',
            backgroundColor: 'var(--glass-panel-bg)',
            backdropFilter: 'var(--glass-panel-blur)',
            border: 'var(--glass-panel-border)',
            borderRadius: 'var(--radius-pill)',
            padding: 'var(--space-2) var(--space-4)',
            boxShadow: 'var(--glass-shadow-soft)',
            display: 'flex',
            alignItems: 'center',
            gap: 'var(--space-4)',
            zIndex: 'var(--z-shell)', // 20
            width: '600px',
            animation: 'fadeIn var(--duration-fluid) var(--easing-spring)'
        }}>
            <button style={{
                backgroundColor: 'rgba(255,255,255,0.05)',
                border: '1px solid var(--glass-border-color)',
                color: 'var(--color-text-main)',
                borderRadius: '50%',
                width: '32px', height: '32px',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                cursor: 'pointer',
                transition: 'all var(--duration-swift)'
            }}>
                <span className="google-symbols" style={{ fontSize: '18px' }}>history</span>
            </button>

            <div style={{ flex: 1, display: 'flex', flexDirection: 'column', gap: '4px' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', font: 'var(--text-micro)', color: 'var(--color-text-dim)' }}>
                    <span>{minDate.toLocaleDateString()}</span>
                    <strong style={{ color: 'var(--color-primary)' }}>{new Date(sliderValue).toLocaleString()}</strong>
                    <span>{maxDate.toLocaleDateString()}</span>
                </div>

                {/* Custom CSS Range Slider */}
                <input
                    type="range"
                    min={minTime}
                    max={maxTime}
                    value={sliderValue}
                    onChange={handleSliderChange}
                    style={{
                        width: '100%',
                        WebkitAppearance: 'none',
                        background: `linear-gradient(to right, var(--color-primary) ${percentage}%, rgba(255,255,255,0.1) ${percentage}%)`,
                        height: '4px',
                        borderRadius: '2px',
                        outline: 'none',
                        cursor: 'pointer'
                    }}
                />
            </div>
        </div>
    );
};
