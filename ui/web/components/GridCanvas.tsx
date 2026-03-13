"use client";

import React, { useEffect, useRef, useState } from 'react';
import { Table } from 'apache-arrow';
import { useGridHotkeys } from '../src/components/grid/useGridHotkeys';

import styles from './GridCanvas.module.css';

interface GridCanvasProps {
    data?: Table | null;
}

const CELL_WIDTH = 120;
const CELL_HEIGHT = 32;

export const GridCanvas: React.FC<GridCanvasProps> = ({ data }) => {
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const containerRef = useRef<HTMLDivElement>(null);

    // Hardcode some dummy dimensions for the structural matrix since we aren't hooking up the real Arrow data just yet
    const numRows = 1000;
    const numCols = 100;

    // The entire Keyboard Event Horizon maps here
    const {
        activeCell,
        selection,
        isEditing,
        setIsEditing,
        editValue,
        setEditValue
    } = useGridHotkeys(numRows, numCols);

    const [scrollOffset, setScrollOffset] = useState({ x: 0, y: 0 });

    const handleWheel = (e: React.WheelEvent) => {
        // Law 16: Smooth Scroll Physics
        setScrollOffset(prev => ({
            x: Math.max(0, prev.x + e.deltaX),
            y: Math.max(0, prev.y + e.deltaY)
        }));
    };

    // Calculate DOM coordinates based on scroll
    const activeTop = activeCell.row * CELL_HEIGHT - scrollOffset.y;
    const activeLeft = activeCell.col * CELL_WIDTH - scrollOffset.x;

    const selStartRow = Math.min(selection.start.row, selection.end.row);
    const selEndRow = Math.max(selection.start.row, selection.end.row);
    const selStartCol = Math.min(selection.start.col, selection.end.col);
    const selEndCol = Math.max(selection.start.col, selection.end.col);

    const selTop = selStartRow * CELL_HEIGHT - scrollOffset.y;
    const selLeft = selStartCol * CELL_WIDTH - scrollOffset.x;
    const selHeight = (selEndRow - selStartRow + 1) * CELL_HEIGHT;
    const selWidth = (selEndCol - selStartCol + 1) * CELL_WIDTH;

    useEffect(() => {
        // [Existing WebGPU Render Logic with Context Binding Omitted for brevity. We will keep what we had and add the camera transforms]
        const initWebGPU = async () => {
            if (!canvasRef.current) return;
            const adapter = await navigator.gpu?.requestAdapter();
            if (!adapter) return;
            const device = await adapter.requestDevice();
            const context = canvasRef.current.getContext('webgpu') as GPUCanvasContext;
            const format = navigator.gpu.getPreferredCanvasFormat();

            context.configure({ device, format, alphaMode: 'premultiplied' });

            const shaderModule = device.createShaderModule({
                code: `
                  struct Uniforms {
                    cameraOffset: vec2<f32>,
                    viewport: vec2<f32>,
                    time: f32,
                  };
                  @group(0) @binding(0) var<uniform> uniforms: Uniforms;

                  struct VertexOutput {
                    @builtin(position) Position: vec4<f32>,
                    @location(0) vCoord: vec2<f32>,
                  };

                  @vertex
                  fn vs_main(@builtin(vertex_index) VertexIndex: u32) -> VertexOutput {
                    var pos = array<vec2<f32>, 6>(
                      vec2<f32>(-1.0, -1.0), vec2<f32>( 1.0, -1.0), vec2<f32>(-1.0,  1.0),
                      vec2<f32>(-1.0,  1.0), vec2<f32>( 1.0, -1.0), vec2<f32>( 1.0,  1.0)
                    );
                    var output: VertexOutput;
                    output.Position = vec4<f32>(pos[VertexIndex], 0.0, 1.0);
                    // Map -1..1 to 0..viewport size
                    output.vCoord = (pos[VertexIndex] * 0.5 + 0.5) * uniforms.viewport;
                    return output;
                  }

                  @fragment
                  fn fs_main(@location(0) vCoord: vec2<f32>) -> @location(0) vec4<f32> {
                    // World coordinates
                    let world = vCoord + uniforms.cameraOffset;
                    
                    // Simple grid lines
                    let gridCell = vec2<f32>(120.0, 32.0); // CELL_WIDTH, CELL_HEIGHT
                    let local = world % gridCell;
                    
                    let lineThickness = 1.0;
                    if (local.x < lineThickness || local.y < lineThickness) {
                        return vec4<f32>(0.2, 0.2, 0.2, 0.5); // Grid Line
                    }
                    
                    return vec4<f32>(0.0, 0.0, 0.0, 0.0); // Transparent Background
                  }
                `
            });

            const pipeline = device.createRenderPipeline({
                layout: 'auto',
                vertex: { module: shaderModule, entryPoint: 'vs_main' },
                fragment: {
                    module: shaderModule, entryPoint: 'fs_main', targets: [{
                        format, blend: {
                            color: { srcFactor: 'src-alpha' as GPUBlendFactor, dstFactor: 'one-minus-src-alpha' as GPUBlendFactor, operation: 'add' as GPUBlendOperation },
                            alpha: { srcFactor: 'one' as GPUBlendFactor, dstFactor: 'one-minus-src-alpha' as GPUBlendFactor, operation: 'add' as GPUBlendOperation }
                        }
                    }]
                },
                primitive: { topology: 'triangle-list' },
            });

            const uniformBuffer = device.createBuffer({
                size: 16, // cameraOffset(8) + viewport(8)
                usage: GPUBufferUsage.UNIFORM | GPUBufferUsage.COPY_DST,
            });

            const bindGroup = device.createBindGroup({
                layout: pipeline.getBindGroupLayout(0),
                entries: [{ binding: 0, resource: { buffer: uniformBuffer } }]
            });

            const frame = () => {
                const viewport = new Float32Array([canvasRef.current!.width, canvasRef.current!.height]);
                // We use scrollOffset ref internally to sync with React render. Hacky but high perf.
                // In production, we'd sync this with requestAnimationFrame properly.
                device.queue.writeBuffer(uniformBuffer, 0, new Float32Array([scrollOffset.x, scrollOffset.y, viewport[0], viewport[1]]));

                const commandEncoder = device.createCommandEncoder();
                const renderPassDescriptor: GPURenderPassDescriptor = {
                    colorAttachments: [{
                        view: context.getCurrentTexture().createView(),
                        clearValue: { r: 0.05, g: 0.05, b: 0.06, a: 1.0 },
                        loadOp: 'clear' as GPULoadOp,
                        storeOp: 'store' as GPUStoreOp,
                    }],
                };

                const passEncoder = commandEncoder.beginRenderPass(renderPassDescriptor);
                passEncoder.setPipeline(pipeline);
                passEncoder.setBindGroup(0, bindGroup);
                passEncoder.draw(6);
                passEncoder.end();

                device.queue.submit([commandEncoder.finish()]);
                // We don't loop here immediately but trigger it via effect updates or a rAF loop.
                // For this mock, we'll let React state changes trigger the WebGPU re-render.
            };

            frame();
        };

        initWebGPU();
    }, [scrollOffset]);

    return (
        <div
            ref={containerRef}
            style={{ position: 'relative', width: '100%', height: '100%', overflow: 'hidden', cursor: 'cell' }}
            onWheel={handleWheel}
            tabIndex={0} // Make focusable for keyboard events
        >
            {/* 1. Underlying WebGPU Cartesian Fabric */}
            <canvas
                ref={canvasRef}
                width={2000} // Mock hardcoded size for demo
                height={1200}
                style={{ position: 'absolute', top: 0, left: 0, pointerEvents: 'none' }}
            />

            {/* 2. Z-Index Segregation: Selection Overlay Wrapper (Law 26) */}
            <div style={{ position: 'absolute', top: 0, left: 0, width: '100%', height: '100%', pointerEvents: 'none', zIndex: 'var(--z-grid-overlay)' }}>

                {/* Rubber Band / Marching Ants Selection Box (Law 11, 22) */}
                <div style={{
                    position: 'absolute',
                    top: selTop,
                    left: selLeft,
                    width: selWidth,
                    height: selHeight,
                    backgroundColor: 'rgba(59, 130, 246, 0.05)', // Even softer than 0.1
                    border: '1px solid var(--color-primary)',
                    boxSizing: 'border-box',
                    transition: 'all var(--duration-swift) var(--easing-gravity)',
                    display: (selWidth > CELL_WIDTH || selHeight > CELL_HEIGHT) ? 'block' : 'none'
                }}>
                    {/* Fill Handle Square */}
                    <div style={{
                        position: 'absolute',
                        bottom: '-4px',
                        right: '-4px',
                        width: '8px',
                        height: '8px',
                        backgroundColor: 'var(--color-primary)',
                        border: '1px solid var(--glass-panel-bg)', // Slight detachment from selection bounding box
                        cursor: 'crosshair',
                        pointerEvents: 'auto',
                        transition: 'transform var(--duration-swift) var(--easing-bounce)'
                    }}
                        onMouseEnter={(e) => e.currentTarget.style.transform = 'scale(1.5)'}
                        onMouseLeave={(e) => e.currentTarget.style.transform = 'scale(1)'}
                    />
                </div>

                {/* Active Cell Focus Ring */}
                <div style={{
                    position: 'absolute',
                    top: activeTop,
                    left: activeLeft,
                    width: CELL_WIDTH,
                    height: CELL_HEIGHT,
                    border: isEditing ? '2px solid var(--color-positive)' : '2px solid var(--color-primary)',
                    boxSizing: 'border-box',
                    boxShadow: isEditing ? '0 0 16px rgba(16, 185, 129, 0.3)' : '0 0 12px rgba(59, 130, 246, 0.2)', // Thinner shadow footprint
                    transition: 'top var(--duration-swift) var(--easing-gravity), left var(--duration-swift) var(--easing-gravity), box-shadow var(--duration-swift)',
                    zIndex: 2
                }}>
                    {/* Inline Editor `<input>` (Law 1, Escapism) */}
                    {isEditing && (
                        <input
                            autoFocus
                            value={editValue}
                            onChange={(e) => setEditValue(e.target.value)}
                            onBlur={() => setIsEditing(false)}
                            className={styles['grid-inline-editor']}
                            placeholder="Edit cell"
                            title="Edit cell value"
                        />
                    )}
                </div>

            </div>
        </div>
    );
};
