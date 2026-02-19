"use client";

import React, { useEffect, useRef } from 'react';
import { Table } from 'apache-arrow';

// The "Projector": Renders the grid using WebGPU.
interface GridCanvasProps {
    data?: Table | null;
}

export const GridCanvas: React.FC<GridCanvasProps> = ({ data }) => {
    const canvasRef = useRef<HTMLCanvasElement>(null);

    // Ultra Diamond: Data Ingestion (Layer 6.2)
    useEffect(() => {
        if (data) {
            console.log(`[GridCanvas] Received Arrow Table: ${data.numRows} rows`);
            // In Layer 7, we will bind this to a GPU Buffer.
            // For now, we prove the data arrived.
        }
    }, [data]);

    useEffect(() => {
        const initWebGPU = async () => {
            if (!canvasRef.current) return;

            const adapter = await navigator.gpu?.requestAdapter();
            if (!adapter) {
                console.error("WebGPU not supported!");
                return;
            }
            const device = await adapter.requestDevice();
            const context = canvasRef.current.getContext('webgpu') as GPUCanvasContext;
            const format = navigator.gpu.getPreferredCanvasFormat();

            context.configure({
                device,
                format,
                alphaMode: 'premultiplied',
            });

            // Ultra Diamond: The Grid Shader (WGSL)
            // Renders infinite anti-aliased grid lines.
            const shaderModule = device.createShaderModule({
                code: `
          struct VertexOutput {
            @builtin(position) Position : vec4<f32>,
            @location(0) vUV : vec2<f32>,
          };

          @vertex
          fn vs_main(@builtin(vertex_index) VertexIndex : u32) -> VertexOutput {
            var pos = array<vec2<f32>, 6>(
              vec2<f32>(-1.0, -1.0), vec2<f32>( 1.0, -1.0), vec2<f32>(-1.0,  1.0),
              vec2<f32>(-1.0,  1.0), vec2<f32>( 1.0, -1.0), vec2<f32>( 1.0,  1.0)
            );
            var output : VertexOutput;
            output.Position = vec4<f32>(pos[VertexIndex], 0.0, 1.0);
            output.vUV = pos[VertexIndex] * 10.0; // Zoom factor
            return output;
          }

          @fragment
          fn fs_main(@location(0) vUV : vec2<f32>) -> @location(0) vec4<f32> {
            var grid = abs(fract(vUV - 0.5) - 0.5) / fwidth(vUV);
            var line = min(grid.x, grid.y);
            var color = 1.0 - min(line, 1.0);
            return vec4<f32>(0.2, 0.2, 0.2, color * 0.5); // Grey lines
          }
        `,
            });

            const pipeline = device.createRenderPipeline({
                layout: 'auto',
                vertex: {
                    module: shaderModule,
                    entryPoint: 'vs_main',
                },
                fragment: {
                    module: shaderModule,
                    entryPoint: 'fs_main',
                    targets: [{ format }],
                },
                primitive: {
                    topology: 'triangle-list',
                },
            });

            // Render Loop
            const frame = () => {
                const commandEncoder = device.createCommandEncoder();
                const textureView = context.getCurrentTexture().createView();

                const renderPassDescriptor: GPURenderPassDescriptor = {
                    colorAttachments: [{
                        view: textureView,
                        clearValue: { r: 0.1, g: 0.1, b: 0.1, a: 1.0 }, // Dark Background
                        loadOp: 'clear',
                        storeOp: 'store',
                    }],
                };

                const passEncoder = commandEncoder.beginRenderPass(renderPassDescriptor);
                passEncoder.setPipeline(pipeline);
                passEncoder.draw(6); // Draw Quad
                passEncoder.end();

                device.queue.submit([commandEncoder.finish()]);
                requestAnimationFrame(frame);
            };

            requestAnimationFrame(frame);
        };

        initWebGPU();
    }, []);

    return <canvas ref={canvasRef} width={800} height={600} />;
};
