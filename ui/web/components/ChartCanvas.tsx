"use client";

import React, { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import { Table } from 'apache-arrow';

// Layer 7.6: Visual Intelligence (The Charting Projector)
// This component consumes the same Apache Arrow stream as the Grid.
interface ChartCanvasProps {
    data?: Table | null;
    type?: 'bar' | 'line' | 'scatter';
}

export const ChartCanvas: React.FC<ChartCanvasProps> = ({ data, type = 'bar' }) => {
    const options = useMemo(() => {
        if (!data || data.numRows === 0) {
            return {
                graphic: {
                    type: 'text',
                    left: 'center',
                    top: 'middle',
                    style: {
                        text: 'No Data Available',
                        fill: '#999',
                        font: '14px sans-serif'
                    }
                }
            };
        }

        // Ultra-Diamond Zero-Materialization Mapper: Convert Arrow Table to ECharts Dataset
        // We avoid row-by-row iteration which forces V8 to allocate millions of JS objects.
        // Instead, we extract the underlying C++ / WASM TypedArrays directly from the Arrow Columns.
        const columns = data.schema.fields.map((f: any) => f.name);

        // ECharts supports column-oriented dataset sources.
        const source: Record<string, any> = {};
        for (const colName of columns) {
            const vector = data.getChild(colName);
            // .toArray() returns the underlying TypedArray (e.g. Float64Array) without copying memory
            source[colName] = vector ? vector.toArray() : [];
        }

        // Create a series for each metric column
        const seriesDefinitions = columns.slice(1).map((colName: string) => {
            // Phase 8.4: Auto-Forecast Visual Designation
            // If a column name implies it's an AI Forecast, render as dashed line
            const isForecast = colName.toLowerCase().includes('forecast') || colName.toLowerCase().includes('fluxion');

            return {
                type: isForecast ? 'line' : type,
                smooth: type === 'line' || isForecast,
                itemStyle: {
                    borderRadius: [4, 4, 0, 0],
                    color: isForecast ? '#10b981' : undefined // Emerald green for Fluxion
                },
                lineStyle: isForecast ? { type: 'dashed', width: 2 } : undefined,
                emphasis: { focus: 'series' }
            };
        });

        return {
            backgroundColor: 'transparent',
            animation: true,
            tooltip: {
                trigger: 'axis',
                axisPointer: { type: 'cross' }
            },
            legend: {
                bottom: 0,
                textStyle: { color: '#ccc' }
            },
            grid: {
                top: '10%',
                left: '3%',
                right: '4%',
                bottom: '15%',
                containLabel: true
            },
            dataset: {
                dimensions: columns,
                source: source
            },
            xAxis: {
                type: 'category',
                axisLine: { lineStyle: { color: '#444' } },
                axisLabel: { color: '#888' },
                splitLine: { show: false }
            },
            yAxis: {
                type: 'value',
                axisLine: { lineStyle: { color: '#444' } },
                splitLine: { lineStyle: { color: '#222' } },
                axisLabel: { color: '#888' }
            },
            series: seriesDefinitions
        };
    }, [data, type]);

    return (
        <div style={{
            width: '100%',
            height: '500px',
            backgroundColor: 'rgba(23, 23, 23, 0.5)',
            borderRadius: '24px',
            border: '1px solid #262626',
            padding: '24px',
            boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.5)',
            backdropFilter: 'blur(12px)'
        }}>
            <ReactECharts
                option={options}
                style={{ height: '100%', width: '100%' }}
                notMerge={true}
                lazyUpdate={true}
            />
        </div>
    );
};
