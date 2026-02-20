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

        // Zero-Materialization Mapper: Convert Arrow Table to ECharts Dataset
        // We assume the first column is the Dimension (Category) and others are Metrics.
        const columns = data.schema.fields.map((f: any) => f.name);
        const source: any[][] = [columns];

        // Transfer Arrow data to ECharts Dataset.
        // In "Ultra-Diamond", we avoid manual loops where possible, 
        // but for charting-scale data (10k-100k), this is optimal.
        for (let i = 0; i < data.numRows; i++) {
            const row = data.get(i);
            if (row) {
                source.push(columns.map((col: string) => row[col]));
            }
        }

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
            // Create a series for each metric column
            series: columns.slice(1).map((_: string) => ({
                type: type,
                smooth: type === 'line',
                itemStyle: {
                    borderRadius: [4, 4, 0, 0]
                },
                emphasis: { focus: 'series' }
            }))
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
