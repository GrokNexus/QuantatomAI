"use client";

import React, { useEffect, useRef } from 'react';
import { Table } from 'apache-arrow';

// The "Projector": Renders the grid using WebGPU.
interface GridCanvasProps {
    data?: Table | null;
}

export const GridCanvas: React.FC<GridCanvasProps> = ({ data }) => {
    const canvasRef = useRef<HTMLCanvasElement>(null);

    // Ultra Diamond Vector 4: WebSocket Connection Saturation (10Hz Debounce)
    const cursorTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    const handleMouseMove = (e: React.MouseEvent<HTMLCanvasElement>) => {
        if (cursorTimeoutRef.current) return; // Drop events if we are within the 100ms throttle window

        cursorTimeoutRef.current = setTimeout(() => {
            // Emit PresenceUpdate over ConnectRPC WebSockets
            // e.g. stream.send({ ClientID: session.user, x: e.clientX, y: e.clientY })
            cursorTimeoutRef.current = null;
        }, 100); // 10Hz max transmission rate
    };

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

          struct EntropyUniforms {
            anomalyActive : f32, // 1.0 if anomaly, 0.0 otherwise
            time : f32,
          };
          @group(0) @binding(0) var<uniform> entropyState : EntropyUniforms;

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
            var colorVal = 1.0 - min(line, 1.0);
            
            // Base grid (grey)
            var baseColor = vec4<f32>(0.2, 0.2, 0.2, colorVal * 0.5); 
            
            // Phase 8.2 & 8.4: Entropy Spatial Shader (Orange Pulse on Anomaly)
            var anomalyColor = vec4<f32>(1.0, 0.5, 0.0, colorVal * (0.5 + 0.5 * sin(entropyState.time * 5.0))); 
            
            return mix(baseColor, anomalyColor, entropyState.anomalyActive);
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

            // Create Uniform Buffer for Entropy State
            const uniformBuffer = device.createBuffer({
                size: 8, // 2 f32s (anomalyActive padding, time)
                usage: GPUBufferUsage.UNIFORM | GPUBufferUsage.COPY_DST,
            });

            const bindGroup = device.createBindGroup({
                layout: pipeline.getBindGroupLayout(0),
                entries: [
                    {
                        binding: 0,
                        resource: { buffer: uniformBuffer }
                    }
                ]
            });

            // For simulation of the anomaly flash triggered by Redpanda streamer
            let anomalyActive = 0.0;
            // E.g., we could set anomalyActive = 1.0 when an event over WebSockets arrives.
            // For now, let's pulse it briefly every 5 seconds to show the capability
            setInterval(() => {
                anomalyActive = 1.0;
                setTimeout(() => anomalyActive = 0.0, 1000);
            }, 5000);

            let startTime = performance.now();

            // Render Loop
            const frame = () => {
                const currentTime = (performance.now() - startTime) / 1000.0;
                device.queue.writeBuffer(uniformBuffer, 0, new Float32Array([anomalyActive, currentTime]));

                const commandEncoder = device.createCommandEncoder();
                const textureView = context.getCurrentTexture().createView();

                const renderPassDescriptor: GPURenderPassDescriptor = {
                    colorAttachments: [{
                        view: textureView,
                        clearValue: { r: 0.1, g: 0.1, b: 0.1, a: 1.0 }, // Dark Background
                        loadOp: 'clear' as GPULoadOp,
                        storeOp: 'store' as GPUStoreOp,
                    }],
                };

                const passEncoder = commandEncoder.beginRenderPass(renderPassDescriptor);
                passEncoder.setPipeline(pipeline);
                passEncoder.setBindGroup(0, bindGroup);
                passEncoder.draw(6); // Draw Quad
                passEncoder.end();

                device.queue.submit([commandEncoder.finish()]);
                requestAnimationFrame(frame);
            };

            requestAnimationFrame(frame);
        };

        initWebGPU();
    }, []);

    return <canvas ref={canvasRef} width={800} height={600} onMouseMove={handleMouseMove} />;
};
